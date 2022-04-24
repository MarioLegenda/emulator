package runner

import (
	"fmt"
	"strings"
)

type RunCommandFactory struct{}

func (cf *RunCommandFactory) CreateProjectNodeCommand(containerName string, projectName string, fileName string, lang *Language) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	cb.
		NewNetwork("none").
		Readonly().
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolume(projectName, "rw").
		Name(containerName).
		Init().
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("node %s > output.txt", fileName)).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreateCCommand(containerName string, projectName string, fileName string, lang *Language) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	execName := strings.Split(fileName, ".")[0]

	cb.
		NewNetwork("none").
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolume(projectName, "rw").
		Name(containerName).
		Init().
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("gcc -o %s %s > output.txt && ./%s > output.txt", execName, fileName, execName)).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreateCProjectCommand(containerName string, projectName string, volumePath string, paths string, lang *Language) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	cb.
		NewNetwork("none").
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolume(volumePath, "rw").
		Name(containerName).
		Init().
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("gcc -o %s %s > output.txt && ./%s > output.txt", projectName, paths, projectName)).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreateCPlusCommand(containerName string, projectName string, fileName string, lang *Language) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	execName := strings.Split(fileName, ".")[0]

	cb.
		NewNetwork("none").
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolume(projectName, "rw").
		Name(containerName).
		Init().
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("g++ -o %s %s > output.txt && ./%s > output.txt", execName, fileName, execName)).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreateCPlusProjectCommand(containerName string, projectName string, volumePath string, paths string, lang *Language) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	cb.
		NewNetwork("none").
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolume(volumePath, "rw").
		Name(containerName).
		Init().
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("g++ -o %s %s > output.txt && ./%s > output.txt", projectName, paths, projectName)).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreateHaskellProjectCommand(containerName string, volumePath string, lang *Language) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	cb.
		NewNetwork("none").
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolume(volumePath, "rw").
		Name(containerName).
		Init().
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("ghc main.hs && ./main")).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreateHaskellCommand(containerName string, projectName string, fileName string, lang *Language) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	cb.
		NewNetwork("none").
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolume(projectName, "rw").
		Name(containerName).
		Init().
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("ghc %s && ./%s", fileName, fileName[:len(fileName)-3])).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreateGoLinkedProjectCommand(containerName string, projectName string, lang *Language, directoryName string) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	cb.
		NewNetwork("none").
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolumeFull(projectName, fmt.Sprintf("app/src/%s", directoryName), "rw").
		Name(containerName).
		Init().
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("cd /app/src/%s && go mod init > /dev/null 2>&1 && go run %s | tee output.txt", directoryName, directoryName)).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreateGoCommand(containerName string, projectName string, lang *Language, directoryName string) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	cb.
		NewNetwork("none").
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolumeFull(projectName, fmt.Sprintf("app/src/%s", directoryName), "rw").
		Name(containerName).
		Init().
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("cd /app/src/%s && go mod init > /dev/null 2>&1 && go run %s | tee output.txt", directoryName, directoryName)).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreateCSharpCommand(containerName string, projectName string, lang *Language, directoryName string) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	cb.
		NewNetwork("none").
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolumeFull(projectName, fmt.Sprintf("app/src/%s", directoryName), "rw").
		Name(containerName).
		Init().
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("cd /app/src/%s && mcs -out:app.exe -pkg:dotnet -recurse:'*.cs' && mono app.exe | tee output.txt", directoryName)).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreatePython2Command(containerName string, projectName string, fileName string, lang *Language) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	cb.
		NewNetwork("none").
		Readonly().
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolume(projectName, "rw").
		Name(containerName).
		Init().
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("python %s &> output.txt", fileName)).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreatePython3Command(containerName string, projectName string, fileName string, lang *Language) []string {

	cb := CommandBuilder{Commands: &[]string{}}

	cb.
		NewNetwork("none").
		Readonly().
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolume(projectName, "rw").
		Name(containerName).
		Init().
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("python3 %s &> output.txt", fileName)).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreateRubyCommand(containerName string, projectName string, fileName string, lang *Language) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	cb.
		NewNetwork("none").
		Readonly().
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolume(projectName, "rw").
		Name(containerName).
		Init().
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("ruby %s &> output.txt", fileName)).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreatePHP74Command(containerName string, projectName string, fileName string, lang *Language) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	cb.
		NewNetwork("none").
		Readonly().
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolume(projectName, "rw").
		Name(containerName).
		Init().
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec(fmt.Sprintf("php %s &> output.txt", fileName)).
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreateRustCommand(containerName string, projectName string, fileName string, lang *Language) []string {
	cb := CommandBuilder{Commands: &[]string{}}

	cb.
		NewNetwork("none").
		AllocatePseudoTty().
		RemoveAfterFinished().
		NewVolumeFull(projectName, "app", "rw").
		Name(containerName).
		Init().
		Tag(string(lang.Tag)).
		Shell("/bin/sh").
		Exec("cargo run --quiet | tee output.txt").
		SendToStd("/dev/stderr").
		Run()

	args := *cb.Commands

	return args
}

func (cf *RunCommandFactory) CreateCommand(containerName, projectName string, fileName string, lang *Language, directoryName string) []string {
	if lang.Name == Node14.Name || lang.Name == NodeLts.Name || lang.Name == NodeEsm.Name {
		return cf.CreateProjectNodeCommand(containerName, projectName, fileName, lang)
	} else if lang.Name == GoLang.Name {
		return cf.CreateGoCommand(containerName, projectName, lang, directoryName)
	} else if lang.Name == Python2.Name {
		return cf.CreatePython2Command(containerName, projectName, fileName, lang)
	} else if lang.Name == Python3.Name {
		return cf.CreatePython3Command(containerName, projectName, fileName, lang)
	} else if lang.Name == Ruby.Name {
		return cf.CreateRubyCommand(containerName, projectName, fileName, lang)
	} else if lang.Name == Php74.Name {
		return cf.CreatePHP74Command(containerName, projectName, fileName, lang)
	} else if lang.Name == Rust.Name {
		return cf.CreateRustCommand(containerName, projectName, fileName, lang)
	} else if lang.Name == Haskell.Name {
		return cf.CreateHaskellCommand(containerName, projectName, fileName, lang)
	} else if lang.Name == CLang.Name {
		return cf.CreateCCommand(containerName, projectName, fileName, lang)
	} else if lang.Name == CPlus.Name {
		return cf.CreateCPlusCommand(containerName, projectName, fileName, lang)
	} else if lang.Name == CSharpMono.Name {
		return cf.CreateCSharpCommand(containerName, projectName, lang, directoryName)
	}

	return []string{}
}
