package runner

import (
	"sync"
)

var runnerBalancer *balancer

type job struct {
	containerName string
	pid int
}

type worker struct {
	comm chan job
}

type stopper struct {
	comm chan bool
}

type balancer struct {
	workers []worker
	stoppers []stopper
	counter int
	maxWorkers int
	sync.Mutex
}

func StartContainerBalancer() {
	newBalancer()

	runnerBalancer.StartWorkers()
}

func StopContainerBalancer() {
	runnerBalancer.StopWorkers()
}

func newBalancer() {
	maxWorkers := 100
	runnerBalancer = &balancer{
		workers:  make([]worker, 0),
		stoppers: make([]stopper, 0),
		counter:  0,
		maxWorkers: maxWorkers,
	}

	for i := 0; i < runnerBalancer.maxWorkers; i++ {
		runnerBalancer.workers = append(runnerBalancer.workers, worker{comm: make(chan job)})
		runnerBalancer.stoppers = append(runnerBalancer.stoppers, stopper{comm: make(chan bool)})
	}
}

func (b *balancer) StartWorkers() {
	for i, w := range b.workers {
		go func(count int, balancer *balancer, worker worker) {
			for c := range worker.comm {
				stopDockerContainer(c.containerName, c.pid)
			}
		}(i, b, w)
	}
}

func (b *balancer) StopWorkers() {
	for _, w := range b.workers {
		close(w.comm)
	}
}

func (b *balancer) addJob(j job) {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		b.Lock()
		if b.counter >= b.maxWorkers {
			b.counter = 0
		}

		b.workers[b.counter].comm<- j

		b.counter++

		b.Unlock()
		
		wg.Done()
	}(wg)

	wg.Wait()
}


