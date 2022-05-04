package balancer

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/execution/balancer/builders"
	"therebelsource/emulator/execution/balancer/runners"
	"time"
)

type Balancer interface {
	StartWorkers()
	AddJob(Job)
	Close()
}

type Job struct {
	BuilderType   string
	ExecutionType string

	EmulatorName      string
	EmulatorExtension string
	EmulatorText      string

	Output chan runners.Result
}

type worker struct {
	input chan Job
	name  string
	index int
}

type balancer struct {
	workers    []worker
	name       string
	controller []int32
	lock       sync.Mutex
	close      bool
}

func NewBalancer(name string, initialWorkers int) Balancer {
	b := &balancer{
		workers:    make([]worker, 0),
		controller: make([]int32, 0),
		name:       name,
	}

	for i := 0; i < initialWorkers; i++ {
		b.workers = append(b.workers, worker{
			input: make(chan Job),
			name:  name,
			index: i,
		})

		var p int32

		atomic.StoreInt32(&p, 0)

		b.controller = append(b.controller, p)
	}

	return b
}

func (b *balancer) StartWorkers() {
	wg := sync.WaitGroup{}
	for _, w := range b.workers {
		wg.Add(1)
		go func(worker worker, wg *sync.WaitGroup) {
			wg.Done()

			var current int32 = 0
			for job := range worker.input {
				if current == 0 && b.close {
					job.Output <- runners.Result{
						Result:  "",
						Success: false,
						Error:   appErrors.New(appErrors.ApplicationError, appErrors.ShutdownError, "Worker is shutting down!"),
					}
					return
				}

				build, err := builders.NodeSingleFileBuild(builders.InitNodeParams(
					job.EmulatorExtension,
					job.EmulatorText,
					fmt.Sprintf("%s/%s", os.Getenv("SINGLE_FILE_STATE_DIR"), worker.name),
				))

				if err != nil {
					job.Output <- runners.Result{
						Result:  "",
						Success: false,
						Error:   err,
					}

					b.lock.Lock()
					b.controller[worker.index] = b.controller[worker.index] - 1
					current = b.controller[worker.index]
					b.lock.Unlock()

					return
				}

				res := runners.NodeRunner(runners.NodeExecParams{
					ExecutionDirectory: build.ExecutionDirectory,
					ContainerDirectory: build.ContainerDirectory,
					ExecutionFile:      build.FileName,
					ContainerName:      worker.name,
				})

				job.Output <- res

				b.lock.Lock()
				b.controller[worker.index] = b.controller[worker.index] - 1
				current = b.controller[worker.index]
				b.lock.Unlock()
			}
		}(w, &wg)
	}

	wg.Wait()
}

func (b *balancer) AddJob(j Job) {
	b.lock.Lock()

	if b.close {
		b.lock.Unlock()

		return
	}

	idx := 0
	first := b.controller[0]
	for i, r := range b.controller {
		if r < first {
			idx = i
		}
	}

	b.controller[idx] = b.controller[idx] + 1

	b.lock.Unlock()

	b.workers[idx].input <- j
}

func (b *balancer) Close() {
	b.lock.Lock()
	b.close = true
	b.lock.Unlock()

	isClosed := make(chan bool)
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(20*time.Second))
	go func() {
		for {
			select {
			case <-ctx.Done():
				isClosed <- true
			default:
				l := len(b.controller) - 1
				a := 0
				for _, r := range b.controller {
					if r == 0 {
						a++
					}

					if a == l {
						isClosed <- true
						return
					}
				}
			}
		}
	}()

	<-isClosed

	for _, w := range b.workers {
		close(w.input)
	}
}
