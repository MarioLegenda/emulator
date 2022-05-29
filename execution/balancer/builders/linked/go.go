package linked

import (
	"fmt"
	"github.com/google/uuid"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/repository"
)

type GoProjectBuildResult struct {
	ContainerDirectory string
	ExecutionDirectory string
}

type GoProjectBuildParams struct {
	CodeProject        *repository.CodeProject
	Contents           []*repository.FileContent
	ContainerDirectory string
	Text               string
	PackageName        string
}

func InitGoParams(cp *repository.CodeProject, contents []*repository.FileContent, containerDir string, text string, packageName string) GoProjectBuildParams {
	return GoProjectBuildParams{
		CodeProject:        cp,
		Contents:           contents,
		ContainerDirectory: containerDir,
		Text:               text,
		PackageName:        packageName,
	}
}

func GoProjectBuild(params GoProjectBuildParams) (GoProjectBuildResult, *appErrors.Error) {
	execDirConstant := uuid.New().String()
	executionDir := fmt.Sprintf("%s/%s", params.ContainerDirectory, params.PackageName)
	ft := initFileTraverse(params.CodeProject.Structure, executionDir)

	paths := ft.createPaths()

	if err := createDir(executionDir); err != nil {
		return GoProjectBuildResult{}, err
	}

	if err := createFsSystem(paths, params.Contents); err != nil {
		return GoProjectBuildResult{}, nil
	}

	fileName := fmt.Sprintf("%s.%s", execDirConstant, params.CodeProject.Environment.Extension)
	if err := writeContent(fileName, executionDir, params.Text); err != nil {
		return GoProjectBuildResult{}, err
	}

	return GoProjectBuildResult{
		ContainerDirectory: params.PackageName,
		ExecutionDirectory: executionDir,
	}, nil
}
