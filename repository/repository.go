package repository

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/httpClient"
)

type Repository struct{}

func InitRepository() Repository {
	return Repository{}
}

func CreateApiUrl() string {
	return fmt.Sprintf("%s://%s/%s/%s", os.Getenv("API_PROTOCOL"), os.Getenv("API_HOST"), os.Getenv("API_ENDPOINT"), os.Getenv("API_VERSION"))
}

func (r Repository) GetCodeBlock(sessionUuid string) (*CodeBlock, *appErrors.Error) {
	url := fmt.Sprintf("%s/page/temp-session/single-file", CreateApiUrl())

	client, err := httpClient.NewHttpClient(&tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Cannot create client to get environment data")
	}

	bm := map[string]interface{}{
		"uuid": sessionUuid,
	}

	body, err := json.Marshal(bm)

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Cannot create client to get environment data")
	}

	response, clientError := client.MakeJsonRequest(&httpClient.JsonRequest{
		Url:    url,
		Method: "POST",
		Body:   body,
	})

	if clientError != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, clientError.GetMessage())
	}

	if response.Status != 200 {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Request did not succeed with status: %d", response.Status))
	}

	var apiResponse map[string]interface{}
	err = json.Unmarshal(response.Body, &apiResponse)

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, err.Error())
	}

	d := apiResponse["data"]

	b, err := json.Marshal(d)

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Cannot get code block: %s", err.Error()))
	}

	var codeBlock *CodeBlock
	if err := json.Unmarshal(b, &codeBlock); err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Cannot get code block: %s", err.Error()))
	}

	return codeBlock, nil
}

func (r Repository) GetProjectSessionData(sessionUuid string) (*SessionCodeProjectData, *appErrors.Error) {
	url := fmt.Sprintf("%s/code-project/temp-session/project", CreateApiUrl())

	client, err := httpClient.NewHttpClient(&tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Cannot create client to get environment data")
	}

	bm := map[string]interface{}{
		"uuid": sessionUuid,
	}

	body, err := json.Marshal(bm)

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Cannot create client to get environment data")
	}

	response, clientError := client.MakeJsonRequest(&httpClient.JsonRequest{
		Url:    url,
		Method: "POST",
		Body:   body,
	})

	if clientError != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, clientError.GetMessage())
	}

	if response.Status != 200 {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Request did not succeed with status: %d", response.Status))
	}

	var apiResponse map[string]interface{}
	err = json.Unmarshal(response.Body, &apiResponse)

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, err.Error())
	}

	d := apiResponse["data"]

	b, err := json.Marshal(d)

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Cannot get session data: %s", err.Error()))
	}

	var sessionData *SessionCodeProjectData
	if err := json.Unmarshal(b, &sessionData); err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Cannot get session data: %s", err.Error()))
	}

	return sessionData, nil
}

func (r Repository) GetSingleFileSessionData(sessionUuid string) (*SingleFileSessionData, *appErrors.Error) {
	url := fmt.Sprintf("%s/code-project/temp-session/single-file", CreateApiUrl())

	client, err := httpClient.NewHttpClient(&tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Cannot create client to get environment data")
	}

	bm := map[string]interface{}{
		"uuid": sessionUuid,
	}

	body, err := json.Marshal(bm)

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Cannot create client to get environment data")
	}

	response, clientError := client.MakeJsonRequest(&httpClient.JsonRequest{
		Url:    url,
		Method: "POST",
		Body:   body,
	})

	if clientError != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, clientError.GetMessage())
	}

	if response.Status != 200 {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Request did not succeed with status: %d", response.Status))
	}

	var apiResponse map[string]interface{}
	err = json.Unmarshal(response.Body, &apiResponse)

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, err.Error())
	}

	d := apiResponse["data"]

	b, err := json.Marshal(d)

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Cannot get session data: %s", err.Error()))
	}

	var sessionData *SingleFileSessionData
	if err := json.Unmarshal(b, &sessionData); err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Cannot get session data: %s", err.Error()))
	}

	return sessionData, nil
}

func (r Repository) GetLinkedSessionData(sessionUuid string) (*LinkedSessionData, *appErrors.Error) {
	url := fmt.Sprintf("%s/code-project/temp-session/linked-code-block", CreateApiUrl())

	client, err := httpClient.NewHttpClient(&tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Cannot create client to get environment data")
	}

	bm := map[string]interface{}{
		"uuid": sessionUuid,
	}

	body, err := json.Marshal(bm)

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Cannot create client to get environment data")
	}

	response, clientError := client.MakeJsonRequest(&httpClient.JsonRequest{
		Url:    url,
		Method: "POST",
		Body:   body,
	})

	if clientError != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, clientError.GetMessage())
	}

	if response.Status != 200 {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Request did not succeed with status: %d", response.Status))
	}

	var apiResponse map[string]interface{}
	err = json.Unmarshal(response.Body, &apiResponse)

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, err.Error())
	}

	d := apiResponse["data"]

	b, err := json.Marshal(d)

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Cannot get session data: %s", err.Error()))
	}

	var sessionData *LinkedSessionData
	if err := json.Unmarshal(b, &sessionData); err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Cannot get session data: %s", err.Error()))
	}

	return sessionData, nil
}

