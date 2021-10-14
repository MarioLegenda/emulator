package runner

import (
	"io/ioutil"
	"os/exec"
	"syscall"
	"therebelsource/emulator/appErrors"
	"time"
)

func getTimeout(t string, execState string, userState string) time.Duration {
	if execState == "documentation_private" {
		return 30 * time.Second
	}

	if userState == "anonymous" {
		return 5 * time.Second
	}

	if userState == "authenticated" {
		return 15 * time.Second
	}


	return 5 * time.Second
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
				// TODO: notify by slack that the container could not be stopped and time happened
			}
		}
	}
}

