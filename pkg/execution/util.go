package execution

import (
	"emulator/pkg/logger"
	"fmt"
	"os/exec"
)

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

func FinalCleanup(log bool) {
	stopAll := exec.Command("docker", "stop", "$(docker ps -a -q)")
	err := stopAll.Run()

	if err != nil {
		if log {
			logger.Warn(fmt.Sprintf("Cannot stop all containers with docker stop command with error: %s. Containers were probably already been stopped. Continuing...", err.Error()))
		}
	}

	rmAll := exec.Command("/bin/bash", []string{"-c", "docker rm -f $(docker ps -a -q)"}...)
	err = rmAll.Run()

	if err != nil {
		if log {
			logger.Warn(fmt.Sprintf("Cannot remove all containers with error: %s. Containers were probably already removed. Continuing...", err.Error()))
		}
	}
}
