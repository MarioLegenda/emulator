package runners

import "therebelsource/emulator/appErrors"

type Result struct {
	Result  string
	Success bool
	Error   *appErrors.Error
}

type name string

type language struct {
	name name `json:"name"`
}

var nodeLts = language{
	name: "node_latest",
}

var nodeEsm = language{
	name: "node_latest_esm",
}

var goLang = language{
	name: "go",
}

var ruby = language{
	name: "ruby",
}

var php = language{
	name: "php74",
}

var python2 = language{
	name: "python2",
}

var python3 = language{
	name: "python3",
}

var julia = language{
	name: "julia",
}

var csharpMono = language{
	name: "c_sharp_mono",
}

var haskell = language{
	name: "haskell",
}

var cLang = language{
	name: "c",
}

var cPlus = language{
	name: "c++",
}

var rust = language{
	name: "rust",
}
