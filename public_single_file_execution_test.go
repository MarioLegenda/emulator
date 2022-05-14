package main

import (
	"bytes"
	"encoding/json"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"therebelsource/emulator/httpUtil"
	"therebelsource/emulator/runner"
	"therebelsource/emulator/staticTypes"
)

var _ = GinkgoDescribe("Single file public execution tests", func() {
	GinkgoIt("Should public execute a single file in a node LTS environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `console.log("mile")`, runner.NodeLts, activeSession)
		sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

		bm := map[string]interface{}{
			"uuid": sessionUuid,
			"text": `
console.log('mile');
`,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/public/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executePublicSingleFileRunResult)

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
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))
	})

	GinkgoIt("Should public execute a single file in a node LTS ESM environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `console.log("mile")`, runner.NodeEsm, activeSession)
		sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

		bm := map[string]interface{}{
			"uuid": sessionUuid,
			"text": `
console.log('mile');
`,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/public/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executePublicSingleFileRunResult)

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
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))
	})

	GinkgoIt("Should public execute a single file in a C# (Mono) environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), ``, runner.CSharpMono, activeSession)
		sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

		bm := map[string]interface{}{
			"uuid": sessionUuid,
			"text": `
using System;

public class HelloWorld
{
    public static void Main(string[] args)
    {
        Console.WriteLine ("Hello World");
    }
}
`,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/public/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executePublicSingleFileRunResult)

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

	GinkgoIt("Should execute a single file in a node 14.x environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `console.log("mile")`, runner.NodeLts, activeSession)
		sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

		bm := map[string]interface{}{
			"uuid": sessionUuid,
			"text": `
console.log('mile');
`,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/public/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executePublicSingleFileRunResult)

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
		gomega.Expect(result.Result).Should(gomega.Equal("mile\n"))
	})

	GinkgoIt("Should execute a single file in a PHP environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `
<?php

echo "mile";
`, runner.Php74, activeSession)
		sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

		bm := map[string]interface{}{
			"uuid": sessionUuid,
			"text": `
<?php

echo "mile";
`,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/public/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executePublicSingleFileRunResult)

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
		gomega.Expect(result.Result).Should(gomega.Equal("\nmile"))
	})

	GinkgoIt("Should execute a single file in a Ruby environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `puts "mile"`, runner.Ruby, activeSession)
		sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

		bm := map[string]interface{}{
			"uuid": sessionUuid,
			"text": `
print "mile"
`,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/public/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executePublicSingleFileRunResult)

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
		gomega.Expect(result.Result).Should(gomega.Equal("mile"))
	})

	GinkgoIt("Should execute a single file in a Go environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `
package main

import "fmt"

func main() {
	fmt.Println("mile")
}
`, runner.GoLang, activeSession)
		sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

		bm := map[string]interface{}{
			"uuid": sessionUuid,
			"text": `
package main

import "fmt"

func main() {
	fmt.Println("something else")
}
`,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/public/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executePublicSingleFileRunResult)

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
		gomega.Expect(result.Result).Should(gomega.Equal("something else\n"))
	})

	GinkgoIt("Should execute a single file in a Python2 environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `print("mile")`, runner.Python2, activeSession)
		sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

		bm := map[string]interface{}{
			"uuid": sessionUuid,
			"text": `
print("This line will be printed.")
`,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/public/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executePublicSingleFileRunResult)

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
		gomega.Expect(result.Result).Should(gomega.Equal("This line will be printed.\n"))
	})

	GinkgoIt("Should execute a single file in a Python3 environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `print("mile")`, runner.Python2, activeSession)
		sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

		bm := map[string]interface{}{
			"uuid": sessionUuid,
			"text": `
print("This line will be printed.")
`,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/public/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executePublicSingleFileRunResult)

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
		gomega.Expect(result.Result).Should(gomega.Equal("This line will be printed.\n"))
	})

	GinkgoIt("Should execute a single file in a Haskell environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `main = putStrLn "mile"`, runner.Haskell, activeSession)
		sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

		bm := map[string]interface{}{
			"uuid": sessionUuid,
			"text": `
main = putStrLn "mile"
`,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/public/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executePublicSingleFileRunResult)

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
		gomega.Expect(result.Result).Should(gomega.Equal("\nmile\n"))
	})

	GinkgoIt("Should execute a single file in a C environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `
#include <stdio.h>
int main() {
   printf("mile");
   return 0;
}
`, runner.CLang, activeSession)
		sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

		bm := map[string]interface{}{
			"uuid": sessionUuid,
			"text": `
#include <stdio.h>
int main() {
   printf("shit");
   return 0;
}
`,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/public/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executePublicSingleFileRunResult)

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
		gomega.Expect(result.Result).Should(gomega.Equal("shit"))
	})

	GinkgoIt("Should execute a single file in a C++ environment", func() {
		ginkgo.Skip("")

		testPrepare()
		defer testCleanup()

		activeSession := testCreateAccount()

		pg := testCreateEmptyPage(activeSession)
		cb := testCreateCodeBlock(pg["uuid"].(string), activeSession)
		testAddEmulatorToCodeBlock(pg["uuid"].(string), cb["uuid"].(string), `
#include <iostream>

int main() {
    std::cout << "mile";
    return 0;
}
`, runner.CPlus, activeSession)
		sessionUuid := testCreateTemporarySession(activeSession, pg["uuid"].(string), cb["uuid"].(string), "single_file")

		bm := map[string]interface{}{
			"uuid": sessionUuid,
			"text": `
#include <iostream>

int main() {
    std::cout << "shit";
    return 0;
}
`,
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/public/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executePublicSingleFileRunResult)

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
		gomega.Expect(result.Result).Should(gomega.Equal("shit"))

		if _, err := os.Stat("/var/www/execution/singleFile"); os.IsNotExist(err) {
			ginkgo.Fail("/var/www/execution/singleFile directory does not exist but should")
		}

		files, err := ioutil.ReadDir("/var/www/execution/singleFile")

		gomega.Expect(err).Should(gomega.BeNil())
		gomega.Expect(files).Should(gomega.HaveLen(0))
	})
})
