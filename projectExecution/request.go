package projectExecution

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/microcosm-cc/bluemonday"
	"strings"
	"therebelsource/emulator/repository"
)

type CodeProjectRunRequest struct {
	CodeProjectUuid           string          `json:"codeProjectUuid"`
	FileUuid                  string          `json:"fileUuid"`
	Type             string          `json:"type"`

	codeProject *repository.CodeProject
	executingFile *repository.File
}

func (l *CodeProjectRunRequest) Sanitize() {
	p := bluemonday.StrictPolicy()

	l.CodeProjectUuid = p.Sanitize(l.CodeProjectUuid)
	l.Type = p.Sanitize(l.Type)
}

func (l *CodeProjectRunRequest) Validate() error {
	if err := validation.ValidateStruct(l,
		validation.Field(&l.CodeProjectUuid, validation.Required, is.UUID),
	); err != nil {
		return err
	}

	codeProjectExists := func(request interface{}) error {
		codeProjectUuid := request.(string)

		repo := repository.InitRepository()

		codeProject, err := repo.GetCodeProject(codeProjectUuid)

		if err != nil {
			return errors.New(fmt.Sprintf("Code project %s to be executed does not exist", codeProjectUuid))
		}

		l.codeProject = codeProject

		return nil
	}

	fileExists := func(request interface{}) error {
		if l.codeProject != nil {
			data := request.(struct{
				fileUuid string
			})

			found := false
			for _, f := range l.codeProject.Structure {
				if f.Uuid == data.fileUuid {
					found = true
					
					l.executingFile = f

					break
				}
			}

			if !found {
				return errors.New(fmt.Sprintf("File to be executed %s does not exist", data.fileUuid))
			}
		}

		return nil
	}

	typeValid := func(request interface{}) error {
		t := request.(string)

		validTypes := []string{"session"}

		for _, k := range validTypes {
			if k == t {
				return nil
			}
		}

		return errors.New(fmt.Sprintf("Invalid type. Valid types are: %s", strings.Join(validTypes, ",")))
	}

	if err := validation.Validate(map[string]interface{} {
		"codeProjectExists": l.CodeProjectUuid,
		"fileExists": struct {
			fileUuid string
		}{
			fileUuid: l.FileUuid,
		},
		"typeValid": l.Type,
	}, validation.Map(
		validation.Key("codeProjectExists", validation.By(codeProjectExists)),
		validation.Key("fileExists", validation.By(fileExists)),
		validation.Key("typeValid", validation.By(typeValid)),
	)); err != nil {
		return err
	}

	return nil
}
