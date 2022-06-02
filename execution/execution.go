package execution

import (
	"fmt"
	"os"
	"sync"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/execution/balancer"
	"therebelsource/emulator/execution/balancer/runners"
	"therebelsource/emulator/execution/containerFactory"
	"therebelsource/emulator/logger"
	"therebelsource/emulator/repository"
	"therebelsource/emulator/slack"
)

var services map[string]Execution

type ProjectExecutionData struct {
}

type Job struct {
	BuilderType   string
	ExecutionType string

	EmulatorName      string
	EmulatorExtension string
	EmulatorTag       string
	EmulatorText      string
	PackageName       string

	CodeProject   *repository.CodeProject
	Contents      []*repository.FileContent
	ExecutingFile *repository.File
}

type Execution interface {
	Close()
	Closed() bool
	RunJob(j Job) runners.Result
}

type execution struct {
	controller map[string][]int32
	balancers  map[string][]balancer.Balancer
	lock       sync.Mutex
	close      bool
	name       string
}

type ContainerBlueprint struct {
	WorkerNum    int
	ContainerNum int
	Tag          string
}

func Init(name string, blueprints []ContainerBlueprint) *appErrors.Error {
	if services == nil {
		services = make(map[string]Execution)
	}

	containerFactory.Init(name)
	s := &execution{
		balancers:  make(map[string][]balancer.Balancer),
		controller: make(map[string][]int32),
		name:       name,
	}

	services[name] = s

	err := s.init(name, blueprints)

	if err != nil {
		return err
	}

	return nil
}

func Service(name string) Execution {
	return services[name]
}

func (e *execution) Closed() bool {
	return e.close
}

func (e *execution) RunJob(j Job) runners.Result {
	defer func() {
		if err := recover(); err != nil {
			slack.SendLog("Error", fmt.Sprintf("A panic occurred while running a job. The server will close right away and you have to clean up after it. Error: %v", err), "critical_log")
			logger.Error(fmt.Sprintf("A panic occurred while running a job. The server will close right away and you have to clean up after it. Error: %v", err))
			os.Exit(125)
		}
	}()

	e.lock.Lock()

	balancers := e.balancers[j.EmulatorTag]
	controller := e.controller[j.EmulatorTag]

	if e.close {
		e.lock.Unlock()

		return runners.Result{
			Result:  "",
			Success: false,
			Error:   appErrors.New(appErrors.ApplicationError, appErrors.TimeoutError, "Closing executioner"),
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

	b := balancers[idx]

	e.lock.Unlock()

	output := make(chan runners.Result)
	b.AddJob(balancer.Job{
		BuilderType:       j.BuilderType,
		ExecutionType:     j.ExecutionType,
		EmulatorName:      j.EmulatorName,
		EmulatorExtension: j.EmulatorExtension,
		EmulatorText:      j.EmulatorText,

		CodeProject:   j.CodeProject,
		ExecutingFile: j.ExecutingFile,
		Contents:      j.Contents,
		PackageName:   j.PackageName,

		Output: output,
	})

	out := <-output
	close(output)

	e.lock.Lock()
	e.controller[j.EmulatorTag][idx] = e.controller[j.EmulatorTag][idx] - 1
	e.lock.Unlock()

	return out
}

func (e *execution) Close() {
	e.lock.Lock()
	e.close = true
	e.lock.Unlock()

	for _, balancers := range e.balancers {
		for _, b := range balancers {
			b.Close()
		}
	}

	containerFactory.Service(e.name).Close()
}

func (e *execution) init(name string, blueprints []ContainerBlueprint) *appErrors.Error {
	workers := make(map[string]int)
	for _, blueprint := range blueprints {
		errs := containerFactory.Service(name).CreateContainers(blueprint.Tag, blueprint.ContainerNum)

		if len(errs) != 0 {
			e.Close()

			slack.SendLog("Error", fmt.Sprintf("Cannot boot container for tag %s: %v", blueprint.Tag, errs), "critical_log")
			logger.Error(fmt.Sprintf("Cannot boot container for tag %s: %v", blueprint.Tag, errs))

			return appErrors.New(appErrors.ServerError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Cannot boot container for tag %s", blueprint.Tag))
		}

		workers[blueprint.Tag] = blueprint.WorkerNum
	}

	containers := containerFactory.Service(name).Containers()

	for _, c := range containers {
		workerNum := workers[c.Tag]
		b := balancer.NewBalancer(c.Name, workerNum)
		b.StartWorkers()
		e.balancers[c.Tag] = make([]balancer.Balancer, 0)
		e.balancers[c.Tag] = append(e.balancers[c.Tag], b)

		e.controller[c.Tag] = make([]int32, 0)
		e.controller[c.Tag] = append(e.controller[c.Tag], 0)
	}

	return nil
}
