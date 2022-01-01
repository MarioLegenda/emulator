package projectExecution

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/microcosm-cc/bluemonday"
	"therebelsource/emulator/repository"
)

type ProjectRunRequest struct {
	Uuid     string `json:"uuid"`
	FileUuid string `json:"fileUuid"`

	sessionData               *repository.SessionCodeProjectData
	executingFile             *repository.File
	validatedTemporarySession repository.ValidatedTemporarySession
}

func (l *ProjectRunRequest) Sanitize() {
	p := bluemonday.StrictPolicy()

	l.Uuid = p.Sanitize(l.Uuid)
}

func (l *ProjectRunRequest) Validate() error {
	if err := validation.ValidateStruct(l,
		validation.Field(&l.Uuid, validation.Required, is.UUID),
	); err != nil {
		return err
	}

	sessionValid := func(request interface{}) error {
		sessionUuid := request.(string)

		repo := repository.InitRepository()

		session, err := repo.ValidateTemporarySession(sessionUuid)

		if err != nil {
			return errors.New("Project does not exist")
		}

		sessionData, err := repo.GetProjectSessionData(sessionUuid)

		if err != nil {
			return errors.New("Project does not exists")
		}

		if err := repo.InvalidateTemporarySession(sessionUuid); err != nil {
			return errors.New("Project does not exist")
		}

		l.validatedTemporarySession = session
		l.sessionData = sessionData

		return nil
	}

	fileExists := func(request interface{}) error {
		if l.sessionData != nil {
			fileUuid := request.(string)

			found := false
			for _, f := range l.sessionData.CodeProject.Structure {
				if f.Uuid == fileUuid {
					found = true

					l.executingFile = f

					break
				}
			}

			if !found {
				return errors.New(fmt.Sprintf("File to be executed %s does not exist", fileUuid))
			}
		}

		return nil
	}

	if err := validation.Validate(map[string]interface{}{
		"sessionValid": l.Uuid,
		"fileExists":   l.FileUuid,
	}, validation.Map(
		validation.Key("sessionValid", validation.By(sessionValid)),
		validation.Key("fileExists", validation.By(fileExists)),
	)); err != nil {
		return err
	}

	return nil
}
