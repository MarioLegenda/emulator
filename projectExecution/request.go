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
	Type             string          `json:"type"`

	codeProject *repository.CodeProject
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
		"typeValid": l.Type,
	}, validation.Map(
		validation.Key("codeProjectExists", validation.By(codeProjectExists)),
		validation.Key("typeValid", validation.By(typeValid)),
	)); err != nil {
		return err
	}

	return nil
}
