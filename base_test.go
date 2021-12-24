package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"io"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/httpClient"
	"therebelsource/emulator/repository"
	"therebelsource/emulator/runner"
	"therebelsource/emulator/singleFileExecution"
	"therebelsource/emulator/staticTypes"
)

var GomegaRegisterFailHandler = gomega.RegisterFailHandler
var GinkgoFail = ginkgo.Fail
var GinkgoRunSpecs = ginkgo.RunSpecs
var GinkgoBeforeSuite = ginkgo.BeforeSuite
var GinkgoAfterSuite = ginkgo.AfterSuite
var GinkgoDescribe = ginkgo.Describe
var GinkgoIt = ginkgo.It

var cancelFn context.CancelFunc

func TestApi(t *testing.T) {
	GomegaRegisterFailHandler(GinkgoFail)
	GinkgoRunSpecs(t, "API Suite")
}

func testPrepare() {
	LoadEnv(staticTypes.APP_DEV_ENV)
	InitRequiredDirectories(false)

	singleFileExecution.InitService()

	runner.StartContainerBalancer()
}

func testCreateSecureRequest(rr *httptest.ResponseRecorder, sessionUuid string, method string, path string, body io.Reader) *http.Request {
	http.SetCookie(rr, &http.Cookie{Name: "session", Value: sessionUuid, Path: "/", MaxAge: 3600, Secure: true, HttpOnly: true})
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		ginkgo.Fail(err.Error())
	}
	req.AddCookie(&http.Cookie{
		Name:       "session",
		Value:      sessionUuid,
		Path:       "/",
		MaxAge:     3600,
		Secure:     true,
		HttpOnly:   true,
	})

	return req
}

func testCleanup() {
	cmd := exec.Command("bash", "-c", "/usr/bin/rm -rf /var/www/projects")
	_, err := cmd.Output()

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot do cleanup: %s", err.Error()))
		return
	}

	runner.StopContainerBalancer()
}

func testCreateEmptyPage(activeSession repository.ActiveSession) map[string]interface{} {
	url := fmt.Sprintf("%s/page", repository.CreateApiUrl())

	client, err := httpClient.NewHttpClient(&tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create page: %s", err.Error()))
	}

	response, clientError := client.MakeJsonRequest(&httpClient.JsonRequest{
		Url:    url,
		Method: "PUT",
		Body:   nil,
		Session: activeSession.Session.Uuid,
	})

	if clientError != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create page: %s", err.Error()))
	}

	if response.Status != 201 {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create page: Response status is %d", response.Status))
	}

	var apiResponse map[string]interface{}
	err = json.Unmarshal(response.Body, &apiResponse)

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create page: %s", err.Error()))
	}

	d := apiResponse["data"]

	b, err := json.Marshal(d)

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create page: %s", err.Error()))
	}

	var data map[string]interface{}
	gomega.Expect(json.Unmarshal(b, &data)).Should(gomega.BeNil())

	return data
}

func testCreateTemporarySession(activeSession repository.ActiveSession, pageUuid string, blockUuid string, t string) repository.TemporarySession {
	url := fmt.Sprintf("%s/auth/single-file-emulator-temporary-session", repository.CreateApiUrl())

	client, err := httpClient.NewHttpClient(&tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create temporary session: %s", err.Error()))
	}

	bm := map[string]interface{}{
		"pageUuid": pageUuid,
		"blockUuid": blockUuid,
		"type": t,
		"accountUuid": activeSession.Account.Uuid,
	}

	body, err := json.Marshal(bm)

	gomega.Expect(err).To(gomega.BeNil())

	response, clientError := client.MakeJsonRequest(&httpClient.JsonRequest{
		Url:    url,
		Method: "POST",
		Body:   body,
	})

	if clientError != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create temporary session: %s", err.Error()))
	}

	if response.Status != 201 {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create temporary session: Response status is %d", response.Status))
	}

	var apiResponse map[string]interface{}
	err = json.Unmarshal(response.Body, &apiResponse)

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create temporary session: %s", err.Error()))
	}

	d := apiResponse["data"]

	b, err := json.Marshal(d)

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create temporary session: %s", err.Error()))
	}

	var data repository.TemporarySession
	gomega.Expect(json.Unmarshal(b, &data)).Should(gomega.BeNil())

	return data
}

