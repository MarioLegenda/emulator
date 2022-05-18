package containerFactory

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"os"
	"os/exec"
	"sync"
	"therebelsource/emulator/appErrors"
	"time"
)

var services map[string]Container

type Container interface {
	CreateContainers(string, int) []*appErrors.Error
	Close()
	Containers() map[string]container
}

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

func Init(name string) {
	if services == nil {
		services = make(map[string]Container)
	}

	s := &service{containers: make(map[string]container)}

	services[name] = s
}

func Service(name string) Container {
	return services[name]
}

func (d *service) Containers() map[string]container {
	return d.containers
}

func (d *service) CreateContainers(tag string, workerNum int) []*appErrors.Error {
	fmt.Println(fmt.Sprintf("Creating %d container(s) for %s", workerNum, tag))

	errs := make([]*appErrors.Error, 0)
	wg := sync.WaitGroup{}
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			name := uuid.New().String()

			containerDir := fmt.Sprintf("%s/%s", os.Getenv("SINGLE_FILE_STATE_DIR"), name)
			fsErr := os.Mkdir(containerDir, os.ModePerm)

			if fsErr != nil {
				errs = append(errs, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Could not start container: %s", fsErr.Error())))

				wg.Done()

				return
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
			case <-time.After(5 * time.Second):
				if !isContainerRunning(name) {
					errs = append(errs, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Container startup timeout"))

					wg.Done()

					return
				}

				close(container.output)
				d.containers[name] = container

				fmt.Println(fmt.Sprintf("Container for %s with name %s started!", container.Tag, container.Name))

				wg.Done()
			case msg := <-container.output:
				if msg.messageType == "error" {
					err := msg.data.(error)

					errs = append(errs, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Could not start container: %s", err.Error())))

					wg.Done()
				}
			}
		}(&wg)
	}

	wg.Wait()

	return errs
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
			"-v",
			fmt.Sprintf("%s:/app:rw", getVolumeDirectory(c.Name)),
			"--name",
			c.Name,
			"--init",
			c.Tag,
			"/bin/sh",
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
