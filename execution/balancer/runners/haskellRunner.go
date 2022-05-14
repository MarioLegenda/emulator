package runners

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"therebelsource/emulator/appErrors"
	"time"
)

type HaskellExecParams struct {
	ContainerName      string
	ExecutionDirectory string
	ContainerDirectory string
	ExecutionFile      string
}

func haskellRunner(params HaskellExecParams) Result {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()

	var outb, errb string
	var runResult Result

	tc := make(chan string)
	pidC := make(chan int, 1)

	go func() {
		cmd := exec.Command("docker", []string{"exec", params.ContainerName, "/bin/bash", "-c", fmt.Sprintf("cd %s && ghc %s > output.txt && ./%s > output.txt", params.ContainerDirectory, params.ExecutionFile, params.ExecutionFile[:len(params.ExecutionFile)-3])}...)
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
		if startErr == nil {
			pidC <- cmd.Process.Pid

			a, _ := io.ReadAll(errPipe)
			b, _ := io.ReadAll(outPipe)
			errb = string(a)
			outb = string(b)

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
			out := makeRunDecision(errb, outb, params.ExecutionDirectory)
			if out != "" {
				runResult.Success = true
				runResult.Result = out
				runResult.Error = nil
			}

			destroyContainerProcess(extractUniqueIdentifier(params.ExecutionFile[:len(params.ExecutionFile)-3], false), true)
			destroy(params.ExecutionDirectory)
			return runResult
		}

		out := makeRunDecision(errb, outb, params.ExecutionDirectory)
		if out != "" {
			runResult.Success = true
			runResult.Result = out
			runResult.Error = nil
		}

		closeExecSession(<-pidC)
		destroy(params.ExecutionDirectory)

		break
	case <-ctx.Done():
		destroyContainerProcess(extractUniqueIdentifier(params.ExecutionFile[:len(params.ExecutionFile)-3], false), true)
		closeExecSession(<-pidC)
		destroy(params.ExecutionDirectory)
		close(pidC)
		return Result{
			Result:  "",
			Success: false,
			Error:   appErrors.New(appErrors.ApplicationError, appErrors.TimeoutError, "Code execution timeout!"),
		}
	}

	runResult.Error = nil

	return runResult
}
