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

type GoExecParams struct {
	ContainerName      string
	ExecutionDirectory string
	ContainerDirectory string
	ExecutionFile      string
}

type name string

type language struct {
	name name `json:"name"`
}

var node14 = language{
	name: "node_v14_x",
}

var nodeLts = language{
	name: "node_latest",
}

var nodeEsm = language{
	name: "node_latest_esm",
}

var goLang = language{
	name: "go",
}
