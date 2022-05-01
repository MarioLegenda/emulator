package containerFactory

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

func getVolumeDirectory(volume string) string {
	return fmt.Sprintf("%s/%s", os.Getenv("SINGLE_FILE_STATE_DIR"), volume)
}

func isContainerRunning(name string) bool {
	cmd := exec.Command("docker", []string{
		"container",
		"inspect",
		"-f",
		"'{{.State.Status}}'",
		name,
	}...)

	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(err)
		return false
	}

	o := strings.Trim(string(out), " ")

	match, _ := regexp.MatchString("running", o)

	return match
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
