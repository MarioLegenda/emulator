package builders

import (
	"fmt"
	"strings"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/repository"
)

type ProjectBuildFn func(*repository.CodeProject, []*repository.FileContent, string, *repository.File) (ProjectBuildResult, *appErrors.Error)
type CProjectBuildFn func(*repository.CodeProject, []*repository.FileContent, string, *repository.File) (CProjectBuildResult, *appErrors.Error)

type ProjectBuildResult struct {
	DirectoryName string
	StateDirectory string
	ExecutionDirectory string
	FileName  string
	Args []string
}

type CProjectBuildResult struct {
	BinaryFileName string
	ResolvedFiles string
	ExecutionDirectory string
	Args []string
}

func createProjectBuilder() ProjectBuildFn {
	return func(cb *repository.CodeProject, contents []*repository.FileContent, state string, executingFile *repository.File) (ProjectBuildResult, *appErrors.Error) {
		executionDir := fmt.Sprintf("%s/%s", getStateDirectory(state), cb.Uuid)
		ft := initFileTraverse(cb.Structure, executionDir)

		paths := ft.createPaths()

		if err := createDir(fmt.Sprintf("%s/%s", getStateDirectory(state), cb.Uuid)); err != nil {
			return ProjectBuildResult{}, err
		}

		if err := createFsSystem(paths, contents); err != nil {
			return ProjectBuildResult{}, nil
		}

		if cb.Environment.Name == "rust" {
			if err := writeContent("Cargo.toml", executionDir, `
[package]
name = "All executions"
version = "0.0.1"
authors = [ "No name" ]

[[bin]]
name = "main"
path = "main.rs"
`); err != nil {
				return ProjectBuildResult{}, err
			}
		}

		return ProjectBuildResult{
			DirectoryName: cb.Uuid,
			StateDirectory: getStateDirectory(state),
			ExecutionDirectory: executionDir,
			FileName: executingFile.Name,
		}, nil
	}
}

func createCLangBuilder() CProjectBuildFn {
	return func(cb *repository.CodeProject, contents []*repository.FileContent, state string, executingFile *repository.File) (CProjectBuildResult, *appErrors.Error) {
		executionDir := fmt.Sprintf("%s/%s", getStateDirectory(state), cb.Uuid)
		ft := initFileTraverse(cb.Structure, executionDir)

		paths := ft.createPaths()

		if err := createDir(fmt.Sprintf("%s/%s", getStateDirectory(state), cb.Uuid)); err != nil {
			return CProjectBuildResult{}, err
		}

		if err := createFsSystem(paths, contents); err != nil {
			return CProjectBuildResult{}, nil
		}

		resolvedFiles := ""
		for dir, files := range paths {
			s := strings.Split(dir, cb.Uuid)
			dockerPath := s[1]

			for _, file := range files {
				if dockerPath == "" {
					resolvedFiles += fmt.Sprintf("%s ", file.Name)
				} else {
					resolvedFiles += fmt.Sprintf("%s/%s ", dockerPath, file.Name)
				}
			}
		}

		return CProjectBuildResult{
			BinaryFileName: cb.Uuid,
			ResolvedFiles: resolvedFiles,
			ExecutionDirectory: executionDir,
		}, nil
	}
}