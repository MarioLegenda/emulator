package builders

import (
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/repository"
)

type ProjectBuildFn func(*repository.CodeProject, []*repository.FileContent, string) (ProjectBuildResult, *appErrors.Error)

type ProjectBuildResult struct {
	DirectoryName string
	StateDirectory string
	ExecutionDirectory string
	FileName  string
}

func createProjectBuilder() ProjectBuildFn {
	return func(cb *repository.CodeProject, contents []*repository.FileContent, state string) (ProjectBuildResult, *appErrors.Error) {
		return ProjectBuildResult{}, nil
	}
}