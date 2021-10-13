package runner

import "fmt"

type CommandBuilder struct {
	Commands *[]string
}

func (cb CommandBuilder) NewVolume(dir string, permission string) CommandBuilder {
	*cb.Commands = append(*cb.Commands, "-v")

	*cb.Commands = append(*cb.Commands, fmt.Sprintf("%s:/app:%s", dir, permission))

	return cb
}

func (cb CommandBuilder) NewVolumeFull(local string, mount string, permission string) CommandBuilder {
	*cb.Commands = append(*cb.Commands, "-v")

	*cb.Commands = append(*cb.Commands, fmt.Sprintf("%s:/%s:%s", local, mount, permission))

	return cb
}

func (cb CommandBuilder) NewNetwork(network string) CommandBuilder {
	*cb.Commands = append(*cb.Commands, fmt.Sprintf("--network=%s", network))

	return cb
}

func (cb CommandBuilder) Init() CommandBuilder {
	*cb.Commands = append(*cb.Commands, "--init")

	return cb
}

func (cb CommandBuilder) Readonly() CommandBuilder {
	*cb.Commands = append(*cb.Commands, "--read-only")

	return cb
}

func (cb CommandBuilder) SecurityOps() CommandBuilder {
	*cb.Commands = append(*cb.Commands, "--security-ops")

	return cb
}

func (cb CommandBuilder) AllocatePseudoTty() CommandBuilder {
	*cb.Commands = append(*cb.Commands, "-t")

	return cb
}

func (cb CommandBuilder) RemoveAfterFinished() CommandBuilder {
	*cb.Commands = append(*cb.Commands, "--rm")

	return cb
}

func (cb CommandBuilder) Name(name string) CommandBuilder {
	*cb.Commands = append(*cb.Commands, "--name")
	*cb.Commands = append(*cb.Commands, name)

	return cb
}

func (cb CommandBuilder) Tag(tag string) CommandBuilder {
	*cb.Commands = append(*cb.Commands, tag)

	return cb
}

func (cb CommandBuilder) Shell(shell string) CommandBuilder {
	*cb.Commands = append(*cb.Commands, shell)

	return cb
}

func (cb CommandBuilder) Exec(exec string) CommandBuilder {
	*cb.Commands = append(*cb.Commands, "-c")
	*cb.Commands = append(*cb.Commands, exec)

	return cb
}

func (cb CommandBuilder) SendToStd(std string) CommandBuilder {
	*cb.Commands = append(*cb.Commands, std)

	return cb
}

func (cb CommandBuilder) User(user string, group string) CommandBuilder {
	*cb.Commands = append(*cb.Commands, "--user")
	*cb.Commands = append(*cb.Commands, fmt.Sprintf("%s:%s", user, group))

	return cb
}

func (cb CommandBuilder) Run() {
	*cb.Commands = append([]string{"run"}, *cb.Commands...)
}

