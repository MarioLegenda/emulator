package builders

import (
	"fmt"
	"os"
	"therebelsource/emulator/appErrors"
)

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

func getStateDirectory(state string) string {
	if state == "dev" {
		return os.Getenv("DEV_STATE_DIR")
	} else if state == "prod" {
		return os.Getenv("PROD_STATE_DIR")
	} else if state == "session" {
		return os.Getenv("SESSION_STATE_DIR")
	} else if state == "single_file" {
		return os.Getenv("SINGLE_FILE_STATE_DIR")
	}

	return ""
}


