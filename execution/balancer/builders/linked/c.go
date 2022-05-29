package linked

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/repository"
)

type CProjectBuildResult struct {
	BinaryFileName     string
	ResolvedFiles      string
	ExecutionDirectory string
	ContainerDirectory string
}

type CProjectBuildParams struct {
	CodeProject        *repository.CodeProject
	Contents           []*repository.FileContent
	ContainerDirectory string
	Text               string
}

func InitCParams(cp *repository.CodeProject, contents []*repository.FileContent, containerDir string, text string) CProjectBuildParams {
	return CProjectBuildParams{
		CodeProject:        cp,
		Contents:           contents,
		ContainerDirectory: containerDir,
		Text:               text,
	}
}

func CProjectBuild(params CProjectBuildParams) (CProjectBuildResult, *appErrors.Error) {
	execDirConstant := uuid.New().String()
	executionDir := fmt.Sprintf("%s/%s", params.ContainerDirectory, execDirConstant)
	ft := initFileTraverse(params.CodeProject.Structure, executionDir)

	paths := ft.createPaths()

	if err := createDir(executionDir); err != nil {
		return CProjectBuildResult{}, err
	}

	if err := createFsSystem(paths, params.Contents); err != nil {
		return CProjectBuildResult{}, nil
	}

	fileName := fmt.Sprintf("%s.%s", execDirConstant, params.CodeProject.Environment.Extension)
	if err := writeContent(fileName, executionDir, params.Text); err != nil {
		return CProjectBuildResult{}, err
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
		BinaryFileName:     params.CodeProject.Uuid,
		ResolvedFiles:      resolvedFiles,
		ExecutionDirectory: executionDir,
		ContainerDirectory: execDirConstant,
	}, nil
}
