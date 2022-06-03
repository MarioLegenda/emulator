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
	"therebelsource/emulator/logger"
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
	lock       sync.Mutex
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
	logger.Info(fmt.Sprintf("Creating %d container(s) for %s", workerNum, tag))

	blocks := makeBlocks(workerNum, 10)

	errs := make([]*appErrors.Error, 0)
	for _, block := range blocks {
		wg := sync.WaitGroup{}

		for _ = range block {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				name := uuid.New().String()

				containerDir := fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), name)
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
				case <-time.After(1 * time.Second):
					if !isContainerRunning(name) {
						errs = append(errs, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Container startup timeout: Tag: %s, Name: %s", container.Tag, container.Name)))

						wg.Done()

						return
					}

					close(container.output)
					d.lock.Lock()
					d.containers[name] = container
					d.lock.Unlock()

					wg.Done()
				case msg := <-container.output:
					if msg.messageType == "error" {
						err := msg.data.(error)
						close(container.output)

						errs = append(errs, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Could not start container; Name: %s, Tag: %s: %s", container.Name, container.Tag, err.Error())))

						wg.Done()
					}
				}
			}(&wg)
		}

		wg.Wait()
	}

	time.Sleep(2 * time.Second)

	return errs
}

func (d service) Close() {
	contArr := containersToSlice(d.containers)
	blocks := makeBlocks(len(contArr), 10)

	for _, block := range blocks {
		wg := sync.WaitGroup{}

		for _, b := range block {
			wg.Add(1)
			c := contArr[b]

			go func(c container, wg *sync.WaitGroup) {
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

				stopDockerContainer(c.Name, pid)

				close(c.pid)

				err := os.RemoveAll(c.dir)

				if err == nil {
					// TODO: // send slack error and log
				}

				if err != nil {
					cmd := exec.Command("rm", []string{"-rf", c.dir}...)

					err := cmd.Run()

					if err != nil {
						wg.Done()
						// TODO: send slack error and log
						return
					}
				}

				wg.Done()
			}(c, &wg)
		}

		wg.Wait()
	}

	cmd := exec.Command("docker", []string{"volume", "rm", "$(docker volume ls -q)"}...)
	cmd.Run()
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
