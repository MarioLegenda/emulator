package builders

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"strings"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/repository"
)

type ProjectBuildFn func(*repository.CodeProject, []*repository.FileContent, string, *repository.File, string) (ProjectBuildResult, *appErrors.Error)
type CProjectBuildFn func(*repository.CodeProject, []*repository.FileContent, string, *repository.File) (CProjectBuildResult, *appErrors.Error)
type ProjectDestroyFn func(dir string) *appErrors.Error

type LinkedBuildFn func(*repository.CodeProject, []*repository.FileContent, string, *repository.CodeBlock) (LinkedProjectBuildResult, *appErrors.Error)
type LinkedInterpretedBuildFn func(*repository.CodeProject, []*repository.FileContent, string, *repository.CodeBlock) (ProjectBuildResult, *appErrors.Error)

type ProjectBuildResult struct {
	DirectoryName      string
	StateDirectory     string
	ExecutionDirectory string
	FileName           string
	Args               []string
}

type CProjectBuildResult struct {
	BinaryFileName     string
	ResolvedFiles      string
	ExecutionDirectory string
	Args               []string
}

type LinkedProjectBuildResult struct {
	BinaryFileName     string
	ResolvedFiles      string
	ExecutionDirectory string
	Args               []string
}

func createProjectBuilder() ProjectBuildFn {
	return func(cb *repository.CodeProject, contents []*repository.FileContent, state string, executingFile *repository.File, executingDir string) (ProjectBuildResult, *appErrors.Error) {
		execDirConstant := executingDir

		executionDir := fmt.Sprintf("%s/%s", getStateDirectory(state), execDirConstant)
		ft := initFileTraverse(cb.Structure, executionDir)

		paths := ft.createPaths()

		if err := createDir(fmt.Sprintf("%s/%s", getStateDirectory(state), cb.Uuid)); err != nil {
			return ProjectBuildResult{}, err
		}

		if err := createFsSystem(paths, contents); err != nil {
			return ProjectBuildResult{}, nil
		}

		fileName := executingFile.Name

		if executingFile.Depth != 1 {
			for path, files := range paths {
				for _, file := range files {
					if file.Uuid == executingFile.Uuid {
						s := strings.Split(path, execDirConstant)

						fileName = fmt.Sprintf("/app%s/%s", s[1], executingFile.Name)
					}
				}
			}
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

		directoryName := cb.Uuid
		if cb.Environment.Name == "go" {
			directoryName = executingDir
		}

		return ProjectBuildResult{
			DirectoryName:      directoryName,
			StateDirectory:     getStateDirectory(state),
			ExecutionDirectory: executionDir,
			FileName:           fileName,
		}, nil
	}
}

func createCLangBuilder() CProjectBuildFn {
	return func(cb *repository.CodeProject, contents []*repository.FileContent, state string, executingFile *repository.File) (CProjectBuildResult, *appErrors.Error) {
		execDirConstant := uuid.New().String()
		executionDir := fmt.Sprintf("%s/%s", getStateDirectory(state), execDirConstant)
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
			s := strings.Split(dir, execDirConstant)
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
			BinaryFileName:     cb.Uuid,
			ResolvedFiles:      resolvedFiles,
			ExecutionDirectory: executionDir,
		}, nil
	}
}

func createCompiledProject() LinkedBuildFn {
	return func(cb *repository.CodeProject, contents []*repository.FileContent, state string, executingFile *repository.CodeBlock) (LinkedProjectBuildResult, *appErrors.Error) {
		execDirConstant := uuid.New().String()
		executionDir := fmt.Sprintf("%s/%s", getStateDirectory(state), execDirConstant)
		ft := initFileTraverse(cb.Structure, executionDir)

		paths := ft.createPaths()

		if err := createDir(executionDir); err != nil {
			return LinkedProjectBuildResult{}, err
		}

		if err := createFsSystem(paths, contents); err != nil {
			return LinkedProjectBuildResult{}, nil
		}

		fileName := fmt.Sprintf("%s.%s", executingFile.Uuid, cb.Environment.Extension)
		if err := writeContent(fileName, executionDir, executingFile.Text); err != nil {
			return LinkedProjectBuildResult{}, nil
		}

		resolvedFiles := ""
		for dir, files := range paths {
			s := strings.Split(dir, execDirConstant)
			dockerPath := s[1]

			for _, file := range files {
				if dockerPath == "" {
					resolvedFiles += fmt.Sprintf("%s ", file.Name)
				} else {
					resolvedFiles += fmt.Sprintf("%s/%s ", dockerPath, file.Name)
				}
			}
		}

		resolvedFiles += fileName

		return LinkedProjectBuildResult{
			BinaryFileName:     cb.Uuid,
			ResolvedFiles:      resolvedFiles,
			ExecutionDirectory: executionDir,
		}, nil
	}
}

func createLinkedInterpretedBuildResult() LinkedInterpretedBuildFn {
	return func(cb *repository.CodeProject, contents []*repository.FileContent, state string, codeBlock *repository.CodeBlock) (ProjectBuildResult, *appErrors.Error) {
		execDirConstant := uuid.New().String()
		executionDir := fmt.Sprintf("%s/%s", getStateDirectory(state), execDirConstant)

		var paths map[string][]*repository.File
		if cb.Environment.Name == "go" {
			ft := initFileTraverse(cb.Structure, fmt.Sprintf("%s/%s", executionDir, cb.Name))

			paths = ft.createPaths()
		} else {
			ft := initFileTraverse(cb.Structure, executionDir)

			paths = ft.createPaths()
		}

		if err := createDir(executionDir); err != nil {
			return ProjectBuildResult{}, err
		}

		if err := createFsSystem(paths, contents); err != nil {
			return ProjectBuildResult{}, nil
		}

		fileName := fmt.Sprintf("%s.%s", codeBlock.Uuid, cb.Environment.Extension)
		if err := writeContent(fileName, executionDir, codeBlock.Text); err != nil {
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
			DirectoryName:      cb.Uuid,
			StateDirectory:     getStateDirectory(state),
			ExecutionDirectory: executionDir,
			FileName:           fileName,
		}, nil
	}
}

func createProjectDestroyer() ProjectDestroyFn {
	return func(dir string) *appErrors.Error {
		if err := os.RemoveAll(dir); err != nil {
			return appErrors.New(appErrors.ApplicationError, appErrors.FilesystemError, fmt.Sprintf("Cannot remove project directory: %s", dir))
		}

		return nil
	}
}
