package http

import (
	httpUtil2 "emulator/pkg/httpUtil"
	"emulator/pkg/linkedProjectExecution"
	"emulator/pkg/projectExecution"
	"emulator/pkg/repository"
	"emulator/pkg/singleFileExecution"
	"emulator/pkg/staticTypes"
	"net/http"
)

func getEnvironmentsHandler(w http.ResponseWriter, r *http.Request) {
	cr := httpUtil2.InitCurrentRequest(w, r)

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
	languages = append(languages, repository.Julia)
	languages = append(languages, repository.Ruby)
	languages = append(languages, repository.Php74)
	languages = append(languages, repository.Rust)
	languages = append(languages, repository.PerlLts)
	languages = append(languages, repository.Lua)

	apiResponse := httpUtil2.CreateSuccessResponse(cr, staticTypes.RESPONSE_RESOURCE, languages, http.StatusOK, "An instance of file content")

	cr.SendResponse(apiResponse)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	cr := httpUtil2.InitCurrentRequest(w, r)

	apiResponse := httpUtil2.CreateSuccessResponse(cr, staticTypes.RESPONSE_RESOURCE, nil, http.StatusOK, "Server healthy")

	cr.SendResponse(apiResponse)
}

func executeSingleCodeBlockHandler(w http.ResponseWriter, r *http.Request) {
	cr := httpUtil2.InitCurrentRequest(w, r)

	requestModel := cr.ReadSingleFileRunRequest()

	if requestModel == nil {
		return
	}

	runResult, err := singleFileExecution.SingleFileExecutionService.RunSingleFile(requestModel)

	if err != nil {
		apiResponse := httpUtil2.CreateErrorResponse(cr, err, err.GetData())

		cr.SendResponse(apiResponse)

		return
	}

	apiResponse := httpUtil2.CreateSuccessResponse(cr, staticTypes.RESPONSE_RESOURCE, runResult, http.StatusOK, "Emulator run result")

	cr.SendResponse(apiResponse)
}

func executeSnippet(w http.ResponseWriter, r *http.Request) {
	cr := httpUtil2.InitCurrentRequest(w, r)

	requestModel := cr.ReadSnippetRequest()

	if requestModel == nil {
		return
	}

	runResult, err := singleFileExecution.SingleFileExecutionService.RunSnippet(requestModel)

	if err != nil {
		apiResponse := httpUtil2.CreateErrorResponse(cr, err, err.GetData())

		cr.SendResponse(apiResponse)

		return
	}

	apiResponse := httpUtil2.CreateSuccessResponse(cr, staticTypes.RESPONSE_RESOURCE, runResult, http.StatusOK, "Emulator run result")

	cr.SendResponse(apiResponse)
}

func executePublicSnippet(w http.ResponseWriter, r *http.Request) {
	cr := httpUtil2.InitCurrentRequest(w, r)

	requestModel := cr.ReadPublicSnippetRequest()

	if requestModel == nil {
		return
	}

	runResult, err := singleFileExecution.SingleFileExecutionService.RunPublicSnippet(requestModel)

	if err != nil {
		apiResponse := httpUtil2.CreateErrorResponse(cr, err, err.GetData())

		cr.SendResponse(apiResponse)

		return
	}

	apiResponse := httpUtil2.CreateSuccessResponse(cr, staticTypes.RESPONSE_RESOURCE, runResult, http.StatusOK, "Emulator run result")

	cr.SendResponse(apiResponse)
}

func executePublicSingleFileRunResult(w http.ResponseWriter, r *http.Request) {
	cr := httpUtil2.InitCurrentRequest(w, r)

	requestModel := cr.ReadPublicSingleFileRunResult()

	if requestModel == nil {
		return
	}

	runResult, err := singleFileExecution.SingleFileExecutionService.RunPublicSingleFile(requestModel)

	if err != nil {
		apiResponse := httpUtil2.CreateErrorResponse(cr, err, err.GetData())

		cr.SendResponse(apiResponse)

		return
	}

	apiResponse := httpUtil2.CreateSuccessResponse(cr, staticTypes.RESPONSE_RESOURCE, runResult, http.StatusOK, "Emulator run result")

	cr.SendResponse(apiResponse)
}

func executeProjectHandler(w http.ResponseWriter, r *http.Request) {
	cr := httpUtil2.InitCurrentRequest(w, r)

	requestModel := cr.ReadProjectExecutionRequest()

	if requestModel == nil {
		return
	}

	runResult, err := projectExecution.ProjectExecutionService.RunProject(requestModel)

	if err != nil {
		apiResponse := httpUtil2.CreateErrorResponse(cr, err, err.GetData())

		cr.SendResponse(apiResponse)

		return
	}

	apiResponse := httpUtil2.CreateSuccessResponse(cr, staticTypes.RESPONSE_RESOURCE, runResult, http.StatusOK, "Emulator run result")

	cr.SendResponse(apiResponse)
}

func executeLinkedProjectHandler(w http.ResponseWriter, r *http.Request) {
	cr := httpUtil2.InitCurrentRequest(w, r)

	requestModel := cr.ReadLinkedProjectExecutionRequest()

	if requestModel == nil {
		return
	}

	runResult, err := linkedProjectExecution.ExecutionService.RunProject(requestModel)

	if err != nil {
		apiResponse := httpUtil2.CreateErrorResponse(cr, err, err.GetData())

		cr.SendResponse(apiResponse)

		return
	}

	apiResponse := httpUtil2.CreateSuccessResponse(cr, staticTypes.RESPONSE_RESOURCE, runResult, http.StatusOK, "Emulator run result")

	cr.SendResponse(apiResponse)
}

func executePublicLinkedProjectHandler(w http.ResponseWriter, r *http.Request) {
	cr := httpUtil2.InitCurrentRequest(w, r)

	requestModel := cr.ReadPublicLinkedProjectExecution()

	if requestModel == nil {
		return
	}

	runResult, err := linkedProjectExecution.ExecutionService.RunPublicProject(requestModel)

	if err != nil {
		apiResponse := httpUtil2.CreateErrorResponse(cr, err, err.GetData())

		cr.SendResponse(apiResponse)

		return
	}

	apiResponse := httpUtil2.CreateSuccessResponse(cr, staticTypes.RESPONSE_RESOURCE, runResult, http.StatusOK, "Emulator run result")

	cr.SendResponse(apiResponse)
}
