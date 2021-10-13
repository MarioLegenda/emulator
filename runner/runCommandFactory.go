package runner

import (
	"fmt"
	"strings"
)

type RunCommandFactory struct {}

func (cf *RunCommandFactory) CreateProjectNodeCommand(containerName string, projectName string, fileName string, lang Language) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	cb.
		NewNetwork("none").
		Readonly().
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolume(projectName, "rw").
		Name(containerName).
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("node %s &> output.txt", fileName)).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreateCCommand(containerName string, projectName string, fileName string, lang Language) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	execName := strings.Split(fileName, ".")[0]

	cb.
		NewNetwork("none").
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolume(projectName, "rw").
		Name(containerName).
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("gcc -o %s %s &> output.txt && ./%s &> output.txt", execName, fileName, execName)).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreateCPlusCommand(containerName string, projectName string, fileName string, lang Language) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	execName := strings.Split(fileName, ".")[0]

	cb.
		NewNetwork("none").
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolume(projectName, "rw").
		Name(containerName).
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("g++ -o %s %s &> output.txt && ./%s &> output.txt", execName, fileName, execName)).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreateHaskellCommand(containerName string, projectName string, fileName string, lang Language) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	cb.
		NewNetwork("none").
		AllocatePseudoTty().
		RemoveAfterFinished().
		User("dockeruser", "dockerusergroup").
		NewVolume(projectName, "rw").
		Name(containerName).
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("ghc %s && ./%s", fileName, fileName[:len(fileName)-3])).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreateGoCommand(containerName string, projectName string, lang Language) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	cb.
		NewNetwork("none").
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolumeFull(projectName, fmt.Sprintf("app/src/%s", projectName), "rw").
		Name(containerName).
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("go run %s | tee output.txt", projectName)).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreatePython2Command(containerName string, projectName string, fileName string, lang Language) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	cb.
		NewNetwork("none").
		Readonly().
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolume(projectName, "rw").
		Name(containerName).
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("python %s &> output.txt", fileName)).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreatePython3Command(containerName string, projectName string, fileName string, lang Language) []string {

	cb := CommandBuilder{Commands: &[]string{}}

	cb.
		NewNetwork("none").
		Readonly().
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolume(projectName, "rw").
		Name(containerName).
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("python3 %s &> output.txt", fileName)).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreateRubyCommand(containerName string, projectName string, fileName string, lang Language) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	cb.
		NewNetwork("none").
		Readonly().
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolume(projectName, "rw").
		Name(containerName).
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("ruby %s &> output.txt", fileName)).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreatePHP74Command(containerName string, projectName string, fileName string, lang Language) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	cb.
		NewNetwork("none").
		Readonly().
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolume(projectName, "rw").
		Name(containerName).
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("php %s &> output.txt", fileName)).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreateRustCommand(containerName string, projectName string, fileName string, lang Language) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	cb.
		NewNetwork("none").
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolumeFull(projectName, "app", "rw").
		Name(containerName).
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec("cargo run --quiet | tee output.txt").
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}


func (cf *RunCommandFactory) CreateCommand(containerName, projectName string, fileName string, lang Language, state string) []string {
	if lang.Name == node12.Name || lang.Name == nodeLts.Name {
		return cf.CreateProjectNodeCommand(containerName, projectName, fileName, lang)
	} else if lang.Name == goLang.Name {
		return cf.CreateGoCommand(containerName, projectName, lang)
	} else if lang.Name == python2.Name {
		return cf.CreatePython2Command(containerName, projectName, fileName, lang)
	} else if lang.Name == python3.Name {
		return cf.CreatePython3Command(containerName, projectName, fileName, lang)
	} else if lang.Name == ruby.Name {
		return cf.CreateRubyCommand(containerName, projectName, fileName, lang)
	} else if lang.Name == php74.Name {
		return cf.CreatePHP74Command(containerName, projectName, fileName, lang)
	} else if lang.Name == rust.Name {
		return cf.CreateRustCommand(containerName, projectName, fileName, lang)
	} else if lang.Name == haskell.Name {
		return cf.CreateHaskellCommand(containerName, projectName, fileName, lang)
	} else if lang.Name == c.Name {
		return cf.CreateCCommand(containerName, projectName, fileName, lang)
	} else if lang.Name == cPlus.Name {
		return cf.CreateCPlusCommand(containerName, projectName, fileName, lang)
	}

	return []string{}
}

