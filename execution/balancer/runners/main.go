package runners

import (
	"fmt"
	"os"
	"therebelsource/emulator/execution/balancer/builders/linked"
	"therebelsource/emulator/execution/balancer/builders/project"
	"therebelsource/emulator/execution/balancer/builders/single"
	"therebelsource/emulator/repository"
)

type Params struct {
	BuilderType   string
	ExecutionType string

	ContainerName string

	EmulatorName      string
	EmulatorExtension string
	EmulatorText      string

	CodeProject   *repository.CodeProject
	Contents      []*repository.FileContent
	ExecutingFile *repository.File
	PackageName   string
}

func Run(params Params) Result {
	if params.EmulatorName == string(nodeLts.name) && params.BuilderType == "single_file" && params.ExecutionType == "single_file" {
		build, err := single.NodeSingleFileBuild(single.InitNodeParams(
			params.EmulatorExtension,
			params.EmulatorText,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
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

	if params.EmulatorName == string(nodeEsm.name) && params.BuilderType == "single_file" && params.ExecutionType == "single_file" {
		build, err := single.NodeSingleFileBuild(single.InitNodeParams(
			params.EmulatorExtension,
			params.EmulatorText,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
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
		build, err := single.GoSingleFileBuild(single.InitGoParams(
			params.EmulatorExtension,
			params.EmulatorText,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
		))

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

	if params.EmulatorName == string(ruby.name) && params.BuilderType == "single_file" && params.ExecutionType == "single_file" {
		build, err := single.RubySingleFileBuild(single.InitRubyParams(
			params.EmulatorExtension,
			params.EmulatorText,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
		))

		if err != nil {
			return Result{
				Result:  "",
				Success: false,
				Error:   err,
			}
		}

		return rubyRunner(RubyExecParams{
			ExecutionDirectory: build.ExecutionDirectory,
			ContainerDirectory: build.ContainerDirectory,
			ExecutionFile:      build.FileName,
			ContainerName:      params.ContainerName,
		})
	}

	if params.EmulatorName == string(php.name) && params.BuilderType == "single_file" && params.ExecutionType == "single_file" {
		build, err := single.PhpSingleFileBuild(single.InitPhpParams(
			params.EmulatorExtension,
			params.EmulatorText,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
		))

		if err != nil {
			return Result{
				Result:  "",
				Success: false,
				Error:   err,
			}
		}

		return phpRunner(PhpExecParams{
			ExecutionDirectory: build.ExecutionDirectory,
			ContainerDirectory: build.ContainerDirectory,
			ExecutionFile:      build.FileName,
			ContainerName:      params.ContainerName,
		})
	}

	if params.EmulatorName == string(python2.name) && params.BuilderType == "single_file" && params.ExecutionType == "single_file" {
		build, err := single.PythonSingleFileBuild(single.InitPythonParams(
			params.EmulatorExtension,
			params.EmulatorText,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
		))

		if err != nil {
			return Result{
				Result:  "",
				Success: false,
				Error:   err,
			}
		}

		return pythonRunner(PythonExecParams{
			ExecutionDirectory: build.ExecutionDirectory,
			ContainerDirectory: build.ContainerDirectory,
			ExecutionFile:      build.FileName,
			ContainerName:      params.ContainerName,
		})
	}

	if params.EmulatorName == string(python3.name) && params.BuilderType == "single_file" && params.ExecutionType == "single_file" {
		build, err := single.PythonSingleFileBuild(single.InitPythonParams(
			params.EmulatorExtension,
			params.EmulatorText,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
		))

		if err != nil {
			return Result{
				Result:  "",
				Success: false,
				Error:   err,
			}
		}

		return python3Runner(PythonExecParams{
			ExecutionDirectory: build.ExecutionDirectory,
			ContainerDirectory: build.ContainerDirectory,
			ExecutionFile:      build.FileName,
			ContainerName:      params.ContainerName,
		})
	}

	if params.EmulatorName == string(csharpMono.name) && params.BuilderType == "single_file" && params.ExecutionType == "single_file" {
		build, err := single.CsharpSingleFileBuild(single.InitCsharpParams(
			params.EmulatorExtension,
			params.EmulatorText,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
		))

		if err != nil {
			return Result{
				Result:  "",
				Success: false,
				Error:   err,
			}
		}

		return csharpRunner(CsharpExecParams{
			ExecutionDirectory: build.ExecutionDirectory,
			ContainerDirectory: build.ContainerDirectory,
			ExecutionFile:      build.FileName,
			ContainerName:      params.ContainerName,
		})
	}

	if params.EmulatorName == string(haskell.name) && params.BuilderType == "single_file" && params.ExecutionType == "single_file" {
		build, err := single.HaskellSingleFileBuild(single.InitHaskellParams(
			params.EmulatorExtension,
			params.EmulatorText,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
		))

		if err != nil {
			return Result{
				Result:  "",
				Success: false,
				Error:   err,
			}
		}

		return haskellRunner(HaskellExecParams{
			ExecutionDirectory: build.ExecutionDirectory,
			ContainerDirectory: build.ContainerDirectory,
			ExecutionFile:      build.FileName,
			ContainerName:      params.ContainerName,
		})
	}

	if params.EmulatorName == string(cLang.name) && params.BuilderType == "single_file" && params.ExecutionType == "single_file" {
		build, err := single.CSingleFileBuild(single.InitCParams(
			params.EmulatorExtension,
			params.EmulatorText,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
		))

		if err != nil {
			return Result{
				Result:  "",
				Success: false,
				Error:   err,
			}
		}

		return cRunner(CExecParams{
			ExecutionDirectory: build.ExecutionDirectory,
			ContainerDirectory: build.ContainerDirectory,
			ExecutionFile:      build.FileName,
			ContainerName:      params.ContainerName,
		})
	}

	if params.EmulatorName == string(cPlus.name) && params.BuilderType == "single_file" && params.ExecutionType == "single_file" {
		build, err := single.CPlusSingleFileBuild(single.InitCPlusParams(
			params.EmulatorExtension,
			params.EmulatorText,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
		))

		if err != nil {
			return Result{
				Result:  "",
				Success: false,
				Error:   err,
			}
		}

		return cplusRunner(CPlusExecParams{
			ExecutionDirectory: build.ExecutionDirectory,
			ContainerDirectory: build.ContainerDirectory,
			ExecutionFile:      build.FileName,
			ContainerName:      params.ContainerName,
		})
	}

	if params.EmulatorName == string(rust.name) && params.BuilderType == "single_file" && params.ExecutionType == "single_file" {
		build, err := single.RustSingleFileBuild(single.InitRustParams(
			params.EmulatorExtension,
			params.EmulatorText,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
		))

		if err != nil {
			return Result{
				Result:  "",
				Success: false,
				Error:   err,
			}
		}

		return rustRunner(RustExecParams{
			ExecutionDirectory: build.ExecutionDirectory,
			ContainerDirectory: build.ContainerDirectory,
			ExecutionFile:      build.FileName,
			ContainerName:      params.ContainerName,
		})
	}

	if params.EmulatorName == string(nodeLts.name) && params.BuilderType == "project" && params.ExecutionType == "project" {
		build, err := project.NodeProjectBuild(project.InitNodeParams(
			params.CodeProject,
			params.Contents,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
			params.ExecutingFile,
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

	if params.EmulatorName == string(nodeEsm.name) && params.BuilderType == "project" && params.ExecutionType == "project" {
		build, err := project.NodeProjectBuild(project.InitNodeParams(
			params.CodeProject,
			params.Contents,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
			params.ExecutingFile,
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

	if params.EmulatorName == string(goLang.name) && params.BuilderType == "project" && params.ExecutionType == "project" {
		build, err := project.GoProjectBuild(project.InitGoParams(
			params.CodeProject,
			params.Contents,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
			params.ExecutingFile,
			params.PackageName,
		))

		if err != nil {
			return Result{
				Result:  "",
				Success: false,
				Error:   err,
			}
		}

		return goProjectRunner(GoProjectExecParams{
			ExecutionDirectory: build.ExecutionDirectory,
			ContainerDirectory: params.PackageName,
			ExecutionFile:      build.FileName,
			ContainerName:      params.ContainerName,
		})
	}

	if params.EmulatorName == string(ruby.name) && params.BuilderType == "project" && params.ExecutionType == "project" {
		build, err := project.RubyProjectBuild(project.InitRubyParams(
			params.CodeProject,
			params.Contents,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
			params.ExecutingFile,
		))

		if err != nil {
			return Result{
				Result:  "",
				Success: false,
				Error:   err,
			}
		}

		return rubyRunner(RubyExecParams{
			ExecutionDirectory: build.ExecutionDirectory,
			ContainerDirectory: build.ContainerDirectory,
			ExecutionFile:      build.FileName,
			ContainerName:      params.ContainerName,
		})
	}

	if params.EmulatorName == string(cLang.name) && params.BuilderType == "project" && params.ExecutionType == "project" {
		build, err := project.CProjectBuild(project.InitCParams(
			params.CodeProject,
			params.Contents,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
		))

		if err != nil {
			return Result{
				Result:  "",
				Success: false,
				Error:   err,
			}
		}

		return cProjectRunner(CProjectExecParams{
			ContainerName:      params.ContainerName,
			ExecutionDirectory: build.ExecutionDirectory,
			ContainerDirectory: build.ContainerDirectory,
			ResolvedPaths:      build.ResolvedFiles,
			BinaryFileName:     build.BinaryFileName,
		})
	}

	if params.EmulatorName == string(cPlus.name) && params.BuilderType == "project" && params.ExecutionType == "project" {
		build, err := project.CPlusProjectBuild(project.InitCPlusParams(
			params.CodeProject,
			params.Contents,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
		))

		if err != nil {
			return Result{
				Result:  "",
				Success: false,
				Error:   err,
			}
		}

		return cPlusProjectRunner(CPlusProjectExecParams{
			ContainerName:      params.ContainerName,
			ExecutionDirectory: build.ExecutionDirectory,
			ContainerDirectory: build.ContainerDirectory,
			ResolvedPaths:      build.ResolvedFiles,
			BinaryFileName:     build.BinaryFileName,
		})
	}

	if params.EmulatorName == string(csharpMono.name) && params.BuilderType == "project" && params.ExecutionType == "project" {
		build, err := project.CsharpProjectFileBuild(project.InitCsharpProjectParams(
			params.CodeProject,
			params.Contents,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
			params.ExecutingFile,
		))

		if err != nil {
			return Result{
				Result:  "",
				Success: false,
				Error:   err,
			}
		}

		return csharpRunner(CsharpExecParams{
			ExecutionDirectory: build.ExecutionDirectory,
			ContainerDirectory: build.ContainerDirectory,
			ExecutionFile:      build.FileName,
			ContainerName:      params.ContainerName,
		})
	}

	if params.EmulatorName == string(python2.name) && params.BuilderType == "project" && params.ExecutionType == "project" {
		build, err := project.Python2ProjectBuild(project.InitPython2Params(
			params.CodeProject,
			params.Contents,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
			params.ExecutingFile,
		))

		if err != nil {
			return Result{
				Result:  "",
				Success: false,
				Error:   err,
			}
		}

		return pythonRunner(PythonExecParams{
			ExecutionDirectory: build.ExecutionDirectory,
			ContainerDirectory: build.ContainerDirectory,
			ExecutionFile:      build.FileName,
			ContainerName:      params.ContainerName,
		})
	}

	if params.EmulatorName == string(python3.name) && params.BuilderType == "project" && params.ExecutionType == "project" {
		build, err := project.Python3ProjectBuild(project.InitPython3Params(
			params.CodeProject,
			params.Contents,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
			params.ExecutingFile,
		))

		if err != nil {
			return Result{
				Result:  "",
				Success: false,
				Error:   err,
			}
		}

		return python3Runner(PythonExecParams{
			ExecutionDirectory: build.ExecutionDirectory,
			ContainerDirectory: build.ContainerDirectory,
			ExecutionFile:      build.FileName,
			ContainerName:      params.ContainerName,
		})
	}

	if params.EmulatorName == string(haskell.name) && params.BuilderType == "project" && params.ExecutionType == "project" {
		build, err := project.HaskellProjectBuild(project.InitHaskellProjectParams(
			params.CodeProject,
			params.Contents,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
		))

		if err != nil {
			return Result{
				Result:  "",
				Success: false,
				Error:   err,
			}
		}

		return haskellProjectRunner(HaskellExecProjectParams{
			ExecutionDirectory: build.ExecutionDirectory,
			ContainerDirectory: build.ContainerDirectory,
			ContainerName:      params.ContainerName,
		})
	}

	if params.EmulatorName == string(rust.name) && params.BuilderType == "project" && params.ExecutionType == "project" {
		build, err := project.RustProjectBuild(project.InitRustParams(
			params.CodeProject,
			params.Contents,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
		))

		if err != nil {
			return Result{
				Result:  "",
				Success: false,
				Error:   err,
			}
		}

		return rustRunner(RustExecParams{
			ExecutionDirectory: build.ExecutionDirectory,
			ContainerDirectory: build.ContainerDirectory,
			ContainerName:      params.ContainerName,
		})
	}

	if params.EmulatorName == string(php.name) && params.BuilderType == "project" && params.ExecutionType == "project" {
		build, err := project.Php74ProjectBuild(project.InitPhp74Params(
			params.CodeProject,
			params.Contents,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
			params.ExecutingFile,
		))

		if err != nil {
			return Result{
				Result:  "",
				Success: false,
				Error:   err,
			}
		}

		return phpRunner(PhpExecParams{
			ExecutionDirectory: build.ExecutionDirectory,
			ContainerDirectory: build.ContainerDirectory,
			ExecutionFile:      build.FileName,
			ContainerName:      params.ContainerName,
		})
	}

	if params.EmulatorName == string(ruby.name) && params.BuilderType == "linked" && params.ExecutionType == "linked" {
		build, err := linked.RubyProjectBuild(linked.InitRubyParams(
			params.CodeProject,
			params.Contents,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
			params.EmulatorText,
		))

		if err != nil {
			return Result{
				Result:  "",
				Success: false,
				Error:   err,
			}
		}

		return rubyRunner(RubyExecParams{
			ExecutionDirectory: build.ExecutionDirectory,
			ContainerDirectory: build.ContainerDirectory,
			ExecutionFile:      build.FileName,
			ContainerName:      params.ContainerName,
		})
	}

	if params.EmulatorName == string(nodeLts.name) && params.BuilderType == "linked" && params.ExecutionType == "linked" {
		build, err := linked.NodeProjectBuild(linked.InitNodeParams(
			params.CodeProject,
			params.Contents,
			fmt.Sprintf("%s/%s", os.Getenv("EXECUTION_DIR"), params.ContainerName),
			params.EmulatorText,
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

	return Result{}
}
