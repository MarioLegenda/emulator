package singleFileExecution

import (
	"github.com/google/uuid"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/builders"
	"therebelsource/emulator/runner"
)

var SingleFileExecutionService Service

func createCommand(params interface{}, lang *runner.Language, containerName string) []string {
	commandFactory := runner.RunCommandFactory{}

	br := params.(builders.SingleFileBuildResult)

	return commandFactory.CreateCommand(containerName, br.ExecutionDirectory, br.FileName, lang, br.DirectoryName)
}

type Service struct {}

func InitService() {
	SingleFileExecutionService = Service{}
}

func (s Service) RunSingleFile(model *SingleFileRunRequest) (runner.SingleFileRunResult, *appErrors.Error) {
	builder := builders.CreateBuilder("single_file").(builders.SingleFileRunFn)

	buildResult, err := builder(model.codeBlock, model.State)

	if err != nil {
		return runner.SingleFileRunResult{}, err
	}

	builtRunner := runner.CreateRunner("singleFile").(runner.SingleFileRunFn)

	containerName := uuid.New().String()

	runResult, err := builtRunner(runner.SingleFileBuildResult{
		ContainerName: containerName,
		DirectoryName:     buildResult.DirectoryName,
		ExecutionDirectory: buildResult.ExecutionDirectory,
		FileName:           buildResult.FileName,
		Environment: model.codeBlock.Emulator,
		StateDirectory: buildResult.StateDirectory,
		Args: createCommand(buildResult, model.codeBlock.Emulator, containerName),
	})

	if err != nil {
		return runner.SingleFileRunResult{}, err
	}

	destroyRunner := builders.CreateDestroyer().(builders.SingleFileDestroyFn)

	if err := destroyRunner(buildResult); err != nil {
		return runner.SingleFileRunResult{}, nil
	}

	return runResult, nil
}
