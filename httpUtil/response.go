package httpUtil

import (
	"encoding/json"
	"log"
	"net/http"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/staticTypes"
)

type ResponsePagination struct {
	page int
	size int
}

type ApiResponse struct {
	Data       interface{}         `json:"data"`
	Method     string              `json:"method"`
	Status     int                 `json:"status"`
	Type       string              `json:"type"`
	Pagination *ResponsePagination `json:"pagination"`
	Message    string              `json:"message"`
	MasterCode int                 `json:"masterCode"`
	Code       int                 `json:"code"`
	OriginUrl  string              `json:"originUrl"`
}

func (a *ApiResponse) ToJson() []byte {
	b, err := json.Marshal(a)

	if err != nil {
		log.Fatalln(err)
	}

	return b
}

func CreateErrorResponse(cr CurrentHttpRequest, err *appErrors.Error, data appErrors.AppError) *ApiResponse {
	apiResponse := &ApiResponse{}
	status := 400

	if err.Code == appErrors.NotFoundError {
		status = http.StatusNotFound
	}

	apiResponse.Method = cr.r.Method
	apiResponse.MasterCode = err.MasterCode
	apiResponse.Code = err.Code
	apiResponse.Message = err.Message
	apiResponse.Type = staticTypes.RESPONSE_ERROR
	apiResponse.Status = status
	apiResponse.OriginUrl = cr.r.URL.String()
	apiResponse.Data = data

	return apiResponse
}

func CreateSuccessResponse(cr CurrentHttpRequest, t string, data interface{}, status int, msg string) *ApiResponse {
	apiResponse := &ApiResponse{}

	apiResponse.Method = cr.r.Method
	apiResponse.Type = t
	apiResponse.OriginUrl = cr.r.URL.String()
	apiResponse.Data = data
	apiResponse.Status = status
	apiResponse.Message = msg

	return apiResponse
}
