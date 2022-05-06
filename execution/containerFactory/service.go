package containerFactory

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"os"
	"os/exec"
	errorHandler "therebelsource/emulator/appErrors"
	"time"
)

type Service interface {
	CreateContainers(string, int) bool
	Close()
	Containers() map[string]container
}

var PackageService Service

type service struct {
	containers map[string]container
}

type message struct {
	messageType string
	data        interface{}
}

type container struct {
	output chan message
	pid    chan int
	dir    string

	Tag       string
	Name      string
	WorkerNum int
}

func InitService() {
	s := &service{containers: make(map[string]container)}

	PackageService = s
}

func (d *service) Containers() map[string]container {
	return d.containers
}

func (d *service) CreateContainers(tag string, workerNum int) bool {
	fmt.Println(fmt.Sprintf("Creating %d workers for %s", workerNum, tag))
	for i := 0; i < workerNum; i++ {
		name := uuid.New().String()

		containerDir := fmt.Sprintf("%s/%s", os.Getenv("SINGLE_FILE_STATE_DIR"), name)
		fsErr := os.Mkdir(containerDir, os.ModePerm)

		if fsErr != nil {
			d.Close()

			errorHandler.TerminateWithMessage(fmt.Sprintf("Cannot create %s directory", containerDir))
		}

		container := container{
			output: make(chan message),
			pid:    make(chan int),
			dir:    containerDir,

			Tag:       tag,
			Name:      name,
			WorkerNum: workerNum,
		}

		createContainer(container)

		select {
		case <-time.After(1 * time.Second):
			if !isContainerRunning(name) {
				d.Close()

				return false
			}

			close(container.output)
			d.containers[name] = container

			fmt.Println(fmt.Sprintf("Container for %s with name %s started!", container.Tag, container.Name))
		case msg := <-container.output:
			if msg.messageType == "error" {
				err := msg.data.(error)

				fmt.Println(fmt.Sprintf("Container with name %s could not start: %s; Cannot continue, stopping all container!", name, err.Error()))

				d.Close()

				return false
			}
		}
	}

	return true
}

func (d service) Close() {
	for _, c := range d.containers {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(20*time.Second))
		out := make(chan int)

		go func(pidCh chan int, out chan int) {
			select {
			case <-ctx.Done():
				out <- 0
			case pid := <-c.pid:
				cancel()
				out <- pid
			}
		}(c.pid, out)

		pid := <-out

		fmt.Println(fmt.Sprintf("Stopping container! Name: %s, PID: %d", c.Name, pid))
		stopDockerContainer(c.Name, pid)

		close(c.pid)

		fmt.Println("Removing associated file system volume...")
		err := os.RemoveAll(c.dir)

		if err == nil {
			fmt.Println(fmt.Sprintf("Volume for %s removed", c.Name))
		}

		if err != nil {
			cmd := exec.Command("rm", []string{"-rf", c.dir}...)

			err := cmd.Run()

			if err != nil {
				fmt.Println(fmt.Sprintf("Failed to remove volume for %s: %s", c.Name, err.Error()))

				return
			}

			fmt.Println(fmt.Sprintf("Volume for %s removed", c.Name))
		}
	}

	fmt.Println("Removing leftover volumes...")
	cmd := exec.Command("docker", []string{"volume", "rm", "$(docker volume ls -q)"}...)
	err := cmd.Run()

	if err != nil {
		fmt.Println("Volumes could not be removed. Either there are no volumes to be removed or you must do it manually!")
	}

	fmt.Println("Volumes removed!")
}

func createContainer(c container) {
	go func(c container) {
		args := []string{
			"run",
			"-d",
			"-t",
			"--network=none",
			"--read-only",
			"--rm",
			"-v",
			fmt.Sprintf("%s:/app:rw", getVolumeDirectory(c.Name)),
			"--name",
			c.Name,
			"--init",
			c.Tag,
		}

		cmd := exec.Command("docker", args...)

		var outb, errb bytes.Buffer

		cmd.Stderr = &errb
		cmd.Stdout = &outb

		startErr := cmd.Run()

		if startErr != nil {
			c.output <- message{
				messageType: "error",
				data:        startErr,
			}

			return
		}

		c.pid <- cmd.Process.Pid
	}(c)
}
