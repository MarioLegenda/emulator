package execution

import (
	"fmt"
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
	EmulatorText      string
}

type Execution interface {
	Close()
	RunJob(j Job) runners.Result
}

type execution struct {
	balancers map[string]balancer.Balancer
}

type containerBlueprint struct {
	workerNum int
	tag       string
}

func Init(workerNum int) *appErrors.Error {
	containerFactory.InitService(workerNum)
	s := &execution{
		balancers: make(map[string]balancer.Balancer),
	}

	err := s.init()

	if err != nil {
		return err
	}

	PackageService = s

	return nil
}

func (e *execution) RunJob(j Job) runners.Result {
	for _, b := range e.balancers {
		output := make(chan runners.Result)
		b.AddJob(balancer.Job{
			BuilderType:       j.BuilderType,
			ExecutionType:     j.ExecutionType,
			EmulatorName:      j.EmulatorName,
			EmulatorExtension: j.EmulatorExtension,
			EmulatorText:      j.EmulatorText,
			Output:            output,
		})

		return <-output
	}

	return runners.Result{}
}

func (e *execution) Close() {
	for _, b := range e.balancers {
		b.Close()
	}

	containerFactory.PackageService.Close()
}

func (e *execution) init() *appErrors.Error {
	blueprints := []containerBlueprint{
		{
			workerNum: 5,
			tag:       string(runner.NodeLts.Tag),
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
		b := balancer.NewBalancer(c.Name, 5)
		b.StartWorkers()
		e.balancers[c.Name] = b
	}

	return nil
}
