package projectExecution

import (
	"github.com/google/uuid"
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

func createCommand(params interface{}, lang *runner.Language, containerName string) []string {
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

	return commandFactory.CreateCommand(containerName, br.ExecutionDirectory, br.FileName, lang, br.DirectoryName)
}

func (s Service) RunProject(model *CodeProjectRunRequest) (runner.ProjectRunResult, *appErrors.Error) {
	repository := repository.InitRepository()

	contents, err := repository.GetAllFileContent(model.CodeProjectUuid)

	if err != nil {
		return runner.ProjectRunResult{}, err
	}

	if model.codeProject.Environment.Name == "c" {
		projectBuilder := builders.CreateBuilder("c_project").(builders.CProjectBuildFn)

		containerName := uuid.New().String()

		buildResult, err := projectBuilder(model.codeProject, contents, "session", model.executingFile)

		args := createCommand(buildResult, model.codeProject.Environment, containerName)

		buildResult.Args = args

		if err != nil {
			return runner.ProjectRunResult{}, err
		}

		builtRunner := runner.CreateRunner("singleFile").(runner.SingleFileRunFn)

		runResult, err := builtRunner(runner.SingleFileBuildResult{
			ContainerName: containerName,
			ExecutionDirectory: buildResult.ExecutionDirectory,
			Environment:     model.codeProject.Environment,
			Args: args,
		})

		if err != nil {
			return runner.ProjectRunResult{}, err
		}

		return runner.ProjectRunResult{
			Success: runResult.Success,
			Result:  runResult.Result,
			Timeout: runResult.Timeout,
		}, nil
	}

	if model.codeProject.Environment.Name == "c++" {
		projectBuilder := builders.CreateBuilder("c_project").(builders.CProjectBuildFn)

		containerName := uuid.New().String()

		buildResult, err := projectBuilder(model.codeProject, contents, "session", model.executingFile)

		args := createCommand(buildResult, model.codeProject.Environment, containerName)

		buildResult.Args = args

		if err != nil {
			return runner.ProjectRunResult{}, err
		}

		builtRunner := runner.CreateRunner("singleFile").(runner.SingleFileRunFn)

		runResult, err := builtRunner(runner.SingleFileBuildResult{
			ContainerName: containerName,
			ExecutionDirectory: buildResult.ExecutionDirectory,
			Environment:     model.codeProject.Environment,
			Args: args,
		})

		if err != nil {
			return runner.ProjectRunResult{}, err
		}

		return runner.ProjectRunResult{
			Success: runResult.Success,
			Result:  runResult.Result,
			Timeout: runResult.Timeout,
		}, nil
	}

	if model.codeProject.Environment.Name == "haskell" {
		projectBuilder := builders.CreateBuilder("c_project").(builders.CProjectBuildFn)

		containerName := uuid.New().String()

		buildResult, err := projectBuilder(model.codeProject, contents, "session", model.executingFile)

		args := createCommand(buildResult, model.codeProject.Environment, containerName)

		buildResult.Args = args

		if err != nil {
			return runner.ProjectRunResult{}, err
		}

		builtRunner := runner.CreateRunner("singleFile").(runner.SingleFileRunFn)

		runResult, err := builtRunner(runner.SingleFileBuildResult{
			ContainerName: containerName,
			ExecutionDirectory: buildResult.ExecutionDirectory,
			Environment:     model.codeProject.Environment,
			Args: args,
		})

		if err != nil {
			return runner.ProjectRunResult{}, err
		}

		return runner.ProjectRunResult{
			Success: runResult.Success,
			Result:  runResult.Result,
			Timeout: runResult.Timeout,
		}, nil
	}

	projectBuilder := builders.CreateBuilder("project").(builders.ProjectBuildFn)

	buildResult, err := projectBuilder(model.codeProject, contents, "session", model.executingFile)

	if err != nil {
		return runner.ProjectRunResult{}, err
	}

	builtRunner := runner.CreateRunner("singleFile").(runner.SingleFileRunFn)

	containerName := uuid.New().String()

	runResult, err := builtRunner(runner.SingleFileBuildResult{
		ContainerName: containerName,
		DirectoryName:     buildResult.DirectoryName,
		ExecutionDirectory: buildResult.ExecutionDirectory,
		FileName:           buildResult.FileName,
		Environment:     model.codeProject.Environment,
		StateDirectory: buildResult.StateDirectory,
		Args: createCommand(buildResult, model.codeProject.Environment, containerName),
	})

	if err != nil {
		return runner.ProjectRunResult{}, err
	}

	return runner.ProjectRunResult{
		Success: runResult.Success,
		Result:  runResult.Result,
		Timeout: runResult.Timeout,
	}, nil
}
