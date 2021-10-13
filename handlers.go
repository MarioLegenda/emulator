package main

import (
	"net/http"
	"therebelsource/emulator/httpUtil"
	"therebelsource/emulator/singleFileExecution"
	"therebelsource/emulator/staticTypes"
)

func getEnvironmentsHandler(w http.ResponseWriter, r *http.Request) {
	cr := httpUtil.InitCurrentRequest(w, r)

	var languages []Language

	languages = append(languages, node12)
	languages = append(languages, nodeLts)
	languages = append(languages, haskell)
	languages = append(languages, c)
	languages = append(languages, cPlus)
	languages = append(languages, goLang)
	languages = append(languages, python2)
	languages = append(languages, python3)
	languages = append(languages, ruby)
	languages = append(languages, php74)
	languages = append(languages, rust)

	apiResponse := httpUtil.CreateSuccessResponse(cr, staticTypes.RESPONSE_RESOURCE, languages, http.StatusOK, "An instance of file content")

	cr.SendResponse(apiResponse)
}

func executeSingleCodeBlockHandler(w http.ResponseWriter, r *http.Request) {
	cr := httpUtil.InitCurrentRequest(w, r)

	requestModel := cr.ReadSingleFileRunRequest()

	if requestModel == nil {
		return
	}

	runResult, err := singleFileExecution.SingleFileExecutionService.RunSingleFile(requestModel)

	if err != nil {
		apiResponse := httpUtil.CreateErrorResponse(cr, err, err.GetData())

		cr.SendResponse(apiResponse)

		return
	}

	apiResponse := httpUtil.CreateSuccessResponse(cr, staticTypes.RESPONSE_RESOURCE, runResult, http.StatusOK, "Emulator run result")

	cr.SendResponse(apiResponse)
}

