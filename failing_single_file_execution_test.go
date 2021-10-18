package main

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/httpUtil"
	"therebelsource/emulator/staticTypes"
)

var _ = GinkgoDescribe("Single file execution tests", func() {
	GinkgoIt("Should fail if code block and page does not exist", func() {
		testPrepare()
		defer testCleanup()

		bm := map[string]interface{}{
			"pageUuid":        uuid.New().String(),
			"blockUuid":       uuid.New().String(),
			"state": "single_file",
			"type": "blog",
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executeSingleCodeBlockHandler)

		handler.ServeHTTP(rr, req)

		b := rr.Body.Bytes()

		var apiResponse httpUtil.ApiResponse
		err = json.Unmarshal(b, &apiResponse)

		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusBadRequest))
		gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
		gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_ERROR))
		gomega.Expect(apiResponse.Message).To(gomega.Equal("An appErrors occurred with MasterCode: 1; ApplicationCode: 1; Message: Request data is invalid"))
		gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(appErrors.ApplicationError))
		gomega.Expect(apiResponse.Code).To(gomega.Equal(appErrors.ApplicationRuntimeError))
		gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusBadRequest))
		gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

		gomega.Expect(apiResponse.Data).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Data.(map[string]interface{})["blockExists"]).To(gomega.Not(gomega.BeEmpty()))
	})

	GinkgoIt("Should fail if page does not exist", func() {
		testPrepare()
		defer testCleanup()

		pg := testCreateEmptyPage()
		cb := testCreateCodeBlock(pg["uuid"].(string))

		bm := map[string]interface{}{
			"pageUuid":        uuid.New().String(),
			"blockUuid":       cb["uuid"].(string),
			"state": "single_file",
			"type": "blog",
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executeSingleCodeBlockHandler)

		handler.ServeHTTP(rr, req)

		b := rr.Body.Bytes()

		var apiResponse httpUtil.ApiResponse
		err = json.Unmarshal(b, &apiResponse)

		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusBadRequest))
		gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
		gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_ERROR))
		gomega.Expect(apiResponse.Message).To(gomega.Equal("An appErrors occurred with MasterCode: 1; ApplicationCode: 1; Message: Request data is invalid"))
		gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(appErrors.ApplicationError))
		gomega.Expect(apiResponse.Code).To(gomega.Equal(appErrors.ApplicationRuntimeError))
		gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusBadRequest))
		gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

		gomega.Expect(apiResponse.Data).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Data.(map[string]interface{})["blockExists"]).To(gomega.Not(gomega.BeEmpty()))
	})

	GinkgoIt("Should fail if code block does not exist", func() {
		testPrepare()
		defer testCleanup()

		pg := testCreateEmptyPage()

		bm := map[string]interface{}{
			"pageUuid":        pg["uuid"].(string),
			"blockUuid":       uuid.New().String(),
			"state": "single_file",
			"type": "blog",
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executeSingleCodeBlockHandler)

		handler.ServeHTTP(rr, req)

		b := rr.Body.Bytes()

		var apiResponse httpUtil.ApiResponse
		err = json.Unmarshal(b, &apiResponse)

		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusBadRequest))
		gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
		gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_ERROR))
		gomega.Expect(apiResponse.Message).To(gomega.Equal("An appErrors occurred with MasterCode: 1; ApplicationCode: 1; Message: Request data is invalid"))
		gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(appErrors.ApplicationError))
		gomega.Expect(apiResponse.Code).To(gomega.Equal(appErrors.ApplicationRuntimeError))
		gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusBadRequest))
		gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

		gomega.Expect(apiResponse.Data).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Data.(map[string]interface{})["blockExists"]).To(gomega.Not(gomega.BeEmpty()))
	})

	GinkgoIt("Should fail if the state is invalid", func() {
		testPrepare()
		defer testCleanup()

		pg := testCreateEmptyPage()
		cb := testCreateCodeBlock(pg["uuid"].(string))

		bm := map[string]interface{}{
			"pageUuid":        pg["uuid"].(string),
			"blockUuid":       cb["uuid"].(string),
			"state": "invalid_state",
			"type": "blog",
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executeSingleCodeBlockHandler)

		handler.ServeHTTP(rr, req)

		b := rr.Body.Bytes()

		var apiResponse httpUtil.ApiResponse
		err = json.Unmarshal(b, &apiResponse)

		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusBadRequest))
		gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
		gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_ERROR))
		gomega.Expect(apiResponse.Message).To(gomega.Equal("An appErrors occurred with MasterCode: 1; ApplicationCode: 1; Message: Request data is invalid"))
		gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(appErrors.ApplicationError))
		gomega.Expect(apiResponse.Code).To(gomega.Equal(appErrors.ApplicationRuntimeError))
		gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusBadRequest))
		gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

		gomega.Expect(apiResponse.Data).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Data.(map[string]interface{})["stateValid"]).To(gomega.Not(gomega.BeEmpty()))
	})

	GinkgoIt("Should fail if the state is invalid", func() {
		testPrepare()
		defer testCleanup()

		pg := testCreateEmptyPage()
		cb := testCreateCodeBlock(pg["uuid"].(string))

		bm := map[string]interface{}{
			"pageUuid":        pg["uuid"].(string),
			"blockUuid":       cb["uuid"].(string),
			"state": "single_file",
			"type": "invalid",
		}

		body, err := json.Marshal(bm)

		gomega.Expect(err).To(gomega.BeNil())

		req, err := http.NewRequest("POST", "/api/environment-emulator/execute/single-file", bytes.NewReader(body))

		if err != nil {
			ginkgo.Fail(err.Error())

			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(executeSingleCodeBlockHandler)

		handler.ServeHTTP(rr, req)

		b := rr.Body.Bytes()

		var apiResponse httpUtil.ApiResponse
		err = json.Unmarshal(b, &apiResponse)

		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(rr.Code).To(gomega.Equal(http.StatusBadRequest))
		gomega.Expect(rr.Body).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Method).To(gomega.Equal("POST"))
		gomega.Expect(apiResponse.Type).To(gomega.Equal(staticTypes.RESPONSE_ERROR))
		gomega.Expect(apiResponse.Message).To(gomega.Equal("An appErrors occurred with MasterCode: 1; ApplicationCode: 1; Message: Request data is invalid"))
		gomega.Expect(apiResponse.MasterCode).To(gomega.Equal(appErrors.ApplicationError))
		gomega.Expect(apiResponse.Code).To(gomega.Equal(appErrors.ApplicationRuntimeError))
		gomega.Expect(apiResponse.Status).To(gomega.Equal(http.StatusBadRequest))
		gomega.Expect(apiResponse.Pagination).To(gomega.BeNil())

		gomega.Expect(apiResponse.Data).To(gomega.Not(gomega.BeNil()))

		gomega.Expect(apiResponse.Data.(map[string]interface{})["typeValid"]).To(gomega.Not(gomega.BeEmpty()))
	})
})

