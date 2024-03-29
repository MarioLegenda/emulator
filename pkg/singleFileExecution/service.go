package singleFileExecution

import (
	"emulator/pkg/appErrors"
	"emulator/pkg/execution"
	"emulator/pkg/repository"
	_var "emulator/var"
)

var SingleFileExecutionService Service

type Service struct{}

func InitService() {
	SingleFileExecutionService = Service{}
}

func (s Service) RunSingleFile(model *SingleFileRunRequest) (repository.RunResult, *appErrors.Error) {
	model.Sanitize()

	res := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
		BuilderType:       "single_file",
		ExecutionType:     "single_file",
		EmulatorName:      string(model.codeBlock.Emulator.Name),
		EmulatorTag:       string(model.codeBlock.Emulator.Tag),
		EmulatorExtension: model.codeBlock.Emulator.Extension,
		EmulatorText:      model.codeBlock.Text,
	})

	result := res.Result

	if result == "" && res.Error != nil && appErrors.TimeoutError == res.Error.Code {
		result = "timeout"
	}

	return repository.RunResult{
		Success: res.Success,
		Result:  result,
		Timeout: 5,
	}, nil
}

func (s Service) RunSnippet(model *SnippetRequest) (repository.RunResult, *appErrors.Error) {
	model.Sanitize()

	res := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
		BuilderType:       "single_file",
		ExecutionType:     "single_file",
		EmulatorName:      string(model.snippet.SnippetCodeBlock.Emulator.Name),
		EmulatorTag:       string(model.snippet.SnippetCodeBlock.Emulator.Tag),
		EmulatorExtension: model.snippet.SnippetCodeBlock.Emulator.Extension,
		EmulatorText:      model.snippet.SnippetCodeBlock.Text,
	})

	result := res.Result

	if result == "" && res.Error != nil && appErrors.TimeoutError == res.Error.Code {
		result = "timeout"
	}

	return repository.RunResult{
		Success: res.Success,
		Result:  result,
		Timeout: 5,
	}, nil
}

func (s Service) RunPublicSnippet(model *PublicSnippetRequest) (repository.RunResult, *appErrors.Error) {
	model.Sanitize()

	res := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
		BuilderType:       "single_file",
		ExecutionType:     "single_file",
		EmulatorName:      string(model.snippet.SnippetCodeBlock.Emulator.Name),
		EmulatorTag:       string(model.snippet.SnippetCodeBlock.Emulator.Tag),
		EmulatorExtension: model.snippet.SnippetCodeBlock.Emulator.Extension,
		EmulatorText:      model.Text,
	})

	result := res.Result

	if result == "" && res.Error != nil && appErrors.TimeoutError == res.Error.Code {
		result = "timeout"
	}

	return repository.RunResult{
		Success: res.Success,
		Result:  result,
		Timeout: 5,
	}, nil
}

func (s Service) RunPublicSingleFile(model *PublicSingleFileRunRequest) (repository.RunResult, *appErrors.Error) {
	model.Sanitize()

	res := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
		BuilderType:       "single_file",
		ExecutionType:     "single_file",
		EmulatorName:      string(model.codeBlock.Emulator.Name),
		EmulatorTag:       string(model.codeBlock.Emulator.Tag),
		EmulatorExtension: model.codeBlock.Emulator.Extension,
		EmulatorText:      model.codeBlock.Text,
	})

	result := res.Result

	if result == "" && res.Error != nil && appErrors.TimeoutError == res.Error.Code {
		result = "timeout"
	}

	return repository.RunResult{
		Success: res.Success,
		Result:  result,
		Timeout: 5,
	}, nil
}
