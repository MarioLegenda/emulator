package projectExecution

import (
	"emulator/pkg/appErrors"
	"emulator/pkg/execution"
	"emulator/pkg/repository"
	_var "emulator/var"
)

var ProjectExecutionService Service

type Service struct{}

func InitService() {
	ProjectExecutionService = Service{}
}

func (s Service) RunProject(model *ProjectRunRequest) (repository.ProjectRunResult, *appErrors.Error) {
	res := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
		BuilderType:   "project",
		ExecutionType: "project",
		EmulatorTag:   string(model.sessionData.CodeProject.Environment.Tag),
		EmulatorName:  string(model.sessionData.CodeProject.Environment.Name),
		CodeProject:   model.sessionData.CodeProject,
		ExecutingFile: model.sessionData.ExecutingFile,
		Contents:      model.sessionData.Content,
		PackageName:   model.sessionData.PackageName,
	})

	result := res.Result

	if result == "" && res.Error != nil && appErrors.TimeoutError == res.Error.Code {
		result = "timeout"
	}

	return repository.ProjectRunResult{
		Success: res.Success,
		Result:  result,
		Timeout: 5,
	}, nil
}
