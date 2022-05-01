package execution

import (
	"fmt"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/execution/containerFactory"
	"therebelsource/emulator/runner"
)

var PackageService Execution

type Execution interface {
	Close()
}

type execution struct{}

func Init(workerNum int) *appErrors.Error {
	containerFactory.InitService(workerNum)
	s := &execution{}

	err := s.init()

	if err != nil {
		return err
	}

	PackageService = s

	return nil
}

func (e *execution) init() *appErrors.Error {
	tags := []string{
		string(runner.NodeEsm.Tag),
	}

	for _, tag := range tags {
		success := containerFactory.PackageService.CreateContainers(tag)

		if !success {
			return appErrors.New(appErrors.ServerError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Cannot boot container for tag %s", tag))
		}
	}

	return nil
}

func (e *execution) Close() {
	containerFactory.PackageService.Close()
}
