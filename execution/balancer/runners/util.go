package runners

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
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

func pidExists(pid int) (bool, error) {
	if pid <= 0 {
		return false, fmt.Errorf("invalid pid %v", pid)
	}
	proc, err := os.FindProcess(int(pid))
	if err != nil {
		return false, err
	}
	err = proc.Signal(syscall.Signal(0))
	if err == nil {
		return true, nil
	}
	if err.Error() == "os: process already finished" {
		return false, nil
	}
	errno, ok := err.(syscall.Errno)
	if !ok {
		return false, err
	}
	switch errno {
	case syscall.ESRCH:
		return false, nil
	case syscall.EPERM:
		return true, nil
	}
	return false, err
}

func destroyContainerProcess(processName string) {
	pids, ok := getContainerProcessPid(processName)

	if !ok {
		return
	}

	for _, pid := range pids {
		exists, err := pidExists(pid)

		if err != nil {
			//TODO: log here to slack, error should not happend
		}

		if exists {
			syscall.Kill(pid, 9)
		}
	}
}

func getContainerProcessPid(processName string) ([]int, bool) {
	cmd := exec.Command("/usr/bin/ps", "aux")
	pids := make([]int, 0)

	out, err := cmd.CombinedOutput()

	if err != nil {
		return []int{}, false
	}

	a := strings.Split(string(out), "\n")

	for _, i := range a {
		match, _ := regexp.MatchString(fmt.Sprintf("(app.*%s)$", processName), i)

		if match {
			m1 := regexp.MustCompile(`\s+`)
			repl := m1.ReplaceAllString(i, " ")

			s := strings.Split(repl, " ")

			if s[1] != "" {
				p, err := strconv.Atoi(s[1])

				if err != nil {
					return []int{}, false
				}

				pids = append(pids, p)
			}
		}
	}

	return pids, true
}

func extractExecDirUniqueIdentifier(dir string) string {
	s := strings.Split(dir, "/")

	return s[len(s)-1]
}
