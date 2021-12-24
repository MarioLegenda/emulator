package httpUtil

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/projectExecution"
	"therebelsource/emulator/singleFileExecution"
)

type CurrentHttpRequest struct {
	w http.ResponseWriter
	r *http.Request
}

func InitCurrentRequest(w http.ResponseWriter, r *http.Request) CurrentHttpRequest {
	return CurrentHttpRequest{
		w: w,
		r: r,
	}
}

func (cr CurrentHttpRequest) SendResponse(response *ApiResponse) {
	cr.w.Header().Set("Content-Type", "application/json")
	cr.w.WriteHeader(response.Status)

	body := response.ToJson()

	_, err := cr.w.Write(body)

	if err != nil {
		log.Fatal(err)
	}
}

func (cr CurrentHttpRequest) ReadBodyOrSendError() []byte {
	if cr.r.Body == nil {
		requestErr := appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Invalid body could not be unpacked")
		apiResponse := CreateErrorResponse(cr, requestErr, nil)

		cr.SendResponse(apiResponse)

		return nil
	}

	body, err := ioutil.ReadAll(cr.r.Body)

	if err != nil {
		requestErr := appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, err.Error())
		apiResponse := CreateErrorResponse(cr, requestErr, nil)

		cr.SendResponse(apiResponse)

		return nil
	}

	return body
}

func (cr CurrentHttpRequest) ReadSingleFileRunRequest() *singleFileExecution.SingleFileRunRequest{
	body := cr.ReadBodyOrSendError()

	if body == nil {
		return nil
	}

	var m singleFileExecution.SingleFileRunRequest
	t := &m

	err := json.Unmarshal(body, t)

	if err != nil {
		requestErr := appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Request data is invalid")
		apiResponse := CreateErrorResponse(cr, requestErr, nil)

		cr.SendResponse(apiResponse)

		return nil
	}

	err = t.Validate()

	if err != nil {
		cr.sendValidationError(err)

		return nil
	}

	return t
}

func (cr CurrentHttpRequest) ReadProjectExecutionRequest() *projectExecution.ProjectRunRequest {
	body := cr.ReadBodyOrSendError()

	if body == nil {
		return nil
	}

	var m projectExecution.ProjectRunRequest
	t := &m

	err := json.Unmarshal(body, t)

	if err != nil {
		requestErr := appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Request data is invalid")
		apiResponse := CreateErrorResponse(cr, requestErr, nil)

		cr.SendResponse(apiResponse)

		return nil
	}

	err = t.Validate()

	if err != nil {
		cr.sendValidationError(err)

		return nil
	}

	return t
}

func (cr CurrentHttpRequest) sendValidationError(err error) {
	b, _ := json.Marshal(err)

	requestErr := appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Request data is invalid")

	var e appErrors.AppError
	json.Unmarshal(b, &e)

	apiResponse := CreateErrorResponse(cr, requestErr, e)

	cr.SendResponse(apiResponse)
}

