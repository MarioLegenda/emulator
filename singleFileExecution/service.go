package singleFileExecution

import (
	"github.com/google/uuid"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/builders"
	"therebelsource/emulator/execution"
	"therebelsource/emulator/runner"
)

var SingleFileExecutionService Service

func createCommand(params interface{}, lang *runner.Language, containerName string) []string {
	commandFactory := runner.RunCommandFactory{}

	br := params.(builders.SingleFileBuildResult)

	return commandFactory.CreateCommand(containerName, br.ExecutionDirectory, br.FileName, lang, br.DirectoryName)
}

type Service struct{}

func InitService() {
	SingleFileExecutionService = Service{}
}

func (s Service) RunSingleFile(model *SingleFileRunRequest) (runner.SingleFileRunResult, *appErrors.Error) {
	model.Sanitize()

	if model.codeBlock.Emulator.Name == "node_latest" {
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

	if model.codeBlock.Emulator.Name == "go" {
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

	if model.codeBlock.Emulator.Name == "ruby" {
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

	if model.codeBlock.Emulator.Name == "php74" {
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

	if model.codeBlock.Emulator.Name == "python2" {
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

	if model.codeBlock.Emulator.Name == "python3" {
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

	if model.codeBlock.Emulator.Name == "c_sharp_mono" {
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

	if model.codeBlock.Emulator.Name == "haskell" {
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

	if model.codeBlock.Emulator.Name == "c" {
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

	if model.codeBlock.Emulator.Name == "c++" {
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

	if model.codeBlock.Emulator.Name == "rust" {
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

	builder := builders.CreateBuilder("single_file").(builders.SingleFileRunFn)

	buildResult, err := builder(model.codeBlock, "single_file")

	if err != nil {
		return runner.SingleFileRunResult{}, err
	}

	builtRunner := runner.CreateRunner("singleFile").(runner.SingleFileRunFn)

	containerName := uuid.New().String()

	runResult, err := builtRunner(runner.SingleFileBuildResult{
		ContainerName:      containerName,
		DirectoryName:      buildResult.DirectoryName,
		ExecutionDirectory: buildResult.ExecutionDirectory,
		FileName:           buildResult.FileName,
		Environment:        model.codeBlock.Emulator,
		StateDirectory:     buildResult.StateDirectory,
		Args:               createCommand(buildResult, model.codeBlock.Emulator, containerName),
		Timeout:            model.validatedTemporarySession.Timeout,
	})

	if err != nil {
		return runner.SingleFileRunResult{}, err
	}

	destroyRunner := builders.CreateDestroyer("single_file").(builders.SingleFileDestroyFn)

	if err := destroyRunner(buildResult); err != nil {
		// log here if it fails, do not tell the user
		return runResult, nil
	}

	return runResult, nil
}

func (s Service) RunPublicSingleFile(model *PublicSingleFileRunRequest) (runner.SingleFileRunResult, *appErrors.Error) {
	model.Sanitize()

	if model.codeBlock.Emulator.Name == "node_latest" {
		res := execution.PackageService.RunJob(execution.Job{
			BuilderType:       "single_file",
			ExecutionType:     "single_file",
			EmulatorName:      string(model.codeBlock.Emulator.Name),
			EmulatorExtension: model.codeBlock.Emulator.Extension,
			EmulatorTag:       string(model.codeBlock.Emulator.Tag),
			EmulatorText:      model.Text,
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

	if model.codeBlock.Emulator.Name == "go" {
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

	if model.codeBlock.Emulator.Name == "ruby" {
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

	if model.codeBlock.Emulator.Name == "php74" {
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

	if model.codeBlock.Emulator.Name == "python2" {
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

	if model.codeBlock.Emulator.Name == "python3" {
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

	if model.codeBlock.Emulator.Name == "c_sharp_mono" {
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

	if model.codeBlock.Emulator.Name == "c++" {
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

	if model.codeBlock.Emulator.Name == "haskell" {
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

	if model.codeBlock.Emulator.Name == "c" {
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

	if model.codeBlock.Emulator.Name == "rust" {
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

	builder := builders.CreateBuilder("single_file").(builders.SingleFileRunFn)

	buildResult, err := builder(model.codeBlock, "single_file")

	if err != nil {
		return runner.SingleFileRunResult{}, err
	}

	builtRunner := runner.CreateRunner("singleFile").(runner.SingleFileRunFn)

	containerName := uuid.New().String()

	runResult, err := builtRunner(runner.SingleFileBuildResult{
		ContainerName:      containerName,
		DirectoryName:      buildResult.DirectoryName,
		ExecutionDirectory: buildResult.ExecutionDirectory,
		FileName:           buildResult.FileName,
		Environment:        model.codeBlock.Emulator,
		StateDirectory:     buildResult.StateDirectory,
		Args:               createCommand(buildResult, model.codeBlock.Emulator, containerName),
		Timeout:            model.validatedTemporarySession.Timeout,
	})

	if err != nil {
		return runner.SingleFileRunResult{}, err
	}

	destroyRunner := builders.CreateDestroyer("single_file").(builders.SingleFileDestroyFn)

	if err := destroyRunner(buildResult); err != nil {
		// log here if it fails, do not tell the user
		return runResult, nil
	}

	return runResult, nil
}
