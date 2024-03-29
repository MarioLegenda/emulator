package containerFactory

import (
	"emulator/pkg/logger"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

func getVolumeDirectory(volume string) string {
	return fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), volume)
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

	if stopErr == nil {
		var rmCmd *exec.Cmd

		rmCmd = exec.Command("docker", []string{"rm", "-f", containerName}...)
		rmErr := rmCmd.Run()

		if rmErr != nil {
			logger.Warn(fmt.Sprintf("Container could not stop. rm error: %s", stopErr.Error()))

			killErr := syscall.Kill(pid, 9)

			if killErr != nil {
				logger.Warn(fmt.Sprintf("Container could not be killed. Kill error: %s", killErr.Error()))
			}
		}
	}

	if stopErr != nil {
		logger.Warn(fmt.Sprintf("Container could not stop. Stop error: %s", stopErr.Error()))
	}
}

func makeBlocks(num int, delimiter int) [][]int {
	portions := num / delimiter
	leftover := num % delimiter
	if leftover != 0 {
		portions++
	}

	blocks := make([][]int, 0)
	current := 0
	for i := 0; i < portions; i++ {
		b := make([]int, 0)
		d := delimiter

		if i == portions-1 && leftover != 0 {
			d = leftover
		}

		for a := 0; a < d; a++ {
			b = append(b, current)
			current++
		}

		blocks = append(blocks, b)
	}

	return blocks
}

func containersToSlice(containers map[string][]container) []container {
	s := make([]container, 0)
	for _, v := range containers {
		s = append(s, v...)
	}

	return s
}
