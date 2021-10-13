package singleFileExecution

import (
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/builders"
	"therebelsource/emulator/runner"
)

var SingleFileExecutionService Service

type Service struct {}

func InitService() {
	SingleFileExecutionService = Service{}
}

func (s Service) RunSingleFile(model *SingleFileRunRequest) (runner.SingleFileRunResult, *appErrors.Error) {
	builder := builders.CreateBuilder(model.Type).(builders.SingleFileRunFn)

	buildResult, err := builder(model.codeBlock, model.State)

	if err != nil {
		return runner.SingleFileRunResult{}, err
	}

	builtRunner := runner.CreateRunner("singleFile").(runner.SingleFileRunFn)

	runResult, err := builtRunner(runner.SingleFileBuildResult{
		DirectoryName:     buildResult.DirectoryName,
		ExecutionDirectory: buildResult.ExecutionDirectory,
		FileName:           buildResult.FileName,
		Environment: model.codeBlock.Emulator,
		StateDirectory: buildResult.StateDirectory,
	})

	if err != nil {
		return runner.SingleFileRunResult{}, err
	}

	return runResult, nil
}
