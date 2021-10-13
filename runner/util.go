package runner

import (
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
	"therebelsource/emulator/appErrors"
)

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

func readFile(path string) (string, *appErrors.Error) {
	buff, err := ioutil.ReadFile(path)

	if err != nil {
		return "", appErrors.New(appErrors.ApplicationError, appErrors.FilesystemError, err.Error())
	}

	return string(buff), nil
}

func stopDockerContainer(containerName string, pid int) {
	var stopCmd *exec.Cmd

	stopCmd = exec.Command("docker", []string{"container", "stop", containerName}...)
	stopErr := stopCmd.Run()

	if stopErr != nil {
		var rmCmd *exec.Cmd

		rmCmd = exec.Command("docker", []string{"rm", "-f", containerName}...)
		rmErr := rmCmd.Run()

		if rmErr != nil {
			killErr := syscall.Kill(pid, 9)

			if killErr != nil {
				// TODO: notify by slack that the container could not be stopped
			}
		}
	}
}

