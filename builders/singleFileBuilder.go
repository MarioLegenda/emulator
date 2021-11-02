package builders

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/repository"
)

type SingleFileRunFn func(cb *repository.CodeBlock, state string) (SingleFileBuildResult, *appErrors.Error)
type SingleFileDestroyFn func(br SingleFileBuildResult) *appErrors.Error

type SingleFileBuildResult struct {
	DirectoryName string
	StateDirectory string
	ExecutionDirectory string
	FileName  string
	Args []string
}

func createSingleFileBuilder() SingleFileRunFn {
	return func(cb *repository.CodeBlock, state string) (SingleFileBuildResult, *appErrors.Error) {
		stateDir := getStateDirectory(state)
		dirName := uuid.New().String()

		tempExecutionDir := fmt.Sprintf("%s/%s", stateDir, dirName)

		fileName := fmt.Sprintf("%s.%s", dirName, cb.Emulator.Extension)

		if err := os.MkdirAll(tempExecutionDir, os.ModePerm); err != nil {
			return SingleFileBuildResult{}, appErrors.New(appErrors.ApplicationError, appErrors.FilesystemError, fmt.Sprintf("Cannot create execution dir: %s", err.Error()))
		}

		if err := writeContent(fileName, tempExecutionDir, cb.Text); err != nil {
			return SingleFileBuildResult{}, err
		}

		if err := writeContent("output.txt", tempExecutionDir, ""); err != nil {
			return SingleFileBuildResult{}, err
		}

		return SingleFileBuildResult{
			DirectoryName: dirName,
			ExecutionDirectory: tempExecutionDir,
			StateDirectory: getStateDirectory(state),
			FileName:  fileName,
		}, nil
	}
}

func createSingleFileDestroyer() SingleFileDestroyFn {
	return func(br SingleFileBuildResult) *appErrors.Error {
		if err := os.RemoveAll(br.ExecutionDirectory); err != nil {
			return appErrors.New(appErrors.ApplicationError, appErrors.FilesystemError, fmt.Sprintf("Cannot remove single_file directory: %s", br.ExecutionDirectory))
		}

		return nil
	}
}
