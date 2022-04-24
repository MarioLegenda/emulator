package builders

import (
	"fmt"
	"os"
	"therebelsource/emulator/appErrors"
)

const PROJECT_EXECUTION_STATE = "project"
const SINGLE_FILE_EXECUTION_STATE = "single_file"

func writeContent(name string, dir string, content string) *appErrors.Error {
	handle, cErr := os.Create(fmt.Sprintf("%s/%s", dir, name))
	if cErr != nil {
		return appErrors.New(appErrors.ApplicationError, appErrors.FilesystemError, fmt.Sprintf("Cannot create file: %s", cErr.Error()))
	}

	_, err := handle.WriteString(content)

	if err != nil {
		if err := handle.Close(); err != nil {
			return appErrors.New(appErrors.ApplicationError, appErrors.FilesystemError, fmt.Sprintf("Cannot close a file after trying to write to it: %s", err.Error()))
		}

		return appErrors.New(appErrors.ApplicationError, appErrors.FilesystemError, fmt.Sprintf("Cannot write to file: %s", err.Error()))
	}

	err = handle.Close()
	if err != nil {
		return appErrors.New(appErrors.ApplicationError, appErrors.FilesystemError, fmt.Sprintf("Cannot close a file: %s", err.Error()))
	}

	return nil
}

func createDir(path string) *appErrors.Error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		cErr := os.MkdirAll(path, os.ModePerm)
		if cErr != nil {
			return appErrors.New(appErrors.ApplicationError, appErrors.FilesystemError, fmt.Sprintf("Cannot create directory: %s", cErr.Error()))
		}
	}

	return nil
}

func getStateDirectory(state string) string {
	if state == PROJECT_EXECUTION_STATE {
		return os.Getenv("CODE_PROJECT_STATE_DIR")
	} else if state == SINGLE_FILE_EXECUTION_STATE {
		return os.Getenv("SINGLE_FILE_STATE_DIR")
	}

	return ""
}
