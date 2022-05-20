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
	"therebelsource/emulator/httpUtil"
	"therebelsource/emulator/repository"
	"therebelsource/emulator/runner"
	"therebelsource/emulator/staticTypes"
)

var _ = GinkgoDescribe("Project execution tests", func() {
	GinkgoIt("Should run a project execution as a session in a NodeJS ESM environment", func() {
		activeSession := testCreateAccount()

		cp := testCreateCodeProject(activeSession, uuid.New().String(), runner.NodeEsm)

		cpUuid := cp["uuid"].(string)

		var rootDirectory map[string]interface{}
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		rootDirectoryFile1 := testCreateFile(activeSession, true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile1.mjs")
		testUpdateFileContent(activeSession, cpUuid, rootDirectoryFile1["uuid"].(string), fmt.Sprintf(`
import {execute} from './rootDirectoryFile2.mjs';
import {subDirDirFileExecute} from './subDir/subSubDir/subSubDirFile.mjs';

execute();

console.log('rootDirectoryFile1');
`))

		rootDirectoryFile2 := testCreateFile(activeSession, true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile2.mjs")
		testUpdateFileContent(activeSession, cpUuid, rootDirectoryFile2["uuid"].(string), `
import {subDirFileExecute} from './subDir/subDirFile.mjs';

function execute() {
    console.log('rootDirectoryFile2');

    subDirFileExecute();
}

export { execute };
`)

		rootDirectorySubDir := testCreateFile(activeSession, false, rootDirectory["uuid"].(string), cpUuid, "subDir")
		subDirFile := testCreateFile(activeSession, true, rootDirectorySubDir["uuid"].(string), cpUuid, "subDirFile.mjs")
		testUpdateFileContent(activeSession, cpUuid, subDirFile["uuid"].(string), `
function subDirFileExecute() {
    console.log('subDirFile');
}

export {subDirFileExecute};
`)

		subDir := testCreateFile(activeSession, false, rootDirectorySubDir["uuid"].(string), cpUuid, "subSubDir")
		subDirSubFile := testCreateFile(activeSession, true, subDir["uuid"].(string), cpUuid, "subSubDirFile.mjs")
		testUpdateFileContent(activeSession, cpUuid, subDirSubFile["uuid"].(string), `
function subDirDirFileExecute() {
    console.log('subSubDirFile');
}

export {subDirDirFileExecute}
`)

		sessionUuid := testCreateProjectTemporarySession(repository.ActiveSession{}, cp["uuid"].(string), rootDirectoryFile1["uuid"].(string))

		bm := map[string]interface{}{
			"uuid": sessionUuid,
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

		var result runner.ProjectRunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("rootDirectoryFile2\nsubDirFile\nrootDirectoryFile1\n"))
	})

	GinkgoIt("Should run a project execution as a session in a NodeJS ESM environment with a file in a deeper directory structure", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()
		cp := testCreateCodeProject(activeSession, uuid.New().String(), runner.NodeEsm)
		cpUuid := cp["uuid"].(string)

		var rootDirectory map[string]interface{}
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		rootDirectoryFile1 := testCreateFile(activeSession, true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile1.mjs")
		testUpdateFileContent(activeSession, cpUuid, rootDirectoryFile1["uuid"].(string), fmt.Sprintf(`
import {execute} from './rootDirectoryFile2.mjs';
import './subDir/subSubDir/subSubDirFile.mjs';

execute();

console.log('rootDirectoryFile1');
`))

		rootDirectoryFile2 := testCreateFile(activeSession, true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile2.mjs")
		testUpdateFileContent(activeSession, cpUuid, rootDirectoryFile2["uuid"].(string), `
import {subDirFileExecute} from './subDir/subDirFile.mjs';

function execute() {
    console.log('rootDirectoryFile2');

    subDirFileExecute();
}

export { execute };
`)

		rootDirectorySubDir := testCreateFile(activeSession, false, rootDirectory["uuid"].(string), cpUuid, "subDir")
		subDirFile := testCreateFile(activeSession, true, rootDirectorySubDir["uuid"].(string), cpUuid, "subDirFile.mjs")
		testUpdateFileContent(activeSession, cpUuid, subDirFile["uuid"].(string), `
function subDirFileExecute() {
    console.log('subDirFile');
}

export { subDirFileExecute }
`)

		subDir := testCreateFile(activeSession, false, rootDirectorySubDir["uuid"].(string), cpUuid, "subSubDir")
		subDirSubFile := testCreateFile(activeSession, true, subDir["uuid"].(string), cpUuid, "subSubDirFile.mjs")
		testUpdateFileContent(activeSession, cpUuid, subDirSubFile["uuid"].(string), `
console.log('subSubDirFile.js is executed');
`)

		sessionUuid := testCreateProjectTemporarySession(repository.ActiveSession{}, cp["uuid"].(string), rootDirectoryFile1["uuid"].(string))
		bm := map[string]interface{}{
			"uuid":     sessionUuid,
			"fileUuid": subDirSubFile["uuid"].(string),
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

		var result runner.ProjectRunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("subSubDirFile.js is executed\nrootDirectoryFile2\nsubDirFile\nrootDirectoryFile1\n"))
	})

	GinkgoIt("Should run a project execution as a session in a NodeJS ESM environment with a file in a deeper directory structure", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()
		cp := testCreateCodeProject(activeSession, uuid.New().String(), runner.CSharpMono)
		cpUuid := cp["uuid"].(string)

		var rootDirectory map[string]interface{}
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		rootDirectoryFile1 := testCreateFile(activeSession, true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile1.cs")
		testUpdateFileContent(activeSession, cpUuid, rootDirectoryFile1["uuid"].(string), fmt.Sprintf(`
public class HelloWorld
{
    public static void Main(string[] args)
    {
        NewClass v = new NewClass();
        v.Fn();
    }
}
`))

		rootDirectorySubDir := testCreateFile(activeSession, false, rootDirectory["uuid"].(string), cpUuid, "subDir")
		subDirFile := testCreateFile(activeSession, true, rootDirectorySubDir["uuid"].(string), cpUuid, "subDirFile.cs")
		testUpdateFileContent(activeSession, cpUuid, subDirFile["uuid"].(string), `
using System;

public class NewClass {
    public void Fn() {
                Console.WriteLine ("Hello World");
    }
}
`)

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

		var result runner.ProjectRunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("Hello World\n"))
	})

	GinkgoIt("Should run a project execution as a session in a Go environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()
		cp := testCreateCodeProject(activeSession, uuid.New().String(), runner.GoLang)
		cpUuid := cp["uuid"].(string)

		var rootDirectory map[string]interface{}
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		rootDirectoryFile1 := testCreateFile(activeSession, true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile1.go")
		testUpdateFileContent(activeSession, cpUuid, rootDirectoryFile1["uuid"].(string), fmt.Sprintf(`
	package main

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

		var result runner.SingleFileRunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("Hello world\n"))
	})

	GinkgoIt("Should run a project execution as a session in a Go environment with a package", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()
		cp := testCreateCodeProject(activeSession, uuid.New().String(), runner.GoLang)
		cpUuid := cp["uuid"].(string)

		var rootDirectory map[string]interface{}
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		packageName := testCreateFile(activeSession, false, rootDirectory["uuid"].(string), cpUuid, "somePackage")
		packageFile := testCreateFile(activeSession, true, packageName["uuid"].(string), cpUuid, "someFunc.go")
		testUpdateFileContent(activeSession, cpUuid, packageFile["uuid"].(string), `
package somePackage

import "fmt"

func MyFunc() {
	fmt.Println("Output from somePackage")
}
`)

		rootDirectoryChildren := rootDirectory["children"].([]interface{})
		testUpdateFileContent(activeSession, cpUuid, rootDirectoryChildren[0].(string), fmt.Sprintf(`
package main

import pcg "%s/somePackage"
import "fmt"

func main() {
    fmt.Println("Hello world")

	pcg.MyFunc()
}
`, cp["packageName"].(string)))

		sessionUuid := testCreateProjectTemporarySession(repository.ActiveSession{}, cp["uuid"].(string), rootDirectoryChildren[0].(string))
		bm := map[string]interface{}{
			"uuid":     sessionUuid,
			"fileUuid": rootDirectoryChildren[0].(string),
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

		var result runner.SingleFileRunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("Hello world\nOutput from somePackage\n"))
	})

	GinkgoIt("Should run a project execution as a session in a Go environment multiple times with the same code project", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()
		cp := testCreateCodeProject(activeSession, uuid.New().String(), runner.GoLang)
		cpUuid := cp["uuid"].(string)

		var rootDirectory map[string]interface{}
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		rootDirectoryFile1 := testCreateFile(activeSession, true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile1.go")
		testUpdateFileContent(activeSession, cpUuid, rootDirectoryFile1["uuid"].(string), fmt.Sprintf(`
	package main

`))

		for i := 0; i < 5; i++ {
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

			var result runner.SingleFileRunResult
			gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

			gomega.Expect(result.Timeout).Should(gomega.Equal(5))
			gomega.Expect(result.Success).Should(gomega.BeTrue())
			gomega.Expect(result.Result).Should(gomega.Equal("Hello world\n"))
		}
	})

	GinkgoIt("Should run a project execution as a session in a Rust environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()
		cp := testCreateCodeProject(activeSession, uuid.New().String(), runner.Rust)
		cpUuid := cp["uuid"].(string)

		var rootDirectory map[string]interface{}
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		rootDirectoryFile1 := testCreateFile(activeSession, true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile1.rs")
		testUpdateFileContent(activeSession, cpUuid, rootDirectoryFile1["uuid"].(string), fmt.Sprintf(`
`))

		testCreateFile(activeSession, true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile2.rs")

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

		var result runner.SingleFileRunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("Hello World!\n"))
	})

	GinkgoIt("Should run a project execution as a session in a C environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()
		cp := testCreateCodeProject(activeSession, uuid.New().String(), runner.CLang)
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

		var result runner.SingleFileRunResult
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
		cp := testCreateCodeProject(activeSession, uuid.New().String(), runner.CPlus)
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

		var result runner.SingleFileRunResult
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
		cp := testCreateCodeProject(activeSession, uuid.New().String(), runner.Haskell)
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

		var result runner.SingleFileRunResult
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
		cp := testCreateCodeProject(activeSession, uuid.New().String(), runner.Ruby)
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

		var result runner.SingleFileRunResult
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
		cp := testCreateCodeProject(activeSession, uuid.New().String(), runner.Php74)
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

		var result runner.SingleFileRunResult
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
		cp := testCreateCodeProject(activeSession, uuid.New().String(), runner.Python2)
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

		var result runner.SingleFileRunResult
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
		cp := testCreateCodeProject(activeSession, uuid.New().String(), runner.Python3)
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

		var result runner.SingleFileRunResult
		gomega.Expect(json.Unmarshal(b, &result)).To(gomega.BeNil())

		gomega.Expect(result.Timeout).Should(gomega.Equal(5))
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		gomega.Expect(result.Result).Should(gomega.Equal("Hello, Jonathan\n"))
	})
})
