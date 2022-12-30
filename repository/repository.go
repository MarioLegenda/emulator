package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/logger"
	newClient "therebelsource/emulator/repository/httpClient"
	"time"
)

func CreateApiUrl() string {
	return fmt.Sprintf("%s://%s/%s/%s", os.Getenv("API_PROTOCOL"), os.Getenv("API_HOST"), os.Getenv("API_ENDPOINT"), os.Getenv("API_VERSION"))
}

func createClient() (*http.Client, *appErrors.Error) {
	var timeout, err = strconv.Atoi(os.Getenv("GLOBAL_REQUEST_TIMEOUT"))

	if err != nil {
		logger.Warn(fmt.Sprintf("Failed creating client: %s", err.Error()))
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, err.Error())
	}

	client := newClient.NewClient(newClient.ClientParams{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * time.Duration(timeout),
	})

	return client, nil
}

func createRequest(url string, method string, body []byte, headers map[string]string) (*http.Request, *appErrors.Error) {
	r, err := newClient.NewRequest(newClient.Request{
		Headers: headers,
		Url:     url,
		Method:  method,
		Body:    body,
	})

	if err != nil {
		logger.Warn(fmt.Sprintf("Failed creating request: %s", err.Error()))
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, err.Error())
	}

	return r, nil
}

func createBody(sessionUuid string) ([]byte, *appErrors.Error) {
	bm := map[string]interface{}{
		"uuid": sessionUuid,
	}

	body, err := json.Marshal(bm)

	if err != nil {
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Cannot create client to get environment data")
	}

	return body, nil
}

func GetCodeBlock(authenticatedSession string, sessionUuid string) (*CodeBlock, *appErrors.Error) {
	logger.Info(fmt.Sprintf("GetCodeBlock: Starting execution"))

	body, err := createBody(sessionUuid)
	if err != nil {
		logger.Warn(fmt.Sprintf("GetCodeBlock: Failed marshaling body to get code block: %s", err.Error()))
		return nil, err
	}

	client, clientErr := createClient()
	if clientErr != nil {
		return nil, clientErr
	}

	r, err := createRequest(fmt.Sprintf("%s/session/single-file-data", CreateApiUrl()), "POST", body, nil)
	if err != nil {
		logger.Warn(fmt.Sprintf("GetCodeBlock: Failed creating request body: %s", err.Error()))
		return nil, err
	}

	if authenticatedSession != "" {
		r.AddCookie(&http.Cookie{Name: "session", Value: authenticatedSession, MaxAge: 3600, Path: "/", HttpOnly: true, Secure: true})
	}

	var codeBlock CodeBlock
	if err := newClient.SendWithBackoff(r, client, &codeBlock); err != nil {
		logger.Warn(fmt.Sprintf("GetCodeBlock: Failed sending request: %s", err.Error()))
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, err.Error())
	}

	logger.Info(fmt.Sprintf("GetCodeBlock: Finished execution"))

	return &codeBlock, nil
}

func GetAnonymousCodeBlock(sessionUuid string) (*CodeBlock, *appErrors.Error) {
	logger.Info(fmt.Sprintf("GetAnonymousCodeBlock: Starging execution"))

	body, err := createBody(sessionUuid)
	if err != nil {
		logger.Warn(fmt.Sprintf("GetAnonymousCodeBlock: Failed marshaling body to get code block: %s", err.Error()))
		return nil, err
	}

	client, clientErr := createClient()
	if clientErr != nil {
		return nil, clientErr
	}

	r, err := createRequest(fmt.Sprintf("%s/session/anonymous/single-file-data", CreateApiUrl()), "POST", body, nil)
	if err != nil {
		logger.Warn(fmt.Sprintf("GetAnonymousCodeBlock: Failed creating request body: %s", err.Error()))
		return nil, err
	}

	var codeBlock CodeBlock
	if err := newClient.SendWithBackoff(r, client, &codeBlock); err != nil {
		logger.Warn(fmt.Sprintf("GetAnonymousCodeBlock: Failed sending request: %s", err.Error()))
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, err.Error())
	}

	logger.Info(fmt.Sprintf("GetAnonymousCodeBlock: Finished execution"))

	return &codeBlock, nil
}

