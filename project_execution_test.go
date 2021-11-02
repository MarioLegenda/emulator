package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"therebelsource/emulator/httpUtil"
	"therebelsource/emulator/runner"
	"therebelsource/emulator/staticTypes"
)

var _ = GinkgoDescribe("Project execution tests", func() {
	GinkgoIt("Should run a project execution as a session in a node14 environment", func() {
		testPrepare()
		defer testCleanup()

		cp := testCreateCodeProject(runner.Node14)
		cpUuid := cp["uuid"].(string)

		var rootDirectory map[string]interface{}
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		rootDirectoryFile1 := testCreateFile(true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile1.js")
		testUpdateFileContent(cpUuid, rootDirectoryFile1["uuid"].(string), fmt.Sprintf(`
const {execute} = require('./rootDirectoryFile2');
const {subDirDirFileExecute} = require('./subDir/subSubDir/subSubDirFile');

execute();

console.log('rootDirectoryFile1');
`))

		rootDirectoryFile2 := testCreateFile(true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile2.js")
		testUpdateFileContent(cpUuid, rootDirectoryFile2["uuid"].(string), `
const {subDirFileExecute} = require('./subDir/subDirFile');

function execute() {
    console.log('rootDirectoryFile2');

    subDirFileExecute();
}

module.exports = {
	execute,
}
`)

		rootDirectorySubDir := testCreateFile(false, rootDirectory["uuid"].(string), cpUuid, "subDir")
		subDirFile := testCreateFile(true, rootDirectorySubDir["uuid"].(string), cpUuid, "subDirFile.js")
		testUpdateFileContent(cpUuid, subDirFile["uuid"].(string), `
function subDirFileExecute() {
    console.log('subDirFile');
}

module.exports = {
	subDirFileExecute,
}
`)

		subDir := testCreateFile(false, rootDirectorySubDir["uuid"].(string), cpUuid, "subSubDir")
		subDirSubFile := testCreateFile(true, subDir["uuid"].(string), cpUuid, "subSubDirFile.js")
		testUpdateFileContent(cpUuid, subDirSubFile["uuid"].(string), `
function subDirDirFileExecute() {
    console.log('subSubDirFile');
}

module.exports = {
	subDirDirFileExecute,
}
`)

		bm := map[string]interface{}{
			"codeProjectUuid": cpUuid,
			"fileUuid": rootDirectoryFile1["uuid"].(string),
			"type": "session",
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

	GinkgoIt("Should run a project execution as a session in a Go environment", func() {
		testPrepare()
		defer testCleanup()

		cp := testCreateCodeProject(runner.GoLang)
		cpUuid := cp["uuid"].(string)

		var rootDirectory map[string]interface{}
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		rootDirectoryFile1 := testCreateFile(true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile1.go")
		testUpdateFileContent(cpUuid, rootDirectoryFile1["uuid"].(string), fmt.Sprintf(`
	package main

`))

		bm := map[string]interface{}{
			"codeProjectUuid": cpUuid,
			"fileUuid": rootDirectoryFile1["uuid"].(string),
			"type": "session",
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
		gomega.Expect(result.Result).Should(gomega.Equal("Hello world\r\n"))
	})

	GinkgoIt("Should run a project execution as a session in a Rust environment", func() {
		testPrepare()
		defer testCleanup()

		cp := testCreateCodeProject(runner.Rust)
		cpUuid := cp["uuid"].(string)

		var rootDirectory map[string]interface{}
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		rootDirectoryFile1 := testCreateFile(true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile1.rs")
		testUpdateFileContent(cpUuid, rootDirectoryFile1["uuid"].(string), fmt.Sprintf(`
`))

		testCreateFile(true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile2.rs")

		bm := map[string]interface{}{
			"codeProjectUuid": cpUuid,
			"fileUuid": rootDirectoryFile1["uuid"].(string),
			"type": "session",
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
		gomega.Expect(result.Result).Should(gomega.Equal("Hello World!\r\n"))
	})

   GinkgoIt("Should run a project execution as a session in a C environment", func() {
		testPrepare()
		defer testCleanup()

		cp := testCreateCodeProject(runner.CLang)
		cpUuid := cp["uuid"].(string)

		var rootDirectory map[string]interface{}
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		rootDirectoryFile1 := testCreateFile(true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile1.c")
		testUpdateFileContent(cpUuid, rootDirectoryFile1["uuid"].(string), fmt.Sprintf(`

`))

		testCreateFile(true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile2.c")

		bm := map[string]interface{}{
			"codeProjectUuid": cpUuid,
			"fileUuid": rootDirectoryFile1["uuid"].(string),
			"type": "session",
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
})

