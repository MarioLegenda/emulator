package runner

var containerWorkers chan struct{
	containerName string
	pid int
}

var containersStoppedSignal chan bool

func WatchContainers() {
	containersStoppedSignal = make(chan bool)
	containerWorkers = make(chan struct{
		containerName string
		pid int
	})
	go func() {
		for c := range containerWorkers {
			stopDockerContainer(c.containerName, c.pid)
		}

		containersStoppedSignal<- true
	}()
}

func StopWatching() (int, int) {
	l := len(containerWorkers)
	close(containerWorkers)
	s := len(containerWorkers)
	<- containersStoppedSignal

	close(containersStoppedSignal)

	return l, s
}
