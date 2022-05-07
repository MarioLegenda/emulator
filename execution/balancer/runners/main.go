package runners

import (
	"fmt"
	"os"
	"therebelsource/emulator/execution/balancer/builders"
)

type Params struct {
	BuilderType   string
	ExecutionType string

	ContainerName string

	EmulatorName      string
	EmulatorExtension string
	EmulatorText      string
}

func Run(params Params) Result {
	if params.EmulatorName == string(nodeLts.name) && params.BuilderType == "single_file" && params.ExecutionType == "single_file" {
		build, err := builders.NodeSingleFileBuild(builders.InitNodeParams(
			params.EmulatorExtension,
			params.EmulatorText,
			fmt.Sprintf("%s/%s", os.Getenv("SINGLE_FILE_STATE_DIR"), params.ContainerName),
		))

		if err != nil {
			return Result{
				Result:  "",
				Success: false,
				Error:   err,
			}
		}

		return nodeRunner(NodeExecParams{
			ExecutionDirectory: build.ExecutionDirectory,
			ContainerDirectory: build.ContainerDirectory,
			ExecutionFile:      build.FileName,
			ContainerName:      params.ContainerName,
		})
	}

	if params.EmulatorName == string(goLang.name) && params.BuilderType == "single_file" && params.ExecutionType == "single_file" {
		build, err := builders.GoSingleFileBuild(builders.InitGoParams(
			params.EmulatorExtension,
			params.EmulatorText,
			fmt.Sprintf("%s/%s", os.Getenv("SINGLE_FILE_STATE_DIR"), params.ContainerName),
		))

		fmt.Println(err)

		if err != nil {
			return Result{
				Result:  "",
				Success: false,
				Error:   err,
			}
		}

		return goRunner(GoExecParams{
			ExecutionDirectory: build.ExecutionDirectory,
			ContainerDirectory: build.ContainerDirectory,
			ExecutionFile:      build.FileName,
			ContainerName:      params.ContainerName,
		})
	}

	return Result{}
}
