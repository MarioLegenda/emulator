package builders

import (
	"therebelsource/emulator/repository"
)

type FileTraverse struct {
	Files []*repository.File
}

func (dth FileTraverse) CreatePaths() map[string][]*repository.File {
	paths := make(map[string][]*repository.File)
	mappedSystem := make(map[string]*repository.File)

	for _, file := range dth.Files {
		mappedSystem[file.Uuid] = file
	}

	return paths
}


