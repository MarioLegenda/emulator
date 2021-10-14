package runner

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"os/exec"
	"strings"
	"therebelsource/emulator/appErrors"
	"time"
)

type SingleFileRunFn func(br SingleFileBuildResult) (SingleFileRunResult, *appErrors.Error)

type SingleFileRunResult struct {
	Success bool `json:"success"`
	Result string `json:"result"`
	Timeout int `json:"timeout"`
}

func createSingleFileRunner() SingleFileRunFn {
	return func(br SingleFileBuildResult) (SingleFileRunResult, *appErrors.Error) {
		containerName := uuid.New().String()

		commandFactory := RunCommandFactory{}

		dockerRunCommand := commandFactory.CreateCommand(containerName, br.ExecutionDirectory, br.FileName, br.Environment, br.DirectoryName)

		context := context.TODO()
		
		timeout := getTimeout("blog", "single_file", "anonymous")

		var outb, errb bytes.Buffer
		var out string
		var success bool
		var runResult SingleFileRunResult

		tc := make(chan string, 1)
		pidC := make(chan int, 1)

		go func() {
			chown := exec.Command("chown", "-R", "dockeruser:dockerusergroup", br.StateDirectory)
			chmod := exec.Command("chmod", "-R", "777", br.ExecutionDirectory)
			chown.Start()
			chmod.Start()

			cmd := exec.Command("docker", dockerRunCommand...)
			
			cmd.Stderr = &errb
			cmd.Stdout = &outb

			startErr := cmd.Start()
			pidC <- cmd.Process.Pid

			if startErr == nil {
				// TODO: Handle wait error properly
				waitErr := cmd.Wait()

				if waitErr != nil {
					//fmt.Printf("Wait error: %s\n", waitErr.Error())
				}
			}

			tc <- "Finished"
		}()

		runResult = SingleFileRunResult{}

		select {
		case res := <-tc:
			res = res
			outE := errb.String()
			outS := outb.String()

			if outE != "" {
				success = false
				out = outE
			} else {
				success = true

				if br.Environment.Name == "go" {
					success = true
					out = outS
				} else if br.Environment.Name == "rust" {
					success = true
					out = outS
				} else if br.Environment.Name == "haskell" {
					split := strings.Split(outS, "...")

					if len(split) == 2 {
						out = split[1]
					} else {
						out = outS
					}
				} else {
					output, err := readFile(fmt.Sprintf("%s/%s", br.ExecutionDirectory, "output.txt"))

					if err != nil {
						success = false
						out = ""
					} else {
						out = output
					}
				}
			}

			break
		case <-time.After(timeout):
			runnerBalancer.addJob(job{
				containerName: containerName,
				pid:           <- pidC,
			})

			runResult.Success = false
			runResult.Result = "timeout"
			runResult.Timeout = int(timeout) / (int(time.Millisecond) * 1000)

			return runResult, nil
		case <-context.Done():
			runnerBalancer.addJob(job{
				containerName: containerName,
				pid:           <- pidC,
			})

			runResult.Success = false
			runResult.Result = "timeout"
			runResult.Timeout = int(timeout) / (int(time.Millisecond) * 1000)

			return runResult, nil
		}

		runResult.Result = out
		runResult.Success = success
		runResult.Timeout = int(timeout) / (int(time.Millisecond) * 1000)

		return runResult, nil
	}
}

