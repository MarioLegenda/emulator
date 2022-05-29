package linkedProjectExecution

import (
	"github.com/google/uuid"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/builders"
	"therebelsource/emulator/execution"
	"therebelsource/emulator/repository"
	"therebelsource/emulator/runner"
	_var "therebelsource/emulator/var"
)

var ExecutionService Service

type Service struct{}

func InitService() {
	ExecutionService = Service{}
}

func destroy(dir string) *appErrors.Error {
	destroyRunner := builders.CreateDestroyer("project").(builders.ProjectDestroyFn)

	if err := destroyRunner(dir); err != nil {
		// log here if it fails, do not tell the user
		return nil
	}

	return nil
}

func goDestroy(dir string) {
	go func(executionDirectory string) {
		if err := destroy(executionDirectory); err != nil {
			// TODO: log and send to slack, big error
		}
	}(dir)
}

func createCommand(params interface{}, lang *repository.Language, containerName string) []string {
	commandFactory := runner.RunCommandFactory{}

	if lang.Name == "c" {
		br := params.(builders.LinkedProjectBuildResult)

		return commandFactory.CreateCProjectCommand(uuid.New().String(), br.BinaryFileName, br.ExecutionDirectory, br.ResolvedFiles, lang)
	}

	if lang.Name == "c++" {
		br := params.(builders.LinkedProjectBuildResult)

		return commandFactory.CreateCPlusProjectCommand(uuid.New().String(), br.BinaryFileName, br.ExecutionDirectory, br.ResolvedFiles, lang)
	}

	if lang.Name == "haskell" {
		br := params.(builders.LinkedProjectBuildResult)

		return commandFactory.CreateHaskellProjectCommand(uuid.New().String(), br.ExecutionDirectory, lang)
	}

	if lang.Name == "go" {
		br := params.(builders.ProjectBuildResult)

		return commandFactory.CreateGoLinkedProjectCommand(containerName, br.ExecutionDirectory, lang, br.DirectoryName)
	}

	br := params.(builders.ProjectBuildResult)

	return commandFactory.CreateCommand(containerName, br.ExecutionDirectory, br.FileName, lang, br.DirectoryName)
}