func GetAuthenticatedSnippet(authenticatedSession string, sessionUuid string) (*Snippet, *appErrors.Error) {
	logger.Info(fmt.Sprintf("GetAnonymousCodeBlock: Starting execution"))

	body, err := createBody(sessionUuid)
	if err != nil {
		logger.Warn(fmt.Sprintf("GetAuthenticatedSnippet: Failed marshaling body to get code block: %s", err.Error()))
		return nil, err
	}

	client, clientErr := createClient()
	if clientErr != nil {
		return nil, clientErr
	}

	r, err := createRequest(fmt.Sprintf("%s/session/snippet-data", CreateApiUrl()), "POST", body, nil)
	if err != nil {
		logger.Warn(fmt.Sprintf("GetAuthenticatedSnippet: Failed creating request body: %s", err.Error()))
		return nil, err
	}

	if authenticatedSession != "" {
		r.AddCookie(&http.Cookie{Name: "session", Value: authenticatedSession, MaxAge: 3600, Path: "/", HttpOnly: true, Secure: true})
	}

	var snippet Snippet
	if err := newClient.SendWithBackoff(r, client, &snippet); err != nil {
		logger.Warn(fmt.Sprintf("GetAuthenticatedSnippet: Failed sending request: %s", err.Error()))
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, err.Error())
	}

	logger.Info(fmt.Sprintf("GetAnonymousCodeBlock: Finished execution"))

	return &snippet, nil
}

func GetAnonymousSnippet(sessionUuid string) (*Snippet, *appErrors.Error) {
	logger.Info(fmt.Sprintf("GetAnonymousSnippet: Starting execution"))

	body, err := createBody(sessionUuid)
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed marshaling body to get code block: %s", err.Error()))
		return nil, err
	}

	client, clientErr := createClient()
	if clientErr != nil {
		return nil, clientErr
	}

	r, err := createRequest(fmt.Sprintf("%s/session/anonymous/snippet-data", CreateApiUrl()), "POST", body, nil)
	if err != nil {
		logger.Warn(fmt.Sprintf("GetAnonymousSnippet: Failed creating request body: %s", err.Error()))
		return nil, err
	}

	var snippet Snippet
	if err := newClient.SendWithBackoff(r, client, &snippet); err != nil {
		logger.Warn(fmt.Sprintf("GetAnonymousSnippet: Failed sending request: %s", err.Error()))
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, err.Error())
	}

	logger.Info(fmt.Sprintf("GetAnonymousSnippet: Finished execution"))

	return &snippet, nil
}

func GetProjectSessionData(authenticatedSession string, sessionUuid string) (*SessionCodeProjectData, *appErrors.Error) {
	logger.Info(fmt.Sprintf("GetProjectSessionData: Starting execution"))

	body, err := createBody(sessionUuid)
	if err != nil {
		logger.Warn(fmt.Sprintf("GetProjectSessionData: Failed marshaling body to get code block: %s", err.Error()))
		return nil, err
	}

	client, clientErr := createClient()
	if clientErr != nil {
		return nil, clientErr
	}

	r, err := createRequest(fmt.Sprintf("%s/session/project-data", CreateApiUrl()), "POST", body, nil)
	if err != nil {
		logger.Warn(fmt.Sprintf("GetProjectSessionData: Failed creating request body: %s", err.Error()))
		return nil, err
	}

	if authenticatedSession != "" {
		r.AddCookie(&http.Cookie{Name: "session", Value: authenticatedSession, MaxAge: 3600, Path: "/", HttpOnly: true, Secure: true})
	}

	var model SessionCodeProjectData
	if err := newClient.SendWithBackoff(r, client, &model); err != nil {
		logger.Warn(fmt.Sprintf("GetProjectSessionData: Failed sending request: %s", err.Error()))
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, err.Error())
	}

	logger.Info(fmt.Sprintf("GetProjectSessionData: Finished execution"))

	return &model, nil
}

func GetLinkedSessionData(authenticatedSession string, sessionUuid string) (*LinkedSessionData, *appErrors.Error) {
	logger.Info(fmt.Sprintf("GetLinkedSessionData: Starting execution"))

	body, err := createBody(sessionUuid)
	if err != nil {
		logger.Warn(fmt.Sprintf("GetLinkedSessionData: Failed marshaling body to get code block: %s", err.Error()))
		return nil, err
	}

	client, clientErr := createClient()
	if clientErr != nil {
		return nil, clientErr
	}

	r, err := createRequest(fmt.Sprintf("%s/session/linked-data", CreateApiUrl()), "POST", body, nil)
	if err != nil {
		logger.Warn(fmt.Sprintf("GetLinkedSessionData: Failed creating request body: %s", err.Error()))
		return nil, err
	}

	if authenticatedSession != "" {
		r.AddCookie(&http.Cookie{Name: "session", Value: authenticatedSession, MaxAge: 3600, Path: "/", HttpOnly: true, Secure: true})
	}

	var model LinkedSessionData
	if err := newClient.SendWithBackoff(r, client, &model); err != nil {
		logger.Warn(fmt.Sprintf("GetLinkedSessionData: Failed sending request: %s", err.Error()))
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, err.Error())
	}

	logger.Info(fmt.Sprintf("GetLinkedSessionData: Finished execution"))

	return &model, nil
}

