package linkedProjectExecution

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/microcosm-cc/bluemonday"
	"therebelsource/emulator/repository"
)

type LinkedProjectRunRequest struct {
	Uuid string `json:"uuid"`

	sessionData               *repository.LinkedSessionData
	validatedTemporarySession repository.ValidatedTemporarySession
}

func (l *LinkedProjectRunRequest) Sanitize() {
	p := bluemonday.StrictPolicy()

	l.Uuid = p.Sanitize(l.Uuid)
}

type PublicLinkedProjectRunRequest struct {
	Uuid string `json:"uuid"`
	Text string `json:"text"`

	sessionData               *repository.LinkedSessionData
	validatedTemporarySession repository.ValidatedTemporarySession
}

func (l *PublicLinkedProjectRunRequest) Sanitize() {
	p := bluemonday.StrictPolicy()

	l.Uuid = p.Sanitize(l.Uuid)
	l.Text = p.Sanitize(l.Text)
}

func (l *PublicLinkedProjectRunRequest) Validate() error {
	if err := validation.ValidateStruct(l,
		validation.Field(&l.Uuid, validation.Required, is.UUID),
		validation.Field(&l.Text, validation.When(l.Text != "", validation.RuneLength(0, 5000))),
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

		sessionData, err := repo.GetLinkedSessionData(session.Session, sessionUuid)

		if err != nil {
			return errors.New("Project does not exists")
		}

		sessionData.CodeBlock.Text = l.Text

		go repo.InvalidateTemporarySession(session.Session, sessionUuid)

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

func (l *LinkedProjectRunRequest) Validate() error {
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

		sessionData, err := repo.GetLinkedSessionData(session.Session, sessionUuid)

		if err != nil {
			return errors.New("Project does not exists")
		}

		go repo.InvalidateTemporarySession(session.Session, sessionUuid)

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
