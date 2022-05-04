package runners

import "therebelsource/emulator/appErrors"

type Result struct {
	Result  string
	Success bool
	Error   *appErrors.Error
}

type NodeExecParams struct {
	ContainerName      string
	ExecutionDirectory string
	ContainerDirectory string
	ExecutionFile      string
}
