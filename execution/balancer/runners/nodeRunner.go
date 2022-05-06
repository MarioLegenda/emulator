package runners

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"therebelsource/emulator/appErrors"
	"time"
)

var t string = "ps aux | awk '/app\\/71db77e0-38bc-44ff-a65f-9d33ed6b7477\\/71db77e0-38bc-44ff-a65f-9d33ed6b7477.js/ { print $2}'"

func NodeRunner(params NodeExecParams) Result {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()

	var outb, errb bytes.Buffer
	var out string
	var success bool
	var runResult Result

	tc := make(chan string)
	pidC := make(chan int, 1)

	go func() {
		cmd := exec.Command("docker", []string{"exec", params.ContainerName, "node", fmt.Sprintf("%s/%s", params.ContainerDirectory, params.ExecutionFile)}...)
		fmt.Println(cmd.String())

		cmd.Stderr = &errb
		cmd.Stdout = &outb

		startErr := cmd.Start()
		pidC <- cmd.Process.Pid

		if startErr == nil {
			waitErr := cmd.Wait()

			if waitErr != nil {
				fmt.Println(waitErr)
				runResult.Error = appErrors.New(appErrors.ApplicationError, appErrors.ExecutionStartError, "Execution failed!")

				tc <- "error"
				//fmt.Printf("Wait error: %s\n", waitErr.Error())
			}
		}

		if startErr != nil {
			runResult.Error = appErrors.New(appErrors.ApplicationError, appErrors.ExecutionStartError, "Execution failed!")

			tc <- "error"

			return
		}

		tc <- "finished"
	}()

	select {
	case res := <-tc:
		if res == "error" {
			destroy(params.ExecutionDirectory)
			destroy(params.ExecutionDirectory)
			return runResult
		}

		outE := errb.String()
		outS := outb.String()

		if outE != "" {
			success = false
			out = outE
		} else if outS != "" {
			success = true
			out = outS
		} else {
			success = true

			output, err := readFile(fmt.Sprintf("%s/%s", params.ExecutionDirectory, "output.txt"))

			if err != nil {
				success = false
				out = ""
			} else {
				out = output
			}
		}

		closeExecSession(<-pidC)
		destroy(params.ExecutionDirectory)

		break
	case <-ctx.Done():
		closeExecSession(<-pidC)
		destroy(params.ExecutionDirectory)
		close(pidC)
		return Result{
			Result:  "",
			Success: false,
			Error:   appErrors.New(appErrors.ApplicationError, appErrors.TimeoutError, "Code execution timeout!"),
		}
	}

	runResult.Result = out
	runResult.Success = success
	runResult.Error = nil

	return runResult
}
