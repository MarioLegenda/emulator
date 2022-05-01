package runner

type Name string
type Text string
type Tag string
type InDevelopment bool
type InMaintenance bool

type Language struct {
	Name          Name          `json:"name"`
	Text          Text          `json:"text"`
	Tag           Tag           `json:"tag"`
	InDevelopment InDevelopment `json:"inDevelopment"`
	InMaintenance InMaintenance `json:"inMaintenance"`
	Language      string        `json:"language"`
	Extension     string        `json:"extension"`
}

var Node14 = Language{
	Name:          "node_v14_x",
	Text:          "Javascript (Node v14.x)",
	Tag:           "node:node_v14_x",
	InDevelopment: false,
	InMaintenance: false,
	Language:      "javascript",
	Extension:     "js",
}

var NodeLts = Language{
	Name:          "node_latest",
	Text:          "Javascript (Node latest)",
	Tag:           "node:node_latest",
	InDevelopment: false,
	InMaintenance: false,
	Language:      "javascript",
	Extension:     "js",
}

var NodeEsm = Language{
	Name:          "node_latest_esm",
	Text:          "Javascript (Node ESM)",
	Tag:           "node:node_latest_esm",
	InDevelopment: false,
	InMaintenance: false,
	Language:      "javascript",
	Extension:     "mjs",
}

var GoLang = Language{
	Name:          "go",
	Text:          "Go v1.*.*",
	Tag:           "go:go_v1_17_6",
	InDevelopment: false,
	InMaintenance: false,
	Language:      "go",
	Extension:     "go",
}

var Python2 = Language{
	Name:          "python2",
	Text:          "Python2",
	Tag:           "python:python2",
	InDevelopment: false,
	InMaintenance: false,
	Language:      "python",
	Extension:     "py",
}

var Python3 = Language{
	Name:          "python3",
	Text:          "Python3",
	Tag:           "python:python3",
	InDevelopment: false,
	InMaintenance: false,
	Language:      "python",
	Extension:     "py",
}

var Ruby = Language{
	Name:          "ruby",
	Text:          "Ruby 2.5.1",
	Tag:           "ruby:ruby",
	InDevelopment: false,
	InMaintenance: false,
	Language:      "ruby",
	Extension:     "rb",
}

var Php74 = Language{
	Name:          "php74",
	Text:          "PHP 7.4",
	Tag:           "php:php7.4",
	InDevelopment: false,
	InMaintenance: false,
	Language:      "php",
	Extension:     "php",
}

var Rust = Language{
	Name:          "rust",
	Text:          "Rust",
	Tag:           "rust:rust",
	InDevelopment: false,
	InMaintenance: false,
	Language:      "rust",
	Extension:     "rs",
}

var Haskell = Language{
	Name:          "haskell",
	Text:          "Haskell",
	Tag:           "haskell:haskell",
	InDevelopment: false,
	InMaintenance: false,
	Language:      "haskell",
	Extension:     "hs",
}

var CLang = Language{
	Name:          "c",
	Text:          "C",
	Tag:           "c:c",
	InDevelopment: false,
	InMaintenance: false,
	Language:      "c",
	Extension:     "c",
}

var CPlus = Language{
	Name:          "c++",
	Text:          "C++",
	Tag:           "c-plus:c-plus",
	InDevelopment: false,
	InMaintenance: false,
	Language:      "c",
	Extension:     "cpp",
}

var CSharpMono = Language{
	Name:          "c_sharp_mono",
	Text:          "C# (Mono)",
	Tag:           "c_sharp_mono:c_sharp_mono",
	InDevelopment: false,
	InMaintenance: false,
	Language:      "csharp",
	Extension:     "cs",
}

type SingleFileBuildResult struct {
	ContainerName      string
	DirectoryName      string
	ExecutionDirectory string
	FileName           string
	Environment        *Language
	StateDirectory     string
	Timeout            int
	Args               []string
}

type ProjectRunResult struct {
	Success bool   `json:"success"`
	Result  string `json:"result"`
	Timeout int    `json:"timeout"`
}
