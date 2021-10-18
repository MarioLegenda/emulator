package projectExecution

import (
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/builders"
	"therebelsource/emulator/repository"
	"therebelsource/emulator/runner"
)

var ProjectExecutionService Service

type Service struct {}

func InitService() {
	ProjectExecutionService = Service{}
}

func (s Service) RunProject(model *CodeProjectRunRequest) (runner.ProjectRunResult, *appErrors.Error) {
	projectBuilder := builders.CreateBuilder("project").(builders.ProjectBuildFn)

	repository := repository.InitRepository()

	contents, err := repository.GetAllFileContent(model.CodeProjectUuid)

	if err != nil {
		return runner.ProjectRunResult{}, err
	}

	projectBuilder(model.codeProject, contents, "session")

	return runner.ProjectRunResult{}, nil
}
