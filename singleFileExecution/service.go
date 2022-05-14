package singleFileExecution

import (
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/execution"
	"therebelsource/emulator/runner"
)

var SingleFileExecutionService Service

type Service struct{}

func InitService() {
	SingleFileExecutionService = Service{}
}

func (s Service) RunSingleFile(model *SingleFileRunRequest) (runner.SingleFileRunResult, *appErrors.Error) {
	model.Sanitize()

	res := execution.PackageService.RunJob(execution.Job{
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

	return runner.SingleFileRunResult{
		Success: res.Success,
		Result:  result,
		Timeout: 5,
	}, nil
}

func (s Service) RunPublicSingleFile(model *PublicSingleFileRunRequest) (runner.SingleFileRunResult, *appErrors.Error) {
	model.Sanitize()

	res := execution.PackageService.RunJob(execution.Job{
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

	return runner.SingleFileRunResult{
		Success: res.Success,
		Result:  result,
		Timeout: 5,
	}, nil
}
