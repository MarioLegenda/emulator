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

var _ = GinkgoDescribe("Linked project execution tests", func() {
	GinkgoIt("Should run a linked code block execution as a session in a C environment", func() {
		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()
		link := testCreateBlog(activeSession)
		pageUuid := link["page"].(map[string]interface{})["uuid"].(string)
		blogUuid := link["blog"].(map[string]interface{})["uuid"].(string)

		cp := testCreateCodeProject(activeSession, uuid.New().String(), runner.CLang)

		cpUuid := cp["uuid"].(string)
		cb := testCreateCodeBlock(pageUuid, activeSession)
		testUpdateCodeBlock(activeSession, pageUuid, cb["uuid"].(string), `
#include <stdio.h>
int main() {
   printf("Hello world!\n");
   return 0;
}
`)
		testLinkCodeProject(activeSession, cp["uuid"].(string), pageUuid, cb["uuid"].(string), blogUuid)

		var rootDirectory map[string]interface{}
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		testUpdateFileContent(activeSession, cpUuid, rootDirectory["children"].([]interface{})[0].(string), fmt.Sprintf(`
`))

		rootDirectoryFile1 := testCreateFile(activeSession, true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile1.c")
		testUpdateFileContent(activeSession, cpUuid, rootDirectoryFile1["uuid"].(string), fmt.Sprintf(`
`))

		testCreateFile(activeSession, true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile2.c")

		sessionUuid := testCreateLinkedSession(activeSession, pageUuid, cb["uuid"].(string))
		bm := map[string]interface{}{
			"uuid": sessionUuid,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/linked-code-block", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executeLinkedProjectHandler)

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

	GinkgoIt("Should run a linked code block execution as a session in a C++ environment", func() {
		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()
		link := testCreateBlog(activeSession)
		pageUuid := link["page"].(map[string]interface{})["uuid"].(string)
		blogUuid := link["blog"].(map[string]interface{})["uuid"].(string)

		cp := testCreateCodeProject(activeSession, uuid.New().String(), runner.CPlus)
		cpUuid := cp["uuid"].(string)
		cb := testCreateCodeBlock(pageUuid, activeSession)
		testUpdateCodeBlock(activeSession, pageUuid, cb["uuid"].(string), `
#include <iostream>

int main() {
    std::cout << "Hello World!";
    return 0;
}
`)
		testLinkCodeProject(activeSession, cp["uuid"].(string), pageUuid, cb["uuid"].(string), blogUuid)

		var rootDirectory map[string]interface{}
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		testUpdateFileContent(activeSession, cpUuid, rootDirectory["children"].([]interface{})[0].(string), fmt.Sprintf(`
`))

		rootDirectoryFile1 := testCreateFile(activeSession, true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile1.cpp")
		testUpdateFileContent(activeSession, cpUuid, rootDirectoryFile1["uuid"].(string), fmt.Sprintf(`
`))

		testCreateFile(activeSession, true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile2.cpp")

		sessionUuid := testCreateLinkedSession(activeSession, pageUuid, cb["uuid"].(string))
		bm := map[string]interface{}{
			"uuid": sessionUuid,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/linked-code-block", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executeLinkedProjectHandler)

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

	GinkgoIt("Should run a linked code block execution as a session in a Haskell environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()
		link := testCreateBlog(activeSession)
		pageUuid := link["page"].(map[string]interface{})["uuid"].(string)
		blogUuid := link["blog"].(map[string]interface{})["uuid"].(string)

		cp := testCreateCodeProject(activeSession, uuid.New().String(), runner.Haskell)
		cpUuid := cp["uuid"].(string)
		cb := testCreateCodeBlock(pageUuid, activeSession)
		testUpdateCodeBlock(activeSession, pageUuid, cb["uuid"].(string), `
import Foo
import Bar.FooBar

main = putStrLn "Hello, World!"
`)
		testLinkCodeProject(activeSession, cp["uuid"].(string), pageUuid, cb["uuid"].(string), blogUuid)

		var rootDirectory *repository.File
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		testUpdateFileContent(activeSession, cpUuid, rootDirectory.Children[0], fmt.Sprintf(`
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

		sessionUuid := testCreateLinkedSession(activeSession, pageUuid, cb["uuid"].(string))
		bm := map[string]interface{}{
			"uuid": sessionUuid,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/linked-code-block", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executeLinkedProjectHandler)

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

	GinkgoIt("Should run a linked code block execution as a session in a Go environment", func() {
		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()
		link := testCreateBlog(activeSession)
		pageUuid := link["page"].(map[string]interface{})["uuid"].(string)
		blogUuid := link["blog"].(map[string]interface{})["uuid"].(string)

		cp := testCreateCodeProject(activeSession, "my_cool_name", runner.GoLang)
		cpUuid := cp["uuid"].(string)
		cb := testCreateCodeBlock(pageUuid, activeSession)
		cbLink := testLinkCodeProject(activeSession, cp["uuid"].(string), pageUuid, cb["uuid"].(string), blogUuid)

		testUpdateCodeBlock(activeSession, pageUuid, cb["uuid"].(string), fmt.Sprintf(`
package main

import c "%s/%s"

func main() {
    c.ExecuteFn()
}
`, cbLink["packageName"].(string), "myPackage"))

		var rootDirectory *repository.File
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		testUpdateFileContent(activeSession, cpUuid, rootDirectory.Children[0], fmt.Sprintf(fmt.Sprintf(`
package main
`)))

		packageDir := testCreateFile(activeSession, false, rootDirectory.Uuid, cpUuid, "myPackage")
		rootDirectoryFile1 := testCreateFile(activeSession, true, packageDir["uuid"].(string), cpUuid, "rootDirectoryFile1.go")
		testUpdateFileContent(activeSession, cpUuid, rootDirectoryFile1["uuid"].(string), fmt.Sprintf(`
package myPackage

import "fmt"

func ExecuteFn() {
    fmt.Println("Executing fn")
}
`))

		sessionUuid := testCreateLinkedSession(activeSession, pageUuid, cb["uuid"].(string))

		bm := map[string]interface{}{
			"uuid": sessionUuid,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/linked-code-block", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executeLinkedProjectHandler)

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
		gomega.Expect(result.Result).Should(gomega.Equal("Executing fn\n"))
	})

	GinkgoIt("Should run a linked code block execution as a session in a NodeJS environment", func() {
		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()
		link := testCreateBlog(activeSession)
		pageUuid := link["page"].(map[string]interface{})["uuid"].(string)
		blogUuid := link["blog"].(map[string]interface{})["uuid"].(string)

		cp := testCreateCodeProject(activeSession, uuid.New().String(), runner.NodeLts)
		cpUuid := cp["uuid"].(string)
		cb := testCreateCodeBlock(pageUuid, activeSession)
		testUpdateCodeBlock(activeSession, pageUuid, cb["uuid"].(string), `
require('./rootDirectoryFile1.js');
`)
		testLinkCodeProject(activeSession, cp["uuid"].(string), pageUuid, cb["uuid"].(string), blogUuid)

		var rootDirectory map[string]interface{}
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		rootDirectoryFile1 := testCreateFile(activeSession, true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile1.js")
		testUpdateFileContent(activeSession, cpUuid, rootDirectoryFile1["uuid"].(string), fmt.Sprintf(`
const {execute} = require('./rootDirectoryFile2');
const {subDirDirFileExecute} = require('./subDir/subSubDir/subSubDirFile');

execute();

console.log('rootDirectoryFile1');
`))

		rootDirectoryFile2 := testCreateFile(activeSession, true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile2.js")
		testUpdateFileContent(activeSession, cpUuid, rootDirectoryFile2["uuid"].(string), `
const {subDirFileExecute} = require('./subDir/subDirFile');

function execute() {
    console.log('rootDirectoryFile2');

    subDirFileExecute();
}

module.exports = {
	execute,
}
`)

		rootDirectorySubDir := testCreateFile(activeSession, false, rootDirectory["uuid"].(string), cpUuid, "subDir")
		subDirFile := testCreateFile(activeSession, true, rootDirectorySubDir["uuid"].(string), cpUuid, "subDirFile.js")
		testUpdateFileContent(activeSession, cpUuid, subDirFile["uuid"].(string), `
function subDirFileExecute() {
    console.log('subDirFile');
}

module.exports = {
	subDirFileExecute,
}
`)

		subDir := testCreateFile(activeSession, false, rootDirectorySubDir["uuid"].(string), cpUuid, "subSubDir")
		subDirSubFile := testCreateFile(activeSession, true, subDir["uuid"].(string), cpUuid, "subSubDirFile.js")
		testUpdateFileContent(activeSession, cpUuid, subDirSubFile["uuid"].(string), `
function subDirDirFileExecute() {
    console.log('subSubDirFile');
}

module.exports = {
	subDirDirFileExecute,
}
`)
		sessionUuid := testCreateLinkedSession(activeSession, pageUuid, cb["uuid"].(string))

		bm := map[string]interface{}{
			"uuid": sessionUuid,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/linked-code-block", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executeLinkedProjectHandler)

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
		gomega.Expect(result.Result).Should(gomega.Equal("rootDirectoryFile2\nsubDirFile\nrootDirectoryFile1\n"))
	})

	GinkgoIt("Should run a linked code block execution as a session in a NodeJS ESM environment", func() {
		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()
		link := testCreateBlog(activeSession)
		pageUuid := link["page"].(map[string]interface{})["uuid"].(string)
		blogUuid := link["blog"].(map[string]interface{})["uuid"].(string)

		cp := testCreateCodeProject(activeSession, uuid.New().String(), runner.NodeEsm)
		cpUuid := cp["uuid"].(string)
		cb := testCreateCodeBlock(pageUuid, activeSession)
		testUpdateCodeBlock(activeSession, pageUuid, cb["uuid"].(string), `
import './rootDirectoryFile1.mjs';
`)
		testLinkCodeProject(activeSession, cp["uuid"].(string), pageUuid, cb["uuid"].(string), blogUuid)

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

export { subDirFileExecute };
`)

		subDir := testCreateFile(activeSession, false, rootDirectorySubDir["uuid"].(string), cpUuid, "subSubDir")
		subDirSubFile := testCreateFile(activeSession, true, subDir["uuid"].(string), cpUuid, "subSubDirFile.mjs")
		testUpdateFileContent(activeSession, cpUuid, subDirSubFile["uuid"].(string), `
function subDirDirFileExecute() {
    console.log('subSubDirFile');
}

export { subDirDirFileExecute };
`)
		sessionUuid := testCreateLinkedSession(activeSession, pageUuid, cb["uuid"].(string))

		bm := map[string]interface{}{
			"uuid": sessionUuid,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/linked-code-block", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executeLinkedProjectHandler)

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
		gomega.Expect(result.Result).Should(gomega.Equal("rootDirectoryFile2\nsubDirFile\nrootDirectoryFile1\n"))
	})

	GinkgoIt("Should run a linked code block execution as a session in a C# (Mono) environment", func() {
		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()
		link := testCreateBlog(activeSession)
		pageUuid := link["page"].(map[string]interface{})["uuid"].(string)
		blogUuid := link["blog"].(map[string]interface{})["uuid"].(string)

		cp := testCreateCodeProject(activeSession, uuid.New().String(), runner.CSharpMono)
		cpUuid := cp["uuid"].(string)
		cb := testCreateCodeBlock(pageUuid, activeSession)
		testUpdateCodeBlock(activeSession, pageUuid, cb["uuid"].(string), `
public class HelloWorld
{
    public static void Main(string[] args)
    {
        NewClass v = new NewClass();
        v.Fn();
    }
}
`)
		testLinkCodeProject(activeSession, cp["uuid"].(string), pageUuid, cb["uuid"].(string), blogUuid)

		var rootDirectory map[string]interface{}
		s, err := json.Marshal(cp["rootDirectory"])
		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(json.Unmarshal(s, &rootDirectory)).Should(gomega.BeNil())

		rootDirectoryFile1 := testCreateFile(activeSession, true, rootDirectory["uuid"].(string), cpUuid, "rootDirectoryFile1.cs")
		testUpdateFileContent(activeSession, cpUuid, rootDirectoryFile1["uuid"].(string), fmt.Sprintf(`
using System;

public class NewClass {
    public void Fn() {
                Console.WriteLine ("Hello World");
    }
}
`))

		sessionUuid := testCreateLinkedSession(activeSession, pageUuid, cb["uuid"].(string))

		bm := map[string]interface{}{
			"uuid": sessionUuid,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/linked-code-block", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executeLinkedProjectHandler)

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
		gomega.Expect(result.Result).Should(gomega.Equal("Hello World\n"))
	})
})
