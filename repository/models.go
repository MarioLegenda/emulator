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

	GistData *GistData `json:"gistData"`
	Emulator runner.Language `json:"emulator"`

	CodeProjectUuid *string `json:"codeProjectUuid"`
	CodeResult      *string `json:"codeResult"`
}

type GistData struct {
	Username string `json:"username"`
	GistId   string `json:"gistId"`
}

type CodeProject struct {
	Uuid        string    `json:"uuid" bson:"uuid"`
	ShortId     string    `json:"shortId" bson:"shortId"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Environment *runner.Language `json:"environment" bson:"environment"`
	Structure []*File `json:"structure" bson:"structure"`
	StructureCount int `json:"structureCount" bson:"structureCount"`

	RootDirectory *File `json:"rootDirectory" bson:"-"`

	CreatedAt *string `json:"createdAt" bson:"createdAt"`
	UpdatedAt *string `json:"updatedAt" bson:"updatedAt"`
}

type File struct {
	Name     string   `json:"name" bson:"name"`

	IsRoot bool `json:"isRoot" bson:"isRoot"`
	Depth  int  `json:"depth" bson:"depth"`
	IsFile bool `json:"isFile" bson:"isFile"`

	Uuid string `json:"uuid" bson:"uuid"`

	Parent   *string  `json:"parent" bson:"parent"`
	Children []string `json:"children" bson:"children"`

	CreatedAt *string `json:"createdAt" bson:"createdAt"`
	UpdatedAt *string `json:"updatedAt" bson:"updatedAt"`
}

type FileContent struct {
	CodeProjectUuid string `json:"codeProjectUuid" bson:"codeProjectUuid"`
	Uuid string `json:"uuid" bson:"uuid"`
	Content string `json:"content" bson:"content"`
}