func (r Repository) ValidateTemporarySession(sessionUuid string) (ValidatedTemporarySession, *appErrors.Error) {
	url := fmt.Sprintf("%s/auth/validate/emulator-temporary-session", CreateApiUrl())

	client, err := httpClient.NewHttpClient(&tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		return ValidatedTemporarySession{}, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Cannot create client to get environment data")
	}

	bm := map[string]interface{}{
		"uuid": sessionUuid,
	}

	body, err := json.Marshal(bm)

	if err != nil {
		return ValidatedTemporarySession{}, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Cannot create client to get environment data")
	}

	response, clientError := client.MakeJsonRequest(&httpClient.JsonRequest{
		Url:    url,
		Method: "POST",
		Body:   body,
	})

	if clientError != nil {
		return ValidatedTemporarySession{}, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, clientError.GetMessage())
	}
	
	if response.Status != 200 {
		return ValidatedTemporarySession{}, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Request did not succeed with status: %d", response.Status))
	}

	var apiResponse map[string]interface{}
	err = json.Unmarshal(response.Body, &apiResponse)

	if err != nil {
		return ValidatedTemporarySession{}, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, err.Error())
	}

	d := apiResponse["data"]

	b, err := json.Marshal(d)

	if err != nil {
		return ValidatedTemporarySession{}, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Cannot get code block: %s", err.Error()))
	}

	var codeBlock ValidatedTemporarySession
	if err := json.Unmarshal(b, &codeBlock); err != nil {
		return ValidatedTemporarySession{}, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Cannot get code block: %s", err.Error()))
	}

	return codeBlock, nil
}

func (r Repository) InvalidateTemporarySession(sessionUuid string) *appErrors.Error {
	url := fmt.Sprintf("%s/auth/invalidate/emulator-temporary-session", CreateApiUrl())

	client, err := httpClient.NewHttpClient(&tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		return appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Cannot create client to get environment data")
	}

	bm := map[string]interface{}{
		"uuid": sessionUuid,
	}

	body, err := json.Marshal(bm)

	if err != nil {
		return appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Cannot create client to get environment data")
	}

	response, clientError := client.MakeJsonRequest(&httpClient.JsonRequest{
		Url:    url,
		Method: "POST",
		Body:   body,
	})

	if clientError != nil {
		return appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, clientError.GetMessage())
	}

	if response.Status != 200 {
		return appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Request did not succeed with status: %d", response.Status))
	}

	var apiResponse map[string]interface{}
	err = json.Unmarshal(response.Body, &apiResponse)

	if err != nil {
		return appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, err.Error())
	}

	d := apiResponse["data"]

	b, err := json.Marshal(d)

	if err != nil {
		return appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Cannot get code block: %s", err.Error()))
	}

	var codeBlock ValidatedTemporarySession
	if err := json.Unmarshal(b, &codeBlock); err != nil {
		return appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Cannot get code block: %s", err.Error()))
	}

	return nil
}

func (r Repository) GetCodeProject(codeProjectUuid string) (*CodeProject, *appErrors.Error) {
	url := fmt.Sprintf("%s/code-project/%s", CreateApiUrl(), codeProjectUuid)

	client, err := httpClient.NewHttpClient(&tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Cannot create client to get environment data")
	}

	response, clientError := client.MakeJsonRequest(&httpClient.JsonRequest{
		Url:    url,
		Method: "GET",
		Body:   nil,
	})

	if clientError != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, clientError.GetMessage())
	}

	if response.Status != 200 {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Request did not succeed with status: %d", response.Status))
	}

	var apiResponse map[string]interface{}
	err = json.Unmarshal(response.Body, &apiResponse)

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, err.Error())
	}

	d := apiResponse["data"]

	b, err := json.Marshal(d)

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Cannot get code project: %s", err.Error()))
	}

	var codeProject *CodeProject
	if err := json.Unmarshal(b, &codeProject); err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Cannot get code project: %s", err.Error()))
	}

	return codeProject, nil
}

func (r Repository) GetAllFileContent(codeProjectUuid string) ([]*FileContent, *appErrors.Error) {
	url := fmt.Sprintf("%s/code-project/file/content/%s", CreateApiUrl(), codeProjectUuid)

	client, err := httpClient.NewHttpClient(&tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Cannot create client to get environment data")
	}

	response, clientError := client.MakeJsonRequest(&httpClient.JsonRequest{
		Url:    url,
		Method: "GET",
		Body:   nil,
	})

	if clientError != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, clientError.GetMessage())
	}

	if response.Status != 200 {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Request did not succeed with status: %d", response.Status))
	}

	var apiResponse map[string]interface{}
	err = json.Unmarshal(response.Body, &apiResponse)

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, err.Error())
	}

	d := apiResponse["data"]

	b, err := json.Marshal(d)

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Cannot get code projects file system contents: %s", err.Error()))
	}

	var contents []*FileContent
	if err := json.Unmarshal(b, &contents); err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, fmt.Sprintf("Cannot get code projects file system contents: %s", err.Error()))
	}

	return contents, nil
}