func testCreateProjectTemporarySession(activeSession repository.ActiveSession, codeProjectUuid string) repository.TemporarySession {
	url := fmt.Sprintf("%s/auth/project-emulator-temporary-session", repository.CreateApiUrl())

	client, err := httpClient.NewHttpClient(&tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create temporary session: %s", err.Error()))
	}

	bm := map[string]interface{}{
		"codeProjectUuid": codeProjectUuid,
		"type": "project",
	}

	body, err := json.Marshal(bm)

	gomega.Expect(err).To(gomega.BeNil())

	response, clientError := client.MakeJsonRequest(&httpClient.JsonRequest{
		Url:    url,
		Method: "POST",
		Body:   body,
	})

	if clientError != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create temporary session: %s", err.Error()))
	}

	if response.Status != 201 {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create temporary session: Response status is %d", response.Status))
	}

	var apiResponse map[string]interface{}
	err = json.Unmarshal(response.Body, &apiResponse)

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create temporary session: %s", err.Error()))
	}

	d := apiResponse["data"]

	b, err := json.Marshal(d)

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create temporary session: %s", err.Error()))
	}

	var data repository.TemporarySession
	gomega.Expect(json.Unmarshal(b, &data)).Should(gomega.BeNil())

	return data
}

func testGetEmail() string {
	return fmt.Sprintf("%s@gmail.com", strings.Split(uuid.New().String(), "-")[0])
}

func testCreateAccount() repository.ActiveSession {
	url := fmt.Sprintf("%s/user", repository.CreateApiUrl())

	client, err := httpClient.NewHttpClient(&tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create account: %s", err.Error()))
	}

	bm := map[string]interface{}{
		"name": "name",
		"lastName": "Last name",
		"email": testGetEmail(),
		"password": "mypassword",
	}

	body, err := json.Marshal(bm)

	gomega.Expect(err).To(gomega.BeNil())

	response, clientError := client.MakeJsonRequest(&httpClient.JsonRequest{
		Url:    url,
		Method: "PUT",
		Body:   body,
	})

	if clientError != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create account: %s", err.Error()))
	}

	var apiResponse map[string]interface{}
	err = json.Unmarshal(response.Body, &apiResponse)

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create account: %s", err.Error()))
	}

	if response.Status != 200 {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create account: Response status is %d", response.Status))
	}

	d := apiResponse["data"]

	b, err := json.Marshal(d)

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create page: %s", err.Error()))
	}

	var data repository.ActiveSession
	gomega.Expect(json.Unmarshal(b, &data)).Should(gomega.BeNil())

	return data
}

func testCreateCodeBlock(pageUuid string, activeSession repository.ActiveSession) map[string]interface{} {
	url := fmt.Sprintf("%s/page/code-block", repository.CreateApiUrl())

	client, err := httpClient.NewHttpClient(&tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create code block: %s", err.Error()))
	}

	bm := map[string]interface{}{
		"pageUuid": pageUuid,
	}

	body, err := json.Marshal(bm)

	gomega.Expect(err).To(gomega.BeNil())

	response, clientError := client.MakeJsonRequest(&httpClient.JsonRequest{
		Url:    url,
		Method: "PUT",
		Body:   body,
		Session: activeSession.Session.Uuid,
	})

	if clientError != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create code block: %s", err.Error()))
	}

	if response.Status != 201 {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create code block: Response status is %d", response.Status))
	}

	var apiResponse map[string]interface{}
	err = json.Unmarshal(response.Body, &apiResponse)

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create code block: %s", err.Error()))
	}

	d := apiResponse["data"]

	b, err := json.Marshal(d)

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create code block: %s", err.Error()))
	}

	var data map[string]interface{}
	gomega.Expect(json.Unmarshal(b, &data)).Should(gomega.BeNil())

	return data
}

func testAddEmulatorToCodeBlock(pageUuid string, blockUuid string, code string, lang runner.Language, activeSession repository.ActiveSession) map[string]interface{} {
	url := fmt.Sprintf("%s/page/code-block", repository.CreateApiUrl())

	client, err := httpClient.NewHttpClient(&tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create page: %s", err.Error()))
	}

	bm := map[string]interface{}{
		"pageUuid":  pageUuid,
		"blockUuid": blockUuid,
		"text":      code,
		"emulator": map[string]interface{}{
			"name": lang.Name,
			"text": lang.Text,
			"tag": lang.Tag,
			"inDevelopment": false,
			"inMaintenance": false,
			"language": lang.Language,
			"extension": lang.Extension,
			"output": "",
			"defaultTimeout": 0,
			"packageTimeout": 0,
		},
		"update": []string{"emulator", "text"},
	}

	body, err := json.Marshal(bm)

	gomega.Expect(err).To(gomega.BeNil())

	response, clientError := client.MakeJsonRequest(&httpClient.JsonRequest{
		Url:    url,
		Method: "POST",
		Body:   body,
		Session: activeSession.Session.Uuid,
	})

	if clientError != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot update page: %s", err.Error()))
	}

	if response.Status != 200 {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot update page: Response status is %d", response.Status))
	}

	var apiResponse map[string]interface{}
	err = json.Unmarshal(response.Body, &apiResponse)

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot update page: %s", err.Error()))
	}

	d := apiResponse["data"]

	b, err := json.Marshal(d)

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot update page: %s", err.Error()))
	}

	var data map[string]interface{}
	gomega.Expect(json.Unmarshal(b, &data)).Should(gomega.BeNil())

	return data
}

