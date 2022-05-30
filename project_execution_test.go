package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"therebelsource/emulator/execution"
	"therebelsource/emulator/httpUtil"
	"therebelsource/emulator/repository"
	"therebelsource/emulator/staticTypes"
	_var "therebelsource/emulator/var"
)

var _ = GinkgoDescribe("Project execution tests", func() {
	GinkgoBeforeEach(func() {
		loadEnv()
		initRequiredDirectories(false)
	})

	GinkgoAfterEach(func() {
		gomega.Expect(os.RemoveAll(os.Getenv("EXECUTION_DIR"))).Should(gomega.BeNil())
	})

	GinkgoAfterAll(func() {
		cmd := exec.Command("/usr/bin/docker", "rm", "-f", "$(docker ps -a -q)")

		err := cmd.Start()

		gomega.Expect(err).Should(gomega.BeNil())
		err = cmd.Wait()
		gomega.Expect(err).Should(gomega.BeNil())
	})

	GinkgoIt("Should run a project execution in NodeJS ESM environment", func() {
		ginkgo.Skip("")

		environment := repository.NodeEsm
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.NodeEsm.Tag),
			},
		})).Should(gomega.BeNil())

		projectName := "project name node esm"
		root := testCreateFileStub(projectName, true, 1, false, nil, []string{})

		rootFile1 := testCreateFileStub("rootDirectoryFile1.mjs", false, 1, true, &root.Uuid, []string{})
		rootFile2 := testCreateFileStub("rootDirectoryFile2.mjs", false, 1, true, &root.Uuid, []string{})

		subDir := testCreateFileStub("subDir", false, 2, false, &root.Uuid, []string{})
		subDirFile1 := testCreateFileStub("subDirFile.mjs", false, 2, true, &subDir.Uuid, []string{})
		subSubDir := testCreateFileStub("subSubDir", false, 3, false, &subDir.Uuid, []string{})

		subSubDirFile := testCreateFileStub("subSubDirFile.mjs", false, 3, true, &subSubDir.Uuid, []string{})

		root.Children = append(root.Children, subDir.Uuid)
		root.Children = append(root.Children, rootFile1.Uuid)
		root.Children = append(root.Children, rootFile2.Uuid)
		subDir.Children = append(subDir.Children, subDirFile1.Uuid)
		subDir.Children = append(subDir.Children, subSubDir.Uuid)
		subSubDir.Children = append(subSubDir.Children, subSubDirFile.Uuid)

		codeProject := testCreateCodeProjectStub(projectName, "", []*repository.File{
			&root,
			&rootFile1,
			&rootFile2,
			&subDir,
			&subDirFile1,
			&subSubDir,
			&subSubDirFile,
		}, &root, &environment)

		content1 := testCreateFileContent(codeProject.Uuid, rootFile1.Uuid, `
import {execute} from './rootDirectoryFile2.mjs';
import {subDirDirFileExecute} from './subDir/subSubDir/subSubDirFile.mjs';

execute();

console.log('rootDirectoryFile1');
`)
		content2 := testCreateFileContent(codeProject.Uuid, rootFile2.Uuid, `
import {subDirFileExecute} from './subDir/subDirFile.mjs';

function execute() {
    console.log('rootDirectoryFile2');

    subDirFileExecute();
}

export { execute };
`)
		content3 := testCreateFileContent(codeProject.Uuid, subDirFile1.Uuid, `
function subDirFileExecute() {
    console.log('subDirFile');
}

export {subDirFileExecute};
`)
		content4 := testCreateFileContent(codeProject.Uuid, subSubDirFile.Uuid, `
function subDirDirFileExecute() {
    console.log('subSubDirFile');
}

export {subDirDirFileExecute}
`)

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "project",
			ExecutionType:     "project",
			EmulatorName:      string(environment.Name),
			EmulatorExtension: string(environment.Extension),
			EmulatorTag:       string(environment.Tag),
			EmulatorText:      "",
			PackageName:       "",
			CodeProject:       &codeProject,
			Contents: []*repository.FileContent{
				&content1,
				&content2,
				&content3,
				&content4,
			},
			ExecutingFile: &rootFile1,
		})

		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("rootDirectoryFile2\nsubDirFile\nrootDirectoryFile1\n"))

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should run a project execution in NodeJS latest environment", func() {
		ginkgo.Skip("")

		environment := repository.NodeLts
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.NodeLts.Tag),
			},
		})).Should(gomega.BeNil())

		projectName := "project name node"
		root := testCreateFileStub(projectName, true, 1, false, nil, []string{})

		rootFile1 := testCreateFileStub("rootDirectoryFile1.js", false, 1, true, &root.Uuid, []string{})
		rootFile2 := testCreateFileStub("rootDirectoryFile2.js", false, 1, true, &root.Uuid, []string{})

		subDir := testCreateFileStub("subDir", false, 2, false, &root.Uuid, []string{})
		subDirFile1 := testCreateFileStub("subDirFile.js", false, 2, true, &subDir.Uuid, []string{})
		subSubDir := testCreateFileStub("subSubDir", false, 3, false, &subDir.Uuid, []string{})

		subSubDirFile := testCreateFileStub("subSubDirFile.js", false, 3, true, &subSubDir.Uuid, []string{})

		root.Children = append(root.Children, subDir.Uuid)
		root.Children = append(root.Children, rootFile1.Uuid)
		root.Children = append(root.Children, rootFile2.Uuid)
		subDir.Children = append(subDir.Children, subDirFile1.Uuid)
		subDir.Children = append(subDir.Children, subSubDir.Uuid)
		subSubDir.Children = append(subSubDir.Children, subSubDirFile.Uuid)

		codeProject := testCreateCodeProjectStub(projectName, "", []*repository.File{
			&root,
			&rootFile1,
			&rootFile2,
			&subDir,
			&subDirFile1,
			&subSubDir,
			&subSubDirFile,
		}, &root, &environment)

		content1 := testCreateFileContent(codeProject.Uuid, rootFile1.Uuid, `
const {execute} = require('./rootDirectoryFile2');
const {subDirDirFileExecute} = require('./subDir/subSubDir/subSubDirFile');

execute();

console.log('rootDirectoryFile1');
`)
		content2 := testCreateFileContent(codeProject.Uuid, rootFile2.Uuid, `
const {subDirFileExecute} = require('./subDir/subDirFile');

function execute() {
    console.log('rootDirectoryFile2');

    subDirFileExecute();
}

module.exports = {
    execute: execute,
};
`)
		content3 := testCreateFileContent(codeProject.Uuid, subDirFile1.Uuid, `
function subDirFileExecute() {
    console.log('subDirFile');
}

module.exports = {
    subDirFileExecute: subDirFileExecute,
};
`)
		content4 := testCreateFileContent(codeProject.Uuid, subSubDirFile.Uuid, `
function subDirDirFileExecute() {
    console.log('subSubDirFile');
}

module.exports = {
    subDirDirFileExecute: subDirDirFileExecute,
};
`)

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "project",
			ExecutionType:     "project",
			EmulatorName:      string(environment.Name),
			EmulatorExtension: string(environment.Extension),
			EmulatorTag:       string(environment.Tag),
			EmulatorText:      "",
			PackageName:       "",
			CodeProject:       &codeProject,
			Contents: []*repository.FileContent{
				&content1,
				&content2,
				&content3,
				&content4,
			},
			ExecutingFile: &rootFile1,
		})

		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("rootDirectoryFile2\nsubDirFile\nrootDirectoryFile1\n"))

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should run a project execution in C# environment", func() {
		ginkgo.Skip("")

		environment := repository.CSharpMono
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.CSharpMono.Tag),
			},
		})).Should(gomega.BeNil())

		projectName := "project name node"
		root := testCreateFileStub(projectName, true, 1, false, nil, []string{})

		rootFile1 := testCreateFileStub("rootDirectoryFile1.cs", false, 1, true, &root.Uuid, []string{})

		subDir := testCreateFileStub("subDir", false, 2, false, &root.Uuid, []string{})
		subDirFile1 := testCreateFileStub("subDirFile.cs", false, 2, true, &subDir.Uuid, []string{})

		root.Children = append(root.Children, subDir.Uuid)
		root.Children = append(root.Children, rootFile1.Uuid)
		subDir.Children = append(subDir.Children, subDirFile1.Uuid)

		codeProject := testCreateCodeProjectStub(projectName, "", []*repository.File{
			&root,
			&rootFile1,
			&subDir,
			&subDirFile1,
		}, &root, &environment)

		content1 := testCreateFileContent(codeProject.Uuid, rootFile1.Uuid, `
public class HelloWorld
{
    public static void Main(string[] args)
    {
        NewClass v = new NewClass();
        v.Fn();
    }
}
`)
		content3 := testCreateFileContent(codeProject.Uuid, subDirFile1.Uuid, `
using System;

public class NewClass {
    public void Fn() {
                Console.WriteLine ("Hello World");
    }
}
`)

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "project",
			ExecutionType:     "project",
			EmulatorName:      string(environment.Name),
			EmulatorExtension: string(environment.Extension),
			EmulatorTag:       string(environment.Tag),
			EmulatorText:      "",
			PackageName:       "",
			CodeProject:       &codeProject,
			Contents: []*repository.FileContent{
				&content1,
				&content3,
			},
			ExecutingFile: &rootFile1,
		})

		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("Hello World\n"))

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should run a project execution in Go environment", func() {
		ginkgo.Skip("")

		environment := repository.GoLang
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.GoLang.Tag),
			},
		})).Should(gomega.BeNil())

		projectName := "project name node"
		root := testCreateFileStub(projectName, true, 1, false, nil, []string{})

		rootFile1 := testCreateFileStub("rootDirectoryFile1.go", false, 1, true, &root.Uuid, []string{})

		subDir := testCreateFileStub("subDir", false, 2, false, &root.Uuid, []string{})
		subDirFile1 := testCreateFileStub("subDirFile.go", false, 2, true, &subDir.Uuid, []string{})

		root.Children = append(root.Children, subDir.Uuid)
		root.Children = append(root.Children, rootFile1.Uuid)
		subDir.Children = append(subDir.Children, subDirFile1.Uuid)

		codeProject := testCreateCodeProjectStub(projectName, "mySuperPackage", []*repository.File{
			&root,
			&rootFile1,
			&subDir,
			&subDirFile1,
		}, &root, &environment)

		content1 := testCreateFileContent(codeProject.Uuid, rootFile1.Uuid, fmt.Sprintf(`
package main

import c "app/%s/subDir"

func main() {
    c.MyFunc()
}
`, codeProject.PackageName))

		content3 := testCreateFileContent(codeProject.Uuid, subDirFile1.Uuid, `
package subDir

import "fmt"

func MyFunc() {
	fmt.Println("Hello World")
}
`)

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "project",
			ExecutionType:     "project",
			EmulatorName:      string(environment.Name),
			EmulatorExtension: string(environment.Extension),
			EmulatorTag:       string(environment.Tag),
			EmulatorText:      "",
			PackageName:       codeProject.PackageName,
			CodeProject:       &codeProject,
			Contents: []*repository.FileContent{
				&content1,
				&content3,
			},
			ExecutingFile: &rootFile1,
		})

		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("Hello World\n"))
		gomega.Expect(result.Error).Should(gomega.BeNil())

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should run a project execution in Rust environment", func() {
		ginkgo.Skip("")

		environment := repository.Rust
		gomega.Expect(execution.Init(_var.PROJECT_EXECUTION, []execution.ContainerBlueprint{
			{
				WorkerNum:    1,
				ContainerNum: 1,
				Tag:          string(repository.Rust.Tag),
			},
		})).Should(gomega.BeNil())

		projectName := "project name rust"
		root := testCreateFileStub(projectName, true, 1, false, nil, []string{})

		rootFile1 := testCreateFileStub("main.rs", false, 1, true, &root.Uuid, []string{})

		subDir := testCreateFileStub("my_mod", false, 2, false, &root.Uuid, []string{})
		modRs := testCreateFileStub("mod.rs", false, 2, true, &subDir.Uuid, []string{})

		root.Children = append(root.Children, subDir.Uuid)
		root.Children = append(root.Children, rootFile1.Uuid)
		subDir.Children = append(subDir.Children, modRs.Uuid)

		codeProject := testCreateCodeProjectStub(projectName, "", []*repository.File{
			&root,
			&rootFile1,
			&subDir,
			&modRs,
		}, &root, &environment)

		content1 := testCreateFileContent(codeProject.Uuid, rootFile1.Uuid, `
mod my_mod;

fn main() {
	my_mod::my_func();
}
`)

		modContent := testCreateFileContent(codeProject.Uuid, modRs.Uuid, `
pub fn my_func() {
    println!("Hello World!");
}
`)

		result := execution.Service(_var.PROJECT_EXECUTION).RunJob(execution.Job{
			BuilderType:       "project",
			ExecutionType:     "project",
			EmulatorName:      string(environment.Name),
			EmulatorExtension: string(environment.Extension),
			EmulatorTag:       string(environment.Tag),
			EmulatorText:      "",
			PackageName:       codeProject.PackageName,
			CodeProject:       &codeProject,
			Contents: []*repository.FileContent{
				&content1,
				&modContent,
			},
			ExecutingFile: &rootFile1,
		})

		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("Hello World!\n"))
		gomega.Expect(result.Error).Should(gomega.BeNil())

		execution.Service(_var.PROJECT_EXECUTION).Close()
	})

	GinkgoIt("Should run a project execution as a session in a C environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()
		cp := testCreateCodeProject(activeSession, uuid.New().String(), repository.CLang)
		cpUuid := cp["uuid"].(string)

		var rootDirectory map[string]interface{}
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		rootDirectoryFile1 := testCreateFile(activeSession, true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile1.c")
		testUpdateFileContent(activeSession, cpUuid, rootDirectoryFile1["uuid"].(string), fmt.Sprintf(`

`))

		testCreateFile(activeSession, true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile2.c")

		sessionUuid := testCreateProjectTemporarySession(repository.ActiveSession{}, cp["uuid"].(string), rootDirectoryFile1["uuid"].(string))
		bm := map[string]interface{}{
			"uuid":     sessionUuid,
			"fileUuid": rootDirectoryFile1["uuid"].(string),
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/project", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executeProjectHandler)

		handler.ServeHTTP(rr, req)

		b := rr.Body.Bytes()

		var apiResponse httpUtil.ApiResponse
		err = json.Unmarshal(b, &apiResponse)

		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
		gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
		gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_RESOURCE))
		gomega.Expect(apiResponse.Message).To(gomega.Equal("Emulator run result"))
		gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Code).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusOK))
		gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

		b, err = json.Marshal(apiResponse.Data)

		gomega.Expect(err).To(gomega.BeNil())

		var result repository.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("Hello world!\n"))
	})

	GinkgoIt("Should run a project execution as a session in a C++ environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()
		cp := testCreateCodeProject(activeSession, uuid.New().String(), repository.CPlus)
		cpUuid := cp["uuid"].(string)

		var rootDirectory map[string]interface{}
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		rootDirectoryFile1 := testCreateFile(activeSession, true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile1.cpp")
		testUpdateFileContent(activeSession, cpUuid, rootDirectoryFile1["uuid"].(string), fmt.Sprintf(`

`))

		testCreateFile(activeSession, true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile2.cpp")

		sessionUuid := testCreateProjectTemporarySession(repository.ActiveSession{}, cp["uuid"].(string), rootDirectoryFile1["uuid"].(string))
		bm := map[string]interface{}{
			"uuid":     sessionUuid,
			"fileUuid": rootDirectoryFile1["uuid"].(string),
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/project", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executeProjectHandler)

		handler.ServeHTTP(rr, req)

		b := rr.Body.Bytes()

		var apiResponse httpUtil.ApiResponse
		err = json.Unmarshal(b, &apiResponse)

		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
		gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
		gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_RESOURCE))
		gomega.Expect(apiResponse.Message).To(gomega.Equal("Emulator run result"))
		gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Code).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusOK))
		gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

		b, err = json.Marshal(apiResponse.Data)

		gomega.Expect(err).To(gomega.BeNil())

		var result repository.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("Hello World!"))
	})

	GinkgoIt("Should run a project execution as a session in a Haskell environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()
		cp := testCreateCodeProject(activeSession, uuid.New().String(), repository.Haskell)
		cpUuid := cp["uuid"].(string)

		var rootDirectory *repository.File
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		testUpdateFileContent(activeSession, cpUuid, rootDirectory.Children[0], fmt.Sprintf(`
import Foo
import Bar.FooBar

main = putStrLn "Hello, World!"
`))

		rootDirectoryFile1 := testCreateFile(activeSession, true, rootDirectory.Uuid, cpUuid, "Foo.hs")
		testUpdateFileContent(activeSession, cpUuid, rootDirectoryFile1["uuid"].(string), fmt.Sprintf(`
module Foo where
`))

		barDir := testCreateFile(activeSession, false, rootDirectory.Uuid, cpUuid, "Bar")

		fooBar := testCreateFile(activeSession, true, barDir["uuid"].(string), cpUuid, "FooBar.hs")
		testUpdateFileContent(activeSession, cpUuid, fooBar["uuid"].(string), fmt.Sprintf(`
module Bar.FooBar where
`))

		sessionUuid := testCreateProjectTemporarySession(repository.ActiveSession{}, cp["uuid"].(string), rootDirectoryFile1["uuid"].(string))
		bm := map[string]interface{}{
			"uuid":     sessionUuid,
			"fileUuid": rootDirectoryFile1["uuid"].(string),
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/project", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executeProjectHandler)

		handler.ServeHTTP(rr, req)

		b := rr.Body.Bytes()

		var apiResponse httpUtil.ApiResponse
		err = json.Unmarshal(b, &apiResponse)

		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
		gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
		gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_RESOURCE))
		gomega.Expect(apiResponse.Message).To(gomega.Equal("Emulator run result"))
		gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Code).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusOK))
		gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

		b, err = json.Marshal(apiResponse.Data)

		gomega.Expect(err).To(gomega.BeNil())

		var result repository.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("\nHello, World!\n"))
	})

	GinkgoIt("Should run a project execution as a session in a Ruby environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()
		cp := testCreateCodeProject(activeSession, uuid.New().String(), repository.Ruby)
		cpUuid := cp["uuid"].(string)

		var rootDirectory *repository.File
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		foo := testCreateFile(activeSession, true, rootDirectory.Uuid, cpUuid, "foo.rb")
		testUpdateFileContent(activeSession, cpUuid, foo["uuid"].(string), fmt.Sprintf(`
class TestClass
    def initialize
        puts "TestClass object created"
    end
end 
`))

		bar := testCreateFile(activeSession, true, rootDirectory.Uuid, cpUuid, "bar.rb")
		testUpdateFileContent(activeSession, cpUuid, bar["uuid"].(string), fmt.Sprintf(`
require "./foo.rb"

puts "Hello world!"
`))

		sessionUuid := testCreateProjectTemporarySession(repository.ActiveSession{}, cp["uuid"].(string), bar["uuid"].(string))
		bm := map[string]interface{}{
			"uuid":     sessionUuid,
			"fileUuid": bar["uuid"].(string),
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/project", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executeProjectHandler)

		handler.ServeHTTP(rr, req)

		b := rr.Body.Bytes()

		var apiResponse httpUtil.ApiResponse
		err = json.Unmarshal(b, &apiResponse)

		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
		gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
		gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_RESOURCE))
		gomega.Expect(apiResponse.Message).To(gomega.Equal("Emulator run result"))
		gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Code).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusOK))
		gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

		b, err = json.Marshal(apiResponse.Data)

		gomega.Expect(err).To(gomega.BeNil())

		var result repository.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("Hello world!\n"))
	})

	GinkgoIt("Should run a project execution as a session in a PHP environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()
		cp := testCreateCodeProject(activeSession, uuid.New().String(), repository.Php74)
		cpUuid := cp["uuid"].(string)

		var rootDirectory *repository.File
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		foo := testCreateFile(activeSession, true, rootDirectory.Uuid, cpUuid, "foo.php")
		testUpdateFileContent(activeSession, cpUuid, foo["uuid"].(string), fmt.Sprintf(`
`))

		bar := testCreateFile(activeSession, true, rootDirectory.Uuid, cpUuid, "bar.php")
		testUpdateFileContent(activeSession, cpUuid, bar["uuid"].(string), fmt.Sprintf(`
<?php

require(__DIR__."/foo.php");

echo "Hello world!";
`))

		sessionUuid := testCreateProjectTemporarySession(repository.ActiveSession{}, cp["uuid"].(string), bar["uuid"].(string))
		bm := map[string]interface{}{
			"uuid":     sessionUuid,
			"fileUuid": bar["uuid"].(string),
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/project", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executeProjectHandler)

		handler.ServeHTTP(rr, req)

		b := rr.Body.Bytes()

		var apiResponse httpUtil.ApiResponse
		err = json.Unmarshal(b, &apiResponse)

		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
		gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
		gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_RESOURCE))
		gomega.Expect(apiResponse.Message).To(gomega.Equal("Emulator run result"))
		gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Code).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusOK))
		gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

		b, err = json.Marshal(apiResponse.Data)

		gomega.Expect(err).To(gomega.BeNil())

		var result repository.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("\n\nHello world!"))
	})

	GinkgoIt("Should run a project execution as a session in a Python2 environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()
		cp := testCreateCodeProject(activeSession, uuid.New().String(), repository.Python2)
		cpUuid := cp["uuid"].(string)

		var rootDirectory *repository.File
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		foo := testCreateFile(activeSession, true, rootDirectory.Uuid, cpUuid, "foo.py")
		testUpdateFileContent(activeSession, cpUuid, foo["uuid"].(string), fmt.Sprintf(`
import foo.bar as bt

bt.greeting("Jonathan")
`))

		foobar := testCreateFile(activeSession, false, rootDirectory.Uuid, cpUuid, "foo")
		testCreateFile(activeSession, true, foobar["uuid"].(string), cpUuid, "__init__.py")

		bar := testCreateFile(activeSession, true, foobar["uuid"].(string), cpUuid, "bar.py")
		testUpdateFileContent(activeSession, cpUuid, bar["uuid"].(string), fmt.Sprintf(`
def greeting(name):
  print("Hello, " + name) 
`))

		sessionUuid := testCreateProjectTemporarySession(repository.ActiveSession{}, cp["uuid"].(string), foo["uuid"].(string))
		bm := map[string]interface{}{
			"uuid":     sessionUuid,
			"fileUuid": foo["uuid"].(string),
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/project", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executeProjectHandler)

		handler.ServeHTTP(rr, req)

		b := rr.Body.Bytes()

		var apiResponse httpUtil.ApiResponse
		err = json.Unmarshal(b, &apiResponse)

		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
		gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
		gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_RESOURCE))
		gomega.Expect(apiResponse.Message).To(gomega.Equal("Emulator run result"))
		gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Code).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusOK))
		gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

		b, err = json.Marshal(apiResponse.Data)

		gomega.Expect(err).To(gomega.BeNil())

		var result repository.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("Hello, Jonathan\n"))
	})

	GinkgoIt("Should run a project execution as a session in a Python3 environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()
		cp := testCreateCodeProject(activeSession, uuid.New().String(), repository.Python3)
		cpUuid := cp["uuid"].(string)

		var rootDirectory *repository.File
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		foo := testCreateFile(activeSession, true, rootDirectory.Uuid, cpUuid, "foo.py")
		testUpdateFileContent(activeSession, cpUuid, foo["uuid"].(string), fmt.Sprintf(`
import foo.bar as bt

bt.greeting("Jonathan")
`))

		foobar := testCreateFile(activeSession, false, rootDirectory.Uuid, cpUuid, "foo")
		testCreateFile(activeSession, true, foobar["uuid"].(string), cpUuid, "__init__.py")

		bar := testCreateFile(activeSession, true, foobar["uuid"].(string), cpUuid, "bar.py")
		testUpdateFileContent(activeSession, cpUuid, bar["uuid"].(string), fmt.Sprintf(`
def greeting(name):
  print("Hello, " + name) 
`))

		sessionUuid := testCreateProjectTemporarySession(repository.ActiveSession{}, cp["uuid"].(string), foo["uuid"].(string))
		bm := map[string]interface{}{
			"uuid":     sessionUuid,
			"fileUuid": foo["uuid"].(string),
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/project", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executeProjectHandler)

		handler.ServeHTTP(rr, req)

		b := rr.Body.Bytes()

		var apiResponse httpUtil.ApiResponse
		err = json.Unmarshal(b, &apiResponse)

		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusOK))
		gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
		gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_RESOURCE))
		gomega.Expect(apiResponse.Message).To(gomega.Equal("Emulator run result"))
		gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Code).To(gomega.Equal(0))
		gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusOK))
		gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

		b, err = json.Marshal(apiResponse.Data)

		gomega.Expect(err).To(gomega.BeNil())

		var result repository.RunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("Hello, Jonathan\n"))
	})
})
