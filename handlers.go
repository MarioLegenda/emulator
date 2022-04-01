package main

import (
	"net/http"
	"therebelsource/emulator/httpUtil"
	"therebelsource/emulator/linkedProjectExecution"
	"therebelsource/emulator/projectExecution"
	"therebelsource/emulator/runner"
	"therebelsource/emulator/singleFileExecution"
	"therebelsource/emulator/staticTypes"
)

func getEnvironmentsHandler(w http.ResponseWriter, r *http.Request) {
	cr := httpUtil.InitCurrentRequest(w, r)

	var languages []runner.Language

	languages = append(languages, runner.Node14)
	languages = append(languages, runner.CSharpMono)
	languages = append(languages, runner.NodeEsm)
	languages = append(languages, runner.NodeLts)
	languages = append(languages, runner.Haskell)
	languages = append(languages, runner.CLang)
	languages = append(languages, runner.CPlus)
	languages = append(languages, runner.GoLang)
	languages = append(languages, runner.Python2)
	languages = append(languages, runner.Python3)
	languages = append(languages, runner.Ruby)
	languages = append(languages, runner.Php74)
	languages = append(languages, runner.Rust)

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
