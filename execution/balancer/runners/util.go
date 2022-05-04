package runners

import (
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
	"therebelsource/emulator/appErrors"
)

func readFile(path string) (string, *appErrors.Error) {
	buff, err := ioutil.ReadFile(path)

	if err != nil {
		return "", appErrors.New(appErrors.ApplicationError, appErrors.FilesystemError, err.Error())
	}

	return string(buff), nil
}

func destroy(path string) {
	err := os.RemoveAll(path)

	if err != nil {
		cmd := exec.Command("rm", []string{"-f", path}...)

		err := cmd.Run()

		if err != nil {
			// TODO: SEND SLACK ERROR AND LOG
		}
	}
}

func closeExecSession(pid int) {
	syscall.Kill(pid, 9)
}
