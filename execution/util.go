package execution

import (
	"fmt"
	"os/exec"
	"therebelsource/emulator/logger"
	"therebelsource/emulator/slack"
)

func makeBlocks(num int, delimiter int) [][]int {
	portions := num / 5
	leftover := num % 5
	if leftover != 0 {
		portions++
	}

	blocks := make([][]int, 0)
	current := 0
	for i := 0; i < portions; i++ {
		b := make([]int, 0)
		d := delimiter

		if i == portions-1 {
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
	stopAll := exec.Command("/bin/bash", []string{"-c", "docker stop $(docker ps -a -q)"}...)
	err := stopAll.Run()

	if err != nil {
		if log {
			logger.Warn(fmt.Sprintf("Cannot stop all containers with error: %s", err.Error()))
			slack.SendLog("Cannot stop containers", err.Error(), "deploy_log")
		}
	}

	rmAll := exec.Command("/bin/bash", []string{"-c", "docker rm -f $(docker ps -a -q)"}...)
	err = rmAll.Run()

	if err != nil {
		if log {
			logger.Warn(fmt.Sprintf("Cannot remove all containers with error: %s", err.Error()))
			slack.SendLog("Cannot remove containers", err.Error(), "deploy_log")
		}
	}
}
