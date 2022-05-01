package containerFactory

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"os"
	"os/exec"
	errorHandler "therebelsource/emulator/appErrors"
	"time"
)

type Service interface {
	CreateContainers(string)
	StopContainers()
}

var PackageService Service

type service struct {
	containers map[string]container
	workerNum  int
}

type message struct {
	messageType string
	data        interface{}
}

type container struct {
	output chan message
	pid    chan int
	name   string
	tag    string
}

func InitService(workerNum int) {
	s := &service{containers: make(map[string]container), workerNum: workerNum}

	PackageService = s
}

func (d *service) CreateContainers(tag string) {
	fmt.Println(fmt.Sprintf("Creating %d workers for %s", d.workerNum, tag))
	for i := 0; i < d.workerNum; i++ {
		name := uuid.New().String()
		container := container{
			output: make(chan message),
			pid:    make(chan int),
			name:   name,
			tag:    tag,
		}

		containerDir := fmt.Sprintf("%s/%s", os.Getenv("SINGLE_FILE_STATE_DIR"), name)
		fsErr := os.Mkdir(containerDir, os.ModePerm)

		if fsErr != nil {
			errorHandler.TerminateWithMessage(fmt.Sprintf("Cannot create %s directory", containerDir))
		}

		createContainer(container)

		select {
		case <-time.After(5 * time.Second):
			if !isContainerRunning(name) {
				d.StopContainers()

				return
			}

			d.containers[name] = container

			fmt.Println(fmt.Sprintf("Container for %s with name %s started!", container.tag, container.name))
		case msg := <-container.output:
			if msg.messageType == "error" {
				err := msg.data.(error)

				fmt.Println(fmt.Sprintf("Container with name %s could not start: %s; Cannot continue, stopping all container!", name, err.Error()))

				d.StopContainers()
			}
		}
	}
}

func (d service) StopContainers() {
	for _, c := range d.containers {
		pid := <-c.pid

		fmt.Println(fmt.Sprintf("Stopping container! Name: %s, PID: %d", c.name, pid))
		stopDockerContainer(c.name, pid)

		close(c.output)
		close(c.pid)
	}
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
			fmt.Sprintf("%s:/app:rw", getVolumeDirectory(c.name)),
			"--name",
			c.name,
			"--init",
			c.tag,
		}

		cmd := exec.Command("docker", args...)

		var outb, errb bytes.Buffer

		cmd.Stderr = &errb
		cmd.Stdout = &outb

		startErr := cmd.Run()
		c.pid <- cmd.Process.Pid

		if startErr != nil {
			c.output <- message{
				messageType: "error",
				data:        startErr,
			}

			return
		}
	}(c)
}