func GetAnonymousLinkedSessionData(sessionUuid string) (*LinkedSessionData, *appErrors.Error) {
	logger.Info(fmt.Sprintf("GetAnonymousLinkedSessionData: Starting execution"))

	body, err := createBody(sessionUuid)
	if err != nil {
		logger.Warn(fmt.Sprintf("GetAnonymousLinkedSessionData: Failed marshaling body to get code block: %s", err.Error()))
		return nil, err
	}

	client, clientErr := createClient()
	if clientErr != nil {
		return nil, clientErr
	}

	r, err := createRequest(fmt.Sprintf("%s/session/anonymous/linked-data", CreateApiUrl()), "POST", body, nil)
	if err != nil {
		logger.Warn(fmt.Sprintf("GetAnonymousLinkedSessionData: Failed creating request body: %s", err.Error()))
		return nil, err
	}

	var model LinkedSessionData
	if err := newClient.SendWithBackoff(r, client, &model); err != nil {
		logger.Warn(fmt.Sprintf("GetAnonymousLinkedSessionData: Failed sending request: %s", err.Error()))
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, err.Error())
	}

	logger.Info(fmt.Sprintf("GetAnonymousLinkedSessionData: Finished execution"))

	return &model, nil
}

func ValidateTemporarySession(sessionUuid string) (ValidatedTemporarySession, *appErrors.Error) {
	logger.Info(fmt.Sprintf("ValidateTemporarySession: Starting execution"))

	body, err := createBody(sessionUuid)
	if err != nil {
		logger.Warn(fmt.Sprintf("Failed marshaling body to get code block: %s", err.Error()))
		return ValidatedTemporarySession{}, err
	}

	client, clientErr := createClient()
	if clientErr != nil {
		return ValidatedTemporarySession{}, clientErr
	}

	r, err := createRequest(fmt.Sprintf("%s/session/validate", CreateApiUrl()), "POST", body, nil)
	if err != nil {
		logger.Warn(fmt.Sprintf("ValidateTemporarySession: Failed creating request body: %s", err.Error()))
		return ValidatedTemporarySession{}, err
	}

	var model ValidatedTemporarySession
	if err := newClient.SendWithBackoff(r, client, &model); err != nil {
		logger.Warn(fmt.Sprintf("ValidateTemporarySession: Failed sending request: %s", err.Error()))
		return ValidatedTemporarySession{}, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, err.Error())
	}

	logger.Info(fmt.Sprintf("ValidateTemporarySession: Finished execution"))

	return model, nil
}

func InvalidateTemporarySession(sessionUuid string) *appErrors.Error {
	logger.Info(fmt.Sprintf("InvalidateTemporarySession: Starting execution"))

	body, err := createBody(sessionUuid)
	if err != nil {
		logger.Warn(fmt.Sprintf("InvalidateTemporarySession: Failed marshaling body to get code block: %s", err.Error()))
		return err
	}

	client, clientErr := createClient()
	if clientErr != nil {
		return clientErr
	}

	r, err := createRequest(fmt.Sprintf("%s/session/invalidate", CreateApiUrl()), "POST", body, nil)
	if err != nil {
		logger.Warn(fmt.Sprintf("InvalidateTemporarySession: Failed creating request body: %s", err.Error()))
		return err
	}

	var model ValidatedTemporarySession
	if err := newClient.SendWithBackoff(r, client, &model); err != nil {
		logger.Warn(fmt.Sprintf("InvalidateTemporarySession: Failed sending request: %s", err.Error()))
		return appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, err.Error())
	}

	logger.Info(fmt.Sprintf("InvalidateTemporarySession: Finished execution"))

	return nil
}

func GetAllFileContent(codeProjectUuid string) ([]*FileContent, *appErrors.Error) {
	logger.Info(fmt.Sprintf("GetAllFileContent: Starting execution"))

	client, clientErr := createClient()
	if clientErr != nil {
		return nil, clientErr
	}

	r, err := createRequest(fmt.Sprintf("%s/code-project/file/content/%s", CreateApiUrl(), codeProjectUuid), "GET", nil, nil)
	if err != nil {
		logger.Warn(fmt.Sprintf("GetAllFileContent: Failed creating request body: %s", err.Error()))
		return nil, err
	}

	var model []*FileContent
	if err := newClient.SendWithBackoff(r, client, &model); err != nil {
		logger.Warn(fmt.Sprintf("GetAllFileContent: Failed sending request: %s", err.Error()))
		return nil, appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, err.Error())
	}

	logger.Info(fmt.Sprintf("GetAllFileContent: Finished execution"))

	return model, nil
}
