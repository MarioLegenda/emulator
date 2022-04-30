package runner

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"therebelsource/emulator/appErrors"
	"time"
)

type SingleFileRunFn func(br SingleFileBuildResult) (SingleFileRunResult, *appErrors.Error)

type SingleFileRunResult struct {
	Success bool   `json:"success"`
	Result  string `json:"result"`
	Timeout int    `json:"timeout"`
}

func createSingleFileRunner() SingleFileRunFn {
	return func(br SingleFileBuildResult) (SingleFileRunResult, *appErrors.Error) {
		context := context.TODO()

		userState := "anonymous"
		if br.Timeout == 15 {
			userState = "authenticated"
		}

		timeout := getTimeout("blog", "single_file", userState)

		var outb, errb bytes.Buffer
		var out string
		var success bool
		var runResult SingleFileRunResult

		tc := make(chan string, 1)
		pidC := make(chan int, 1)

		go func() {
			cmd := exec.Command("docker", br.Args...)
			fmt.Println(cmd.String())

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

					if out == "" {
						output, err := readFile(fmt.Sprintf("%s/%s", br.ExecutionDirectory, "output.txt"))

						if err != nil {
							success = false
							out = ""
						} else {
							out = output
						}
					}
				} else if br.Environment.Name == "rust" {
					success = true
					out = outS
				} else if br.Environment.Name == "c_sharp_mono" {
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
				containerName: br.ContainerName,
				pid:           <-pidC,
			})

			runResult.Success = false
			runResult.Result = "timeout"
			runResult.Timeout = int(timeout) / (int(time.Millisecond) * 1000)

			return runResult, nil
		case <-context.Done():
			runnerBalancer.addJob(job{
				containerName: br.ContainerName,
				pid:           <-pidC,
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
