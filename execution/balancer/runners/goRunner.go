package runners

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"therebelsource/emulator/appErrors"
	"time"
)

func goRunner(params GoExecParams) Result {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()

	var outb, errb string
	var out string
	var success bool
	var runResult Result

	tc := make(chan string)
	pidC := make(chan int, 1)

	fmt.Println(params.ContainerName)

	go func() {
		cmd := exec.Command("docker", []string{"exec", params.ContainerName, fmt.Sprintf("cd %s && go mod init > /dev/null 2>&1 && go run . | tee output.txt", params.ContainerDirectory)}...)
		//cmd := exec.Command("docker", []string{"exec", params.ContainerName, fmt.Sprintf("cd %s", params.ContainerDirectory)}...)
		fmt.Println(cmd.String())
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

		fmt.Println(errb)
		fmt.Println(outb)

		if startErr == nil {
			waitErr := cmd.Wait()

			if waitErr != nil {
				fmt.Println(waitErr)
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
			//destroy(params.ExecutionDirectory)
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
		//destroy(params.ExecutionDirectory)

		break
	case <-ctx.Done():
		destroyContainerProcess(extractExecDirUniqueIdentifier(params.ExecutionFile))
		closeExecSession(<-pidC)
		//destroy(params.ExecutionDirectory)
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
