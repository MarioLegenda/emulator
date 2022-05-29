package linked

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/repository"
)

type CPlusProjectBuildResult struct {
	BinaryFileName     string
	ResolvedFiles      string
	ExecutionDirectory string
	ContainerDirectory string
}

type CPlusProjectBuildParams struct {
	CodeProject        *repository.CodeProject
	Contents           []*repository.FileContent
	ContainerDirectory string
	Text               string
}

func InitCPlusParams(cp *repository.CodeProject, contents []*repository.FileContent, containerDir string, text string) CPlusProjectBuildParams {
	return CPlusProjectBuildParams{
		CodeProject:        cp,
		Contents:           contents,
		ContainerDirectory: containerDir,
		Text:               text,
	}
}

func CPlusProjectBuild(params CPlusProjectBuildParams) (CPlusProjectBuildResult, *appErrors.Error) {
	execDirConstant := uuid.New().String()
	executionDir := fmt.Sprintf("%s/%s", params.ContainerDirectory, execDirConstant)
	ft := initFileTraverse(params.CodeProject.Structure, executionDir)

	paths := ft.createPaths()

	if err := createDir(executionDir); err != nil {
		return CPlusProjectBuildResult{}, err
	}

	if err := createFsSystem(paths, params.Contents); err != nil {
		return CPlusProjectBuildResult{}, nil
	}

	fileName := fmt.Sprintf("%s.%s", execDirConstant, params.CodeProject.Environment.Extension)
	if err := writeContent(fileName, executionDir, params.Text); err != nil {
		return CPlusProjectBuildResult{}, err
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

	return CPlusProjectBuildResult{
		BinaryFileName:     params.CodeProject.Uuid,
		ResolvedFiles:      resolvedFiles,
		ExecutionDirectory: executionDir,
		ContainerDirectory: execDirConstant,
	}, nil
}
