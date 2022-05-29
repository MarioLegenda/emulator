package repository

type Token struct {
	ExpiresAt int64 `json:"expiresAt"`
}

type Session struct {
	Uuid        string   `json:"uuid"`
	Tokens      [3]Token `json:"tokens"`
	Persistent  bool     `json:"persistent"`
	Device      string   `json:"device"`
	Type        string   `json:"type"`
	AccountUuid string   `json:"AccountUuid"`
}

type ActiveSession struct {
	Session Session `json:"session"`
	Account Account `json:"account"`
}

type Account struct {
	Uuid     string `json:"uuid"`
	Name     string `json:"name"`
	LastName string `json:"lastName"`
	Email    string `json:"email"`
	Type     string `json:"type"`
	Provider string `json:"provider"`

	Confirmed bool `json:"confirmed"`

	CreatedAt int64  `json:"createdAt"`
	UpdatedAt *int64 `json:"updatedAt"`
}

type TemporarySession struct {
	Uuid        string                 `json:"uuid"`
	Device      string                 `json:"device"`
	Purpose     string                 `json:"purpose"`
	Permissions []string               `json:"permissions"`
	Data        map[string]interface{} `json:"data"`
}

type ValidatedTemporarySession struct {
	Timeout int `json:"timeout"`
}

type CodeBlock struct {
	Uuid      string  `json:"uuid"`
	PageUuid  string  `json:"pageUuid"`
	Position  int32   `json:"position"`
	BlockType string  `json:"blockType"`
	CreatedAt *string `json:"createdAt"`
	UpdatedAt *string `json:"updatedAt"`

	Text string `json:"text"`

	IsGist      bool   `json:"isGist"`
	IsCode      bool   `json:"isCode"`
	Readonly    bool   `json:"readonly"`
	PackageName string `json:"packageName"`

	GistData *GistData `json:"gistData"`
	Emulator *Language `json:"emulator"`

	CodeProjectUuid *string `json:"codeProjectUuid"`
	CodeResult      *string `json:"codeResult"`
}

type GistData struct {
	Username string `json:"username"`
	GistId   string `json:"gistId"`
}

type CodeProject struct {
	Uuid           string    `json:"uuid"`
	ShortId        string    `json:"shortId"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Environment    *Language `json:"environment"`
	Structure      []*File   `json:"structure"`
	StructureCount int       `json:"structureCount"`
	PackageName    string    `json:"packageName"`

	RootDirectory *File `json:"rootDirectory"`

	CreatedAt int64  `json:"createdAt"`
	UpdatedAt *int64 `json:"updatedAt"`
}

type File struct {
	Name string `json:"name" bson:"name"`

	IsRoot bool `json:"isRoot" bson:"isRoot"`
	Depth  int  `json:"depth" bson:"depth"`
	IsFile bool `json:"isFile" bson:"isFile"`

	Uuid string `json:"uuid" bson:"uuid"`

	Parent   *string  `json:"parent" bson:"parent"`
	Children []string `json:"children" bson:"children"`

	CreatedAt int64  `json:"createdAt" bson:"createdAt"`
	UpdatedAt *int64 `json:"updatedAt" bson:"updatedAt"`
}

type FileContent struct {
	CodeProjectUuid string `json:"codeProjectUuid" bson:"codeProjectUuid"`
	Uuid            string `json:"uuid" bson:"uuid"`
	Content         string `json:"content" bson:"content"`
}

type SessionCodeProjectData struct {
	CodeProject   *CodeProject   `json:"codeProject"`
	Content       []*FileContent `json:"fileContent"`
	ExecutingFile *File          `json:"executingFile"`
	PackageName   string         `json:"packageName"`
}

type SingleFileSessionData struct {
	CodeBlock *CodeBlock `json:"codeBlock"`
}

type LinkedSessionData struct {
	CodeProject *CodeProject   `json:"codeProject"`
	CodeBlock   *CodeBlock     `json:"codeBlock"`
	PackageName string         `json:"packageName"`
	Content     []*FileContent `json:"fileContent"`
}
