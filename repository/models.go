package repository

import "therebelsource/emulator/runner"

type CodeBlock struct {
	Uuid      string  `json:"uuid"`
	PageUuid  string  `json:"pageUuid"`
	Position  int32   `json:"position"`
	BlockType string  `json:"blockType"`
	CreatedAt *string `json:"createdAt"`
	UpdatedAt *string `json:"updatedAt"`

	Text string `json:"text"`

	IsGist   bool `json:"isGist"`
	IsCode   bool `json:"isCode"`
	Readonly bool `json:"readonly"`

	GistData *GistData              `json:"gistData"`
	Emulator runner.Language `json:"emulator"`

	CodeProjectUuid *string `json:"codeProjectUuid"`
	CodeResult      *string `json:"codeResult"`
}

type GistData struct {
	Username string `json:"username"`
	GistId   string `json:"gistId"`
}

