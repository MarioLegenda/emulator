package builders

import (
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/repository"
)

func createFsSystem(paths map[string][]*repository.File, contents []*repository.FileContent) *appErrors.Error {
	contentsMap := make(map[string]*repository.FileContent)

	for _, c := range contents {
		contentsMap[c.Uuid] = c
	}

	for path, files := range paths {
		if err := createDir(path); err != nil {
			return err
		}

		if files != nil && len(files) != 0 {
			for _, f := range files {
				if content, ok := contentsMap[f.Uuid]; ok {
					if err := writeContent(f.Name, path, content.Content); err != nil {
						return err
					}

					continue
				}

				content := &repository.FileContent{
					CodeProjectUuid: "",
					Uuid:            "",
					Content:         "",
				}

				if err := writeContent(f.Name, path, content.Content); err != nil {
					return err
				}
			}
		}
	}

	return nil
}