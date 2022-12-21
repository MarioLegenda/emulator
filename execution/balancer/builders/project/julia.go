package project

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/repository"
)

type JuliaProjectBuildResult struct {
	ContainerDirectory string
	ExecutionDirectory string
	FileName           string
}

type JuliaProjectBuildParams struct {
	CodeProject        *repository.CodeProject
	Contents           []*repository.FileContent
	ContainerDirectory string
	ExecutingFile      *repository.File
}

func InitJuliaParams(cp *repository.CodeProject, contents []*repository.FileContent, containerDir string, executingFile *repository.File) JuliaProjectBuildParams {
	return JuliaProjectBuildParams{
		CodeProject:        cp,
		Contents:           contents,
		ContainerDirectory: containerDir,
		ExecutingFile:      executingFile,
	}
}

func JuliaProjectBuild(params JuliaProjectBuildParams) (JuliaProjectBuildResult, *appErrors.Error) {
	execDirConstant := uuid.New().String()

	executionDir := fmt.Sprintf("%s/%s", params.ContainerDirectory, execDirConstant)
	ft := initFileTraverse(params.CodeProject.Structure, executionDir)

	paths := ft.createPaths()

	if err := createDir(executionDir); err != nil {
		return JuliaProjectBuildResult{}, err
	}

	if err := createFsSystem(paths, params.Contents); err != nil {
		return JuliaProjectBuildResult{}, nil
	}

	fileName := params.ExecutingFile.Name

	if params.ExecutingFile.Depth != 1 {
		for path, files := range paths {
			for _, file := range files {
				if file.Uuid == params.ExecutingFile.Uuid {
					s := strings.Split(path, execDirConstant)

					fileName = fmt.Sprintf("/app%s/%s", s[1], params.ExecutingFile.Name)
				}
			}
		}
	}

	return JuliaProjectBuildResult{
		ContainerDirectory: fmt.Sprintf("/app/%s", execDirConstant),
		ExecutionDirectory: executionDir,
		FileName:           fileName,
	}, nil
}