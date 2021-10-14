package runner

import "time"

type Name string
type Text string
type Tag string
type InDevelopment bool
type InMaintenance bool

type Language struct {
	Name           Name          `json:"name"`
	Text           Text          `json:"text"`
	Tag            Tag           `json:"tag"`
	InDevelopment  InDevelopment `json:"inDevelopment"`
	InMaintenance  InMaintenance `json:"inMaintenance"`
	Language       string        `json:"language"`
	Extension      string        `json:"extension"`
	Output         string        `json:"output"`
	DefaultTimeout time.Duration `json:"defaultTimeout"`
	PackageTimeout time.Duration `json:"packageTimeout"`
}

var Node14 = Language{
	Name:           "node_v14_x",
	Text:           "Javascript (Node v14.x)",
	Tag:            "node:node_v14_x",
	InDevelopment:  false,
	InMaintenance:  false,
	Language:       "javascript",
	Extension:      "js",
	Output:         "",
	DefaultTimeout: 10 * time.Second,
	PackageTimeout: 60 * time.Second,
}

var NodeLts = Language{
	Name:           "node_latest",
	Text:           "Javascript (Node latest)",
	Tag:            "node:node_latest",
	InDevelopment:  false,
	InMaintenance:  false,
	Language:       "javascript",
	Extension:      "js",
	Output:         "",
	DefaultTimeout: 10 * time.Second,
	PackageTimeout: 60 * time.Second,
}

var GoLang = Language{
	Name:           "go",
	Text:           "Go v1.*.*",
	Tag:            "go:go_v1_14_2",
	InDevelopment:  false,
	InMaintenance:  false,
	Language:       "go",
	Extension:      "go",
	Output:         "",
	DefaultTimeout: 15 * time.Second,
	PackageTimeout: 60 * time.Second,
}

var Python2 = Language{
	Name:           "python2",
	Text:           "Python2",
	Tag:            "python:python2",
	InDevelopment:  false,
	InMaintenance:  false,
	Language:       "python",
	Extension:      "py",
	Output:         "",
	DefaultTimeout: 10 * time.Second,
	PackageTimeout: 60 * time.Second,
}

var Python3 = Language{
	Name:           "python3",
	Text:           "Python3",
	Tag:            "python:python3",
	InDevelopment:  false,
	InMaintenance:  false,
	Language:       "python",
	Extension:      "py",
	Output:         "",
	DefaultTimeout: 10 * time.Second,
	PackageTimeout: 60 * time.Second,
}

var Ruby = Language{
	Name:           "ruby",
	Text:           "Ruby 2.5.1",
	Tag:            "ruby:ruby",
	InDevelopment:  false,
	InMaintenance:  false,
	Language:       "ruby",
	Extension:      "rb",
	Output:         "",
	DefaultTimeout: 10 * time.Second,
	PackageTimeout: 60 * time.Second,
}

var Php74 = Language{
	Name:           "php74",
	Text:           "PHP 7.4",
	Tag:            "php:php7.4",
	InDevelopment:  false,
	InMaintenance:  false,
	Language:       "php",
	Extension:      "php",
	Output:         "",
	DefaultTimeout: 10 * time.Second,
	PackageTimeout: 30 * time.Second,
}

var Rust = Language{
	Name:           "rust",
	Text:           "Rust",
	Tag:            "rust:rust",
	InDevelopment:  false,
	InMaintenance:  false,
	Language:       "rust",
	Extension:      "rs",
	Output:         "",
	DefaultTimeout: 15 * time.Second,
	PackageTimeout: 30 * time.Second,
}

var Haskell = Language{
	Name:           "haskell",
	Text:           "Haskell",
	Tag:            "haskell:haskell",
	InDevelopment:  false,
	InMaintenance:  false,
	Language:       "haskell",
	Extension:      "hs",
	Output:         "",
	DefaultTimeout: 15 * time.Second,
	PackageTimeout: 30 * time.Second,
}

var CLang = Language{
	Name:           "c",
	Text:           "C",
	Tag:            "c:c",
	InDevelopment:  false,
	InMaintenance:  false,
	Language:       "c",
	Extension:      "c",
	Output:         "",
	DefaultTimeout: 15 * time.Second,
	PackageTimeout: 30 * time.Second,
}

var CPlus = Language{
	Name:           "c++",
	Text:           "C++",
	Tag:            "c-plus:c-plus",
	InDevelopment:  false,
	InMaintenance:  false,
	Language:       "c",
	Extension:      "cpp",
	Output:         "",
	DefaultTimeout: 15 * time.Second,
	PackageTimeout: 30 * time.Second,
}

type SingleFileBuildResult struct {
	DirectoryName string
	ExecutionDirectory string
	FileName  string
	Environment Language
	StateDirectory string
}