func (s Service) RunProject(model *LinkedProjectRunRequest) (repository.ProjectRunResult, *appErrors.Error) {
	if model.sessionData.CodeProject.Environment.Name == "ruby" {
		res := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:   "linked",
			ExecutionType: "linked",
			EmulatorTag:   string(model.sessionData.CodeProject.Environment.Tag),
			EmulatorName:  string(model.sessionData.CodeProject.Environment.Name),
			EmulatorText:  model.sessionData.CodeBlock.Text,
			CodeProject:   model.sessionData.CodeProject,
			Contents:      model.sessionData.Content,
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

	if model.sessionData.CodeProject.Environment.Name == "node_latest" {
		res := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:   "linked",
			ExecutionType: "linked",
			EmulatorTag:   string(model.sessionData.CodeProject.Environment.Tag),
			EmulatorName:  string(model.sessionData.CodeProject.Environment.Name),
			EmulatorText:  model.sessionData.CodeBlock.Text,
			CodeProject:   model.sessionData.CodeProject,
			Contents:      model.sessionData.Content,
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

	if model.sessionData.CodeProject.Environment.Name == "node_latest_esm" {
		res := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:   "linked",
			ExecutionType: "linked",
			EmulatorTag:   string(model.sessionData.CodeProject.Environment.Tag),
			EmulatorName:  string(model.sessionData.CodeProject.Environment.Name),
			EmulatorText:  model.sessionData.CodeBlock.Text,
			CodeProject:   model.sessionData.CodeProject,
			Contents:      model.sessionData.Content,
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

	if model.sessionData.CodeProject.Environment.Name == "go" {
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

	if model.sessionData.CodeProject.Environment.Name == "haskell" {
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

	if model.sessionData.CodeProject.Environment.Name == "c" {
		projectBuilder := builders.CreateBuilder("linked_compiled_project").(builders.LinkedBuildFn)

		containerName := uuid.New().String()

		buildResult, err := projectBuilder(model.sessionData.CodeProject, model.sessionData.Content, builders.PROJECT_EXECUTION_STATE, model.sessionData.CodeBlock)
		defer goDestroy(buildResult.ExecutionDirectory)

		args := createCommand(buildResult, model.sessionData.CodeProject.Environment, containerName)

		buildResult.Args = args

		if err != nil {
			return repository.ProjectRunResult{}, err
		}

		builtRunner := runner.CreateRunner("singleFile").(runner.SingleFileRunFn)

		runResult, err := builtRunner(repository.SingleFileBuildResult{
			ContainerName:      containerName,
			ExecutionDirectory: buildResult.ExecutionDirectory,
			Environment:        model.sessionData.CodeProject.Environment,
			Args:               args,
			Timeout:            model.validatedTemporarySession.Timeout,
		})

		if err != nil {
			return repository.ProjectRunResult{}, err
		}

		return repository.ProjectRunResult{
			Success: runResult.Success,
			Result:  runResult.Result,
			Timeout: runResult.Timeout,
		}, nil
	}

	if model.sessionData.CodeProject.Environment.Name == "c++" {
		projectBuilder := builders.CreateBuilder("linked_compiled_project").(builders.LinkedBuildFn)

		containerName := uuid.New().String()

		buildResult, err := projectBuilder(model.sessionData.CodeProject, model.sessionData.Content, builders.PROJECT_EXECUTION_STATE, model.sessionData.CodeBlock)
		defer goDestroy(buildResult.ExecutionDirectory)

		args := createCommand(buildResult, model.sessionData.CodeProject.Environment, containerName)

		buildResult.Args = args

		if err != nil {
			return repository.ProjectRunResult{}, err
		}

		builtRunner := runner.CreateRunner("singleFile").(runner.SingleFileRunFn)

		runResult, err := builtRunner(repository.SingleFileBuildResult{
			ContainerName:      containerName,
			ExecutionDirectory: buildResult.ExecutionDirectory,
			Environment:        model.sessionData.CodeProject.Environment,
			Args:               args,
			Timeout:            model.validatedTemporarySession.Timeout,
		})

		if err != nil {
			return repository.ProjectRunResult{}, err
		}

		return repository.ProjectRunResult{
			Success: runResult.Success,
			Result:  runResult.Result,
			Timeout: runResult.Timeout,
		}, nil
	}

	if model.sessionData.CodeProject.Environment.Name == "haskell" {
		projectBuilder := builders.CreateBuilder("linked_compiled_project").(builders.LinkedBuildFn)

		containerName := uuid.New().String()

		buildResult, err := projectBuilder(model.sessionData.CodeProject, model.sessionData.Content, builders.PROJECT_EXECUTION_STATE, model.sessionData.CodeBlock)
		defer goDestroy(buildResult.ExecutionDirectory)

		args := createCommand(buildResult, model.sessionData.CodeProject.Environment, containerName)

		buildResult.Args = args

		if err != nil {
			return repository.ProjectRunResult{}, err
		}

		builtRunner := runner.CreateRunner("singleFile").(runner.SingleFileRunFn)

		runResult, err := builtRunner(repository.SingleFileBuildResult{
			ContainerName:      containerName,
			ExecutionDirectory: buildResult.ExecutionDirectory,
			Environment:        model.sessionData.CodeProject.Environment,
			Args:               args,
			Timeout:            model.validatedTemporarySession.Timeout,
		})

		if err != nil {
			return repository.ProjectRunResult{}, err
		}

		return repository.ProjectRunResult{
			Success: runResult.Success,
			Result:  runResult.Result,
			Timeout: runResult.Timeout,
		}, nil
	}

	projectBuilder := builders.CreateBuilder("linked_interpreted_project").(builders.LinkedInterpretedBuildFn)

	executingDir := uuid.New().String()
	if model.sessionData.CodeProject.Environment.Name == "go" {
		executingDir = model.sessionData.CodeBlock.PackageName
	}

	buildResult, err := projectBuilder(
		model.sessionData.CodeProject,
		model.sessionData.Content,
		builders.PROJECT_EXECUTION_STATE,
		model.sessionData.CodeBlock,
		executingDir,
	)

	defer goDestroy(buildResult.ExecutionDirectory)

	if err != nil {
		return repository.ProjectRunResult{}, err
	}

	builtRunner := runner.CreateRunner("singleFile").(runner.SingleFileRunFn)

	containerName := uuid.New().String()

	runResult, err := builtRunner(repository.SingleFileBuildResult{
		ContainerName:      containerName,
		DirectoryName:      buildResult.DirectoryName,
		ExecutionDirectory: buildResult.ExecutionDirectory,
		FileName:           buildResult.FileName,
		Environment:        model.sessionData.CodeProject.Environment,
		StateDirectory:     buildResult.StateDirectory,
		Args:               createCommand(buildResult, model.sessionData.CodeProject.Environment, containerName),
		Timeout:            model.validatedTemporarySession.Timeout,
	})

	if err != nil {
		return repository.ProjectRunResult{}, err
	}

	return repository.ProjectRunResult{
		Success: runResult.Success,
		Result:  runResult.Result,
		Timeout: runResult.Timeout,
	}, nil
}

func (s Service) RunPublicProject(model *PublicLinkedProjectRunRequest) (repository.ProjectRunResult, *appErrors.Error) {
	if model.sessionData.CodeProject.Environment.Name == "ruby" {
		res := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:   "linked",
			ExecutionType: "linked",
			EmulatorTag:   string(model.sessionData.CodeProject.Environment.Tag),
			EmulatorName:  string(model.sessionData.CodeProject.Environment.Name),
			EmulatorText:  model.Text,
			CodeProject:   model.sessionData.CodeProject,
			Contents:      model.sessionData.Content,
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

	if model.sessionData.CodeProject.Environment.Name == "node_latest" {
		res := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:   "linked",
			ExecutionType: "linked",
			EmulatorTag:   string(model.sessionData.CodeProject.Environment.Tag),
			EmulatorName:  string(model.sessionData.CodeProject.Environment.Name),
			EmulatorText:  model.Text,
			CodeProject:   model.sessionData.CodeProject,
			Contents:      model.sessionData.Content,
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

	if model.sessionData.CodeProject.Environment.Name == "go" {
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

	if model.sessionData.CodeProject.Environment.Name == "haskell" {
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

	if model.sessionData.CodeProject.Environment.Name == "node_latest_esm" {
		res := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:   "linked",
			ExecutionType: "linked",
			EmulatorTag:   string(model.sessionData.CodeProject.Environment.Tag),
			EmulatorName:  string(model.sessionData.CodeProject.Environment.Name),
			EmulatorText:  model.Text,
			CodeProject:   model.sessionData.CodeProject,
			Contents:      model.sessionData.Content,
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

	if model.sessionData.CodeProject.Environment.Name == "c" {
		projectBuilder := builders.CreateBuilder("linked_compiled_project").(builders.LinkedBuildFn)

		containerName := uuid.New().String()

		buildResult, err := projectBuilder(model.sessionData.CodeProject, model.sessionData.Content, builders.PROJECT_EXECUTION_STATE, model.sessionData.CodeBlock)
		defer goDestroy(buildResult.ExecutionDirectory)

		args := createCommand(buildResult, model.sessionData.CodeProject.Environment, containerName)

		buildResult.Args = args

		if err != nil {
			return repository.ProjectRunResult{}, err
		}

		builtRunner := runner.CreateRunner("singleFile").(runner.SingleFileRunFn)

		runResult, err := builtRunner(repository.SingleFileBuildResult{
			ContainerName:      containerName,
			ExecutionDirectory: buildResult.ExecutionDirectory,
			Environment:        model.sessionData.CodeProject.Environment,
			Args:               args,
			Timeout:            model.validatedTemporarySession.Timeout,
		})

		if err != nil {
			return repository.ProjectRunResult{}, err
		}

		return repository.ProjectRunResult{
			Success: runResult.Success,
			Result:  runResult.Result,
			Timeout: runResult.Timeout,
		}, nil
	}

	if model.sessionData.CodeProject.Environment.Name == "c++" {
		projectBuilder := builders.CreateBuilder("linked_compiled_project").(builders.LinkedBuildFn)

		containerName := uuid.New().String()

		buildResult, err := projectBuilder(model.sessionData.CodeProject, model.sessionData.Content, builders.PROJECT_EXECUTION_STATE, model.sessionData.CodeBlock)
		defer goDestroy(buildResult.ExecutionDirectory)

		args := createCommand(buildResult, model.sessionData.CodeProject.Environment, containerName)

		buildResult.Args = args

		if err != nil {
			return repository.ProjectRunResult{}, err
		}

		builtRunner := runner.CreateRunner("singleFile").(runner.SingleFileRunFn)

		runResult, err := builtRunner(repository.SingleFileBuildResult{
			ContainerName:      containerName,
			ExecutionDirectory: buildResult.ExecutionDirectory,
			Environment:        model.sessionData.CodeProject.Environment,
			Args:               args,
			Timeout:            model.validatedTemporarySession.Timeout,
		})

		if err != nil {
			return repository.ProjectRunResult{}, err
		}

		return repository.ProjectRunResult{
			Success: runResult.Success,
			Result:  runResult.Result,
			Timeout: runResult.Timeout,
		}, nil
	}

	if model.sessionData.CodeProject.Environment.Name == "haskell" {
		projectBuilder := builders.CreateBuilder("linked_compiled_project").(builders.LinkedBuildFn)

		containerName := uuid.New().String()

		buildResult, err := projectBuilder(model.sessionData.CodeProject, model.sessionData.Content, builders.PROJECT_EXECUTION_STATE, model.sessionData.CodeBlock)
		defer goDestroy(buildResult.ExecutionDirectory)

		args := createCommand(buildResult, model.sessionData.CodeProject.Environment, containerName)

		buildResult.Args = args

		if err != nil {
			return repository.ProjectRunResult{}, err
		}

		builtRunner := runner.CreateRunner("singleFile").(runner.SingleFileRunFn)

		runResult, err := builtRunner(repository.SingleFileBuildResult{
			ContainerName:      containerName,
			ExecutionDirectory: buildResult.ExecutionDirectory,
			Environment:        model.sessionData.CodeProject.Environment,
			Args:               args,
			Timeout:            model.validatedTemporarySession.Timeout,
		})

		if err != nil {
			return repository.ProjectRunResult{}, err
		}

		return repository.ProjectRunResult{
			Success: runResult.Success,
			Result:  runResult.Result,
			Timeout: runResult.Timeout,
		}, nil
	}

	projectBuilder := builders.CreateBuilder("linked_interpreted_project").(builders.LinkedInterpretedBuildFn)

	executingDir := uuid.New().String()
	if model.sessionData.CodeProject.Environment.Name == "go" {
		executingDir = model.sessionData.CodeBlock.PackageName
	}

	buildResult, err := projectBuilder(
		model.sessionData.CodeProject,
		model.sessionData.Content,
		builders.PROJECT_EXECUTION_STATE,
		model.sessionData.CodeBlock,
		executingDir,
	)

	defer goDestroy(buildResult.ExecutionDirectory)

	if err != nil {
		return repository.ProjectRunResult{}, err
	}

	builtRunner := runner.CreateRunner("singleFile").(runner.SingleFileRunFn)

	containerName := uuid.New().String()

	runResult, err := builtRunner(repository.SingleFileBuildResult{
		ContainerName:      containerName,
		DirectoryName:      buildResult.DirectoryName,
		ExecutionDirectory: buildResult.ExecutionDirectory,
		FileName:           buildResult.FileName,
		Environment:        model.sessionData.CodeProject.Environment,
		StateDirectory:     buildResult.StateDirectory,
		Args:               createCommand(buildResult, model.sessionData.CodeProject.Environment, containerName),
		Timeout:            model.validatedTemporarySession.Timeout,
	})

	if err != nil {
		return repository.ProjectRunResult{}, err
	}

	return repository.ProjectRunResult{
		Success: runResult.Success,
		Result:  runResult.Result,
		Timeout: runResult.Timeout,
	}, nil
}
