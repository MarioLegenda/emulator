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
