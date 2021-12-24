package projectExecution

import (
	"github.com/google/uuid"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/builders"
	"therebelsource/emulator/runner"
)

var ProjectExecutionService Service

type Service struct {}

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

func (s Service) RunProject(model *ProjectRunRequest) (runner.ProjectRunResult, *appErrors.Error) {
	if model.sessionData.CodeProject.Environment.Name == "c" {
		projectBuilder := builders.CreateBuilder("c_project").(builders.CProjectBuildFn)

		containerName := uuid.New().String()

		buildResult, err := projectBuilder(model.sessionData.CodeProject, model.sessionData.Content, "session", model.executingFile)
		defer destroy(buildResult.ExecutionDirectory)

		args := createCommand(buildResult, model.sessionData.CodeProject.Environment, containerName)

		buildResult.Args = args

		if err != nil {
			return runner.ProjectRunResult{}, err
		}

		builtRunner := runner.CreateRunner("singleFile").(runner.SingleFileRunFn)

		runResult, err := builtRunner(runner.SingleFileBuildResult{
			ContainerName: containerName,
			ExecutionDirectory: buildResult.ExecutionDirectory,
			Environment:     model.sessionData.CodeProject.Environment,
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

	if model.sessionData.CodeProject.Environment.Name == "c++" {
		projectBuilder := builders.CreateBuilder("c_project").(builders.CProjectBuildFn)

		containerName := uuid.New().String()

		buildResult, err := projectBuilder(model.sessionData.CodeProject, model.sessionData.Content, "session", model.executingFile)
		defer destroy(buildResult.ExecutionDirectory)

		args := createCommand(buildResult, model.sessionData.CodeProject.Environment, containerName)

		buildResult.Args = args

		if err != nil {
			return runner.ProjectRunResult{}, err
		}

		builtRunner := runner.CreateRunner("singleFile").(runner.SingleFileRunFn)

		runResult, err := builtRunner(runner.SingleFileBuildResult{
			ContainerName: containerName,
			ExecutionDirectory: buildResult.ExecutionDirectory,
			Environment:     model.sessionData.CodeProject.Environment,
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

	if model.sessionData.CodeProject.Environment.Name == "haskell" {
		projectBuilder := builders.CreateBuilder("c_project").(builders.CProjectBuildFn)

		containerName := uuid.New().String()

		buildResult, err := projectBuilder(model.sessionData.CodeProject, model.sessionData.Content, "session", model.executingFile)
		defer destroy(buildResult.ExecutionDirectory)

		args := createCommand(buildResult, model.sessionData.CodeProject.Environment, containerName)

		buildResult.Args = args

		if err != nil {
			return runner.ProjectRunResult{}, err
		}

		builtRunner := runner.CreateRunner("singleFile").(runner.SingleFileRunFn)

		runResult, err := builtRunner(runner.SingleFileBuildResult{
			ContainerName: containerName,
			ExecutionDirectory: buildResult.ExecutionDirectory,
			Environment:     model.sessionData.CodeProject.Environment,
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

	buildResult, err := projectBuilder(model.sessionData.CodeProject, model.sessionData.Content, "session", model.executingFile)
	defer destroy(buildResult.ExecutionDirectory)

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
		Environment:     model.sessionData.CodeProject.Environment,
		StateDirectory: buildResult.StateDirectory,
		Args: createCommand(buildResult, model.sessionData.CodeProject.Environment, containerName),
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
