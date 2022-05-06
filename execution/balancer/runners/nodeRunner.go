package runners

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"therebelsource/emulator/appErrors"
	"time"
)

var t string = "ps aux | awk '/app\\/71db77e0-38bc-44ff-a65f-9d33ed6b7477\\/71db77e0-38bc-44ff-a65f-9d33ed6b7477.js/ { print $2}'"

func NodeRunner(params NodeExecParams) Result {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()

	var outb, errb string
	var out string
	var success bool
	var runResult Result

	tc := make(chan string)
	pidC := make(chan int, 1)

	process := fmt.Sprintf("%s/%s", params.ContainerDirectory, params.ExecutionFile)

	go func() {
		cmd := exec.Command("docker", []string{"exec", params.ContainerName, "node", process}...)

		errPipe, err := cmd.StderrPipe()

		if err != nil {
			runResult.Error = appErrors.New(appErrors.ApplicationError, appErrors.ExecutionStartError, "Execution failed!")

			tc <- "error"

			return
		}

		outPipe, err := cmd.StdoutPipe()

		if err != nil {
			runResult.Error = appErrors.New(appErrors.ApplicationError, appErrors.ExecutionStartError, "Execution failed!")

			tc <- "error"

			return
		}

		startErr := cmd.Start()
		pidC <- cmd.Process.Pid

		a, _ := io.ReadAll(errPipe)
		b, _ := io.ReadAll(outPipe)
		errb = string(a)
		outb = string(b)

		if startErr == nil {
			waitErr := cmd.Wait()

			if waitErr != nil {
				runResult.Error = appErrors.New(appErrors.ApplicationError, appErrors.ExecutionStartError, "Execution failed!")

				tc <- "error"

				return
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
			destroyContainerProcess(extractExecDirUniqueIdentifier(params.ExecutionDirectory))
			destroy(params.ExecutionDirectory)
			return runResult
		}

		if errb != "" {
			success = false
			out = errb
		} else if outb != "" {
			success = true
			out = outb
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
		destroyContainerProcess(extractExecDirUniqueIdentifier(params.ExecutionFile))
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
