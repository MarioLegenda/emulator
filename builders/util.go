package builders

import (
	"fmt"
	"os"
	"path/filepath"
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

func removeDirectories(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
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


