package linkedProjectExecution

import (
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/execution"
	"therebelsource/emulator/repository"
	_var "therebelsource/emulator/var"
)

var ExecutionService Service

type Service struct{}

func InitService() {
	ExecutionService = Service{}
}

func (s Service) RunProject(model *LinkedProjectRunRequest) (repository.ProjectRunResult, *appErrors.Error) {
	res := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
		BuilderType:   "linked",
		ExecutionType: "linked",
		EmulatorTag:   string(model.sessionData.CodeProject.Environment.Tag),
		EmulatorName:  string(model.sessionData.CodeProject.Environment.Name),
		EmulatorText:  model.sessionData.CodeBlock.Text,
		CodeProject:   model.sessionData.CodeProject,
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

func (s Service) RunPublicProject(model *PublicLinkedProjectRunRequest) (repository.ProjectRunResult, *appErrors.Error) {
	res := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
		BuilderType:   "linked",
		ExecutionType: "linked",
		EmulatorTag:   string(model.sessionData.CodeProject.Environment.Tag),
		EmulatorName:  string(model.sessionData.CodeProject.Environment.Name),
		EmulatorText:  model.Text,
		CodeProject:   model.sessionData.CodeProject,
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
