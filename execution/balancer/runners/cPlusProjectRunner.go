package runners

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"therebelsource/emulator/appErrors"
	"time"
)

type CPlusProjectExecParams struct {
	ContainerName       string
	ContainerDirectory  string
	ExecutionDirectory  string
	BinaryFileName      string
	CompilationFileName string
}

func cPlusProjectRunner(params CPlusProjectExecParams) Result {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()

	var outb, errb string
	var runResult Result

	tc := make(chan string)
	pidC := make(chan int, 1)

	process := fmt.Sprintf(
		"cd %s && g++ %s -o %s > output.txt && ./%s > output.txt",
		params.ContainerDirectory,
		params.CompilationFileName,
		params.BinaryFileName,
		params.BinaryFileName,
	)

	go func() {
		cmd := exec.Command("docker", []string{"exec", params.ContainerName, "/bin/sh", "-c", process}...)

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

			destroyContainerProcess(extractUniqueIdentifier(params.BinaryFileName, false), true)
			destroy(params.ExecutionDirectory)
			return runResult
		}

		out := makeRunDecision(errb, outb, params.ExecutionDirectory)
		runResult.Success = true
		runResult.Result = out
		runResult.Error = nil

		closeExecSession(<-pidC)
		destroy(params.ExecutionDirectory)

		break
	case <-ctx.Done():
		destroyContainerProcess(extractUniqueIdentifier(params.BinaryFileName, false), true)
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