func testCreateCodeProject(activeSession repository.ActiveSession, lang runner.Language) map[string]interface{} {
	url := fmt.Sprintf("%s/code-project", repository.CreateApiUrl())

	client, err := httpClient.NewHttpClient(&tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create code project: %s", err.Error()))
	}

	bm := map[string]interface{}{
		"name": uuid.New().String(),
		"description": "description",
		"environment": lang,
	}

	body, err := json.Marshal(bm)

	gomega.Expect(err).To(gomega.BeNil())

	response, clientError := client.MakeJsonRequest(&httpClient.JsonRequest{
		Url:    url,
		Method: "PUT",
		Body:   body,
		Session: activeSession.Session.Uuid,
	})

	if clientError != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create code project: %s", err.Error()))
	}

	if response.Status != 201 {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create code project: Response status is %d", response.Status))
	}

	var apiResponse map[string]interface{}
	err = json.Unmarshal(response.Body, &apiResponse)

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create code project: %s", err.Error()))
	}

	d := apiResponse["data"]

	b, err := json.Marshal(d)

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create code project: %s", err.Error()))
	}

	var data map[string]interface{}
	gomega.Expect(json.Unmarshal(b, &data)).Should(gomega.BeNil())

	return data
}

func testCreateFile(activeSession repository.ActiveSession, isFile bool, parent string, cpUuid string, name string) map[string]interface{} {
	url := fmt.Sprintf("%s/code-project/file", repository.CreateApiUrl())

	client, err := httpClient.NewHttpClient(&tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create file: %s", err.Error()))
	}

	bm := map[string]interface{}{
		"isFile": isFile,
		"parent": parent,
		"codeProjectUuid": cpUuid,
		"name": name,
	}

	body, err := json.Marshal(bm)

	gomega.Expect(err).To(gomega.BeNil())

	response, clientError := client.MakeJsonRequest(&httpClient.JsonRequest{
		Url:    url,
		Method: "PUT",
		Body:   body,
		Session: activeSession.Session.Uuid,
	})

	if clientError != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create file: %s", err.Error()))
	}

	if response.Status != 201 {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create file: Response status is %d", response.Status))
	}

	var apiResponse map[string]interface{}
	err = json.Unmarshal(response.Body, &apiResponse)

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create file: %s", err.Error()))
	}

	d := apiResponse["data"]

	b, err := json.Marshal(d)

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create file: %s", err.Error()))
	}

	var data map[string]interface{}
	gomega.Expect(json.Unmarshal(b, &data)).Should(gomega.BeNil())

	return data
}

func testUpdateFileContent(activeSession repository.ActiveSession, cpUuid string, fileUuid string, content string) map[string]interface{} {
	url := fmt.Sprintf("%s/code-project/file/update-content", repository.CreateApiUrl())

	client, err := httpClient.NewHttpClient(&tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot update file content: %s", err.Error()))
	}

	bm := map[string]interface{}{
		"codeProjectUuid": cpUuid,
		"uuid": fileUuid,
		"content": content,
	}

	body, err := json.Marshal(bm)

	gomega.Expect(err).To(gomega.BeNil())

	response, clientError := client.MakeJsonRequest(&httpClient.JsonRequest{
		Url:    url,
		Method: "POST",
		Body:   body,
		Session: activeSession.Session.Uuid,
	})

	if clientError != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot update file content: %s", err.Error()))
	}

	if response.Status != 200 {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create file content: Response status is %d", response.Status))
	}

	var apiResponse map[string]interface{}
	err = json.Unmarshal(response.Body, &apiResponse)

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create file content: %s", err.Error()))
	}

	d := apiResponse["data"]

	b, err := json.Marshal(d)

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create file content: %s", err.Error()))
	}

	var data map[string]interface{}
	gomega.Expect(json.Unmarshal(b, &data)).Should(gomega.BeNil())

	return data
}

