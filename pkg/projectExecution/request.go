package projectExecution

import (
	repository2 "emulator/pkg/repository"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/microcosm-cc/bluemonday"
)

type ProjectRunRequest struct {
	Uuid     string `json:"uuid"`
	FileUuid string `json:"fileUuid"`

	sessionData               *repository2.SessionCodeProjectData
	validatedTemporarySession repository2.ValidatedTemporarySession
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

		session, err := repository2.ValidateTemporarySession(sessionUuid)

		if err != nil {
			return errors.New("Project does not exist")
		}

		sessionData, err := repository2.GetProjectSessionData(session.Session, sessionUuid)

		if err != nil {
			return errors.New("Project does not exists")
		}

		go repository2.InvalidateTemporarySession(sessionUuid)

		l.validatedTemporarySession = session
		l.sessionData = sessionData

		return nil
	}

	if err := validation.Validate(map[string]interface{}{
		"sessionValid": l.Uuid,
	}, validation.Map(
		validation.Key("sessionValid", validation.By(sessionValid)),
	)); err != nil {
		return err
	}

	return nil
}
