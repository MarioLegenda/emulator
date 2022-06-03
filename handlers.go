package main

import (
	"net/http"
	"therebelsource/emulator/httpUtil"
	"therebelsource/emulator/linkedProjectExecution"
	"therebelsource/emulator/projectExecution"
	"therebelsource/emulator/repository"
	"therebelsource/emulator/singleFileExecution"
	"therebelsource/emulator/staticTypes"
)

func getEnvironmentsHandler(w http.ResponseWriter, r *http.Request) {
	cr := httpUtil.InitCurrentRequest(w, r)

	var languages []repository.Language

	languages = append(languages, repository.CSharpMono)
	languages = append(languages, repository.NodeEsm)
	languages = append(languages, repository.NodeLts)
	languages = append(languages, repository.Haskell)
	languages = append(languages, repository.CLang)
	languages = append(languages, repository.CPlus)
	languages = append(languages, repository.GoLang)
	languages = append(languages, repository.Python2)
	languages = append(languages, repository.Python3)
	languages = append(languages, repository.Ruby)
	languages = append(languages, repository.Php74)
	languages = append(languages, repository.Rust)

	apiResponse := httpUtil.CreateSuccessResponse(cr, staticTypes.RESPONSE_RESOURCE, languages, http.StatusOK, "An instance of file content")

	cr.SendResponse(apiResponse)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	cr := httpUtil.InitCurrentRequest(w, r)

	apiResponse := httpUtil.CreateSuccessResponse(cr, staticTypes.RESPONSE_RESOURCE, nil, http.StatusOK, "Server healthy")

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

func executePublicSingleFileRunResult(w http.ResponseWriter, r *http.Request) {
	cr := httpUtil.InitCurrentRequest(w, r)

	requestModel := cr.ReadPublicSingleFileRunResult()

	if requestModel == nil {
		return
	}

	runResult, err := singleFileExecution.SingleFileExecutionService.RunPublicSingleFile(requestModel)

	if err != nil {
		apiResponse := httpUtil.CreateErrorResponse(cr, err, err.GetData())

		cr.SendResponse(apiResponse)

		return
	}

	apiResponse := httpUtil.CreateSuccessResponse(cr, staticTypes.RESPONSE_RESOURCE, runResult, http.StatusOK, "Emulator run result")

	cr.SendResponse(apiResponse)
}

func executeProjectHandler(w http.ResponseWriter, r *http.Request) {
	cr := httpUtil.InitCurrentRequest(w, r)

	requestModel := cr.ReadProjectExecutionRequest()

	if requestModel == nil {
		return
	}

	runResult, err := projectExecution.ProjectExecutionService.RunProject(requestModel)

	if err != nil {
		apiResponse := httpUtil.CreateErrorResponse(cr, err, err.GetData())

		cr.SendResponse(apiResponse)

		return
	}

	apiResponse := httpUtil.CreateSuccessResponse(cr, staticTypes.RESPONSE_RESOURCE, runResult, http.StatusOK, "Emulator run result")

	cr.SendResponse(apiResponse)
}

func executeLinkedProjectHandler(w http.ResponseWriter, r *http.Request) {
	cr := httpUtil.InitCurrentRequest(w, r)

	requestModel := cr.ReadLinkedProjectExecutionRequest()

	if requestModel == nil {
		return
	}

	runResult, err := linkedProjectExecution.ExecutionService.RunProject(requestModel)

	if err != nil {
		apiResponse := httpUtil.CreateErrorResponse(cr, err, err.GetData())

		cr.SendResponse(apiResponse)

		return
	}

	apiResponse := httpUtil.CreateSuccessResponse(cr, staticTypes.RESPONSE_RESOURCE, runResult, http.StatusOK, "Emulator run result")

	cr.SendResponse(apiResponse)
}

func executePublicLinkedProjectHandler(w http.ResponseWriter, r *http.Request) {
	cr := httpUtil.InitCurrentRequest(w, r)

	requestModel := cr.ReadPublicLinkedProjectExecution()

	if requestModel == nil {
		return
	}

	runResult, err := linkedProjectExecution.ExecutionService.RunPublicProject(requestModel)

	if err != nil {
		apiResponse := httpUtil.CreateErrorResponse(cr, err, err.GetData())

		cr.SendResponse(apiResponse)

		return
	}

	apiResponse := httpUtil.CreateSuccessResponse(cr, staticTypes.RESPONSE_RESOURCE, runResult, http.StatusOK, "Emulator run result")

	cr.SendResponse(apiResponse)
}
