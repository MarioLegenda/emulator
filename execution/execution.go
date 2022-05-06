package execution

import (
	"fmt"
	"sync"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/execution/balancer"
	"therebelsource/emulator/execution/balancer/runners"
	"therebelsource/emulator/execution/containerFactory"
	"therebelsource/emulator/runner"
)

var PackageService Execution

type Job struct {
	BuilderType   string
	ExecutionType string

	EmulatorName      string
	EmulatorExtension string
	EmulatorTag       string
	EmulatorText      string
}

type Execution interface {
	Close()
	RunJob(j Job) runners.Result
}

type execution struct {
	controller map[string][]int32
	balancers  map[string][]balancer.Balancer
	lock       sync.Mutex
	close      bool
}

type containerBlueprint struct {
	workerNum int
	tag       string
}

func Init(workerNum int, containerNum int) *appErrors.Error {
	containerFactory.InitService()
	s := &execution{
		balancers:  make(map[string][]balancer.Balancer),
		controller: make(map[string][]int32),
	}

	err := s.init(workerNum, containerNum)

	if err != nil {
		return err
	}

	PackageService = s

	return nil
}

func (e *execution) RunJob(j Job) runners.Result {
	e.lock.Lock()

	balancers := e.balancers[j.EmulatorTag]
	controller := e.controller[j.EmulatorTag]

	if e.close {
		e.lock.Unlock()

		return runners.Result{
			Result:  "",
			Success: false,
			Error:   appErrors.New(appErrors.ApplicationError, appErrors.TimeoutError, "Code execution timeout!"),
		}
	}

	idx := 0
	first := controller[0]
	for i, r := range controller {
		if r < first {
			idx = i
		}
	}

	e.controller[j.EmulatorTag][idx] = e.controller[j.EmulatorTag][idx] + 1

	e.lock.Unlock()

	b := balancers[idx]

	output := make(chan runners.Result)
	b.AddJob(balancer.Job{
		BuilderType:       j.BuilderType,
		ExecutionType:     j.ExecutionType,
		EmulatorName:      j.EmulatorName,
		EmulatorExtension: j.EmulatorExtension,
		EmulatorText:      j.EmulatorText,
		Output:            output,
	})

	out := <-output
	close(output)

	e.lock.Lock()
	e.controller[j.EmulatorTag][idx] = e.controller[j.EmulatorTag][idx] - 1
	e.lock.Unlock()

	return out
}

func (e *execution) Close() {
	defer func() {
		if err := recover(); err != nil {
			// send to slack/log
		}
	}()

	e.lock.Lock()
	e.close = true
	e.lock.Unlock()

	for _, balancers := range e.balancers {
		for _, b := range balancers {
			b.Close()
		}
	}

	containerFactory.PackageService.Close()
}

func (e *execution) init(workerNum int, containerNum int) *appErrors.Error {
	blueprints := []containerBlueprint{
		{
			workerNum: containerNum,
			tag:       string(runner.NodeLts.Tag),
		},
		{
			workerNum: containerNum,
			tag:       string(runner.GoLang.Tag),
		},
	}

	for _, blueprint := range blueprints {
		success := containerFactory.PackageService.CreateContainers(blueprint.tag, blueprint.workerNum)

		if !success {
			return appErrors.New(appErrors.ServerError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Cannot boot container for tag %s", blueprint.tag))
		}
	}

	containers := containerFactory.PackageService.Containers()

	for _, c := range containers {
		b := balancer.NewBalancer(c.Name, workerNum)
		b.StartWorkers()
		e.balancers[c.Tag] = make([]balancer.Balancer, 0)
		e.balancers[c.Tag] = append(e.balancers[c.Tag], b)

		e.controller[c.Tag] = make([]int32, 0)
		e.controller[c.Tag] = append(e.controller[c.Tag], 0)
	}

	return nil
}
