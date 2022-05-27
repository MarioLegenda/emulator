package projectExecution

import (
	"github.com/google/uuid"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/builders"
	"therebelsource/emulator/execution"
	"therebelsource/emulator/repository"
	"therebelsource/emulator/runner"
	_var "therebelsource/emulator/var"
)

var ProjectExecutionService Service

type Service struct{}

func InitService() {
	ProjectExecutionService = Service{}
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
		br := params.(builders.CProjectBuildResult)

		return commandFactory.CreateCProjectCommand(uuid.New().String(), br.BinaryFileName, br.ExecutionDirectory, br.ResolvedFiles, lang)
	}

	if lang.Name == "c++" {
		br := params.(builders.CProjectBuildResult)

		return commandFactory.CreateCPlusProjectCommand(uuid.New().String(), br.BinaryFileName, br.ExecutionDirectory, br.ResolvedFiles, lang)
	}

	if lang.Name == "haskell" {
		br := params.(builders.CProjectBuildResult)

		return commandFactory.CreateHaskellProjectCommand(uuid.New().String(), br.ExecutionDirectory, lang)
	}

	br := params.(builders.ProjectBuildResult)

	cmd := commandFactory.CreateCommand(containerName, br.ExecutionDirectory, br.FileName, lang, br.DirectoryName)

	return cmd
}

func (s Service) RunProject(model *ProjectRunRequest) (repository.ProjectRunResult, *appErrors.Error) {
	if model.sessionData.CodeProject.Environment.Name == "node_latest" {
		res := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:   "project",
			ExecutionType: "project",
			EmulatorTag:   string(model.sessionData.CodeProject.Environment.Tag),
			EmulatorName:  string(model.sessionData.CodeProject.Environment.Name),
			CodeProject:   model.sessionData.CodeProject,
			ExecutingFile: model.sessionData.ExecutingFile,
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
			BuilderType:   "project",
			ExecutionType: "project",
			EmulatorTag:   string(model.sessionData.CodeProject.Environment.Tag),
			EmulatorName:  string(model.sessionData.CodeProject.Environment.Name),
			CodeProject:   model.sessionData.CodeProject,
			ExecutingFile: model.sessionData.ExecutingFile,
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

	if model.sessionData.CodeProject.Environment.Name == "ruby" {
		res := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:   "project",
			ExecutionType: "project",
			EmulatorTag:   string(model.sessionData.CodeProject.Environment.Tag),
			EmulatorName:  string(model.sessionData.CodeProject.Environment.Name),
			CodeProject:   model.sessionData.CodeProject,
			ExecutingFile: model.sessionData.ExecutingFile,
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

	if model.sessionData.CodeProject.Environment.Name == "c" {
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

	if model.sessionData.CodeProject.Environment.Name == "c++" {
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

	if model.sessionData.CodeProject.Environment.Name == "c_sharp_mono" {
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

	if model.sessionData.CodeProject.Environment.Name == "python2" {
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

	if model.sessionData.CodeProject.Environment.Name == "python3" {
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

	if model.sessionData.CodeProject.Environment.Name == "haskell" {
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

	if model.sessionData.CodeProject.Environment.Name == "c" {
		projectBuilder := builders.CreateBuilder("c_project").(builders.CProjectBuildFn)

		containerName := uuid.New().String()

		buildResult, err := projectBuilder(model.sessionData.CodeProject, model.sessionData.Content, builders.PROJECT_EXECUTION_STATE, nil)
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
		projectBuilder := builders.CreateBuilder("c_project").(builders.CProjectBuildFn)

		containerName := uuid.New().String()

		buildResult, err := projectBuilder(model.sessionData.CodeProject, model.sessionData.Content, builders.PROJECT_EXECUTION_STATE, nil)
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
		projectBuilder := builders.CreateBuilder("c_project").(builders.CProjectBuildFn)

		containerName := uuid.New().String()

		buildResult, err := projectBuilder(model.sessionData.CodeProject, model.sessionData.Content, builders.PROJECT_EXECUTION_STATE, nil)
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

	projectBuilder := builders.CreateBuilder("project").(builders.ProjectBuildFn)

	executingDir := uuid.New().String()
	if model.sessionData.CodeProject.Environment.Name == "go" {
		executingDir = model.sessionData.PackageName
	}

	buildResult, err := projectBuilder(
		model.sessionData.CodeProject,
		model.sessionData.Content,
		builders.PROJECT_EXECUTION_STATE,
		model.sessionData.ExecutingFile,
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
