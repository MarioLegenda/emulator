package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"os/exec"
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

func testCleanup() {
	cmd := exec.Command("bash", "-c", "/usr/bin/rm -rf /var/www/projects")
	_, err := cmd.Output()

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot do cleanup: %s", err.Error()))
		return
	}

	runner.StopContainerBalancer()
}

func testCreateEmptyPage() map[string]interface{} {
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

func testCreateCodeBlock(pageUuid string) map[string]interface{} {
	url := fmt.Sprintf("%s/page/code-block", repository.CreateApiUrl())

	client, err := httpClient.NewHttpClient(&tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		appErrors.TerminateWithMessage(fmt.Sprintf("Cannot create page: %s", err.Error()))
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

func testAddEmulatorToCodeBlock(pageUuid string, blockUuid string, code string, lang runner.Language) map[string]interface{} {
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