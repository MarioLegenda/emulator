package builders

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"therebelsource/emulator/appErrors"
)

type NodeSingleFileBuildResult struct {
	Directory string
	FileName  string
}

type NodeSingleFileBuildParams struct {
	Extension string
	Text      string
	StateDir  string
}

func InitParams(ext string, text string, stateDir string) NodeSingleFileBuildParams {
	return NodeSingleFileBuildParams{
		Extension: ext,
		Text:      text,
		StateDir:  stateDir,
	}
}

func NodeSingleFileBuild(params NodeSingleFileBuildParams) (NodeSingleFileBuildResult, *appErrors.Error) {
	dirName := uuid.New().String()
	tempExecutionDir := fmt.Sprintf("%s/%s", params.StateDir, dirName)
	fileName := fmt.Sprintf("%s.%s", dirName, params.Extension)

	if err := os.MkdirAll(tempExecutionDir, os.ModePerm); err != nil {
		return NodeSingleFileBuildResult{}, appErrors.New(appErrors.ApplicationError, appErrors.FilesystemError, fmt.Sprintf("Cannot create execution dir: %s", err.Error()))
	}

	if err := writeContent(fileName, tempExecutionDir, params.Text); err != nil {
		return NodeSingleFileBuildResult{}, err
	}

	if err := writeContent("output.txt", tempExecutionDir, ""); err != nil {
		return NodeSingleFileBuildResult{}, err
	}

	return NodeSingleFileBuildResult{
		Directory: tempExecutionDir,
		FileName:  fileName,
	}, nil
}
