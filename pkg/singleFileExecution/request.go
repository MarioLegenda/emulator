package singleFileExecution

import (
	repository2 "emulator/pkg/repository"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/microcosm-cc/bluemonday"
)

type SingleFileRunRequest struct {
	Uuid string `json:"uuid"`

	codeBlock                 *repository2.CodeBlock
	validatedTemporarySession repository2.ValidatedTemporarySession
}

type SnippetRequest struct {
	Uuid string `json:"uuid"`

	snippet                   *repository2.Snippet
	validatedTemporarySession repository2.ValidatedTemporarySession
}

type PublicSnippetRequest struct {
	Uuid string `json:"uuid"`
	Text string `json:"text"`

	snippet                   *repository2.Snippet
	validatedTemporarySession repository2.ValidatedTemporarySession
}

type PublicSingleFileRunRequest struct {
	Uuid string `json:"uuid"`
	Text string `json:"text"`

	codeBlock                 *repository2.CodeBlock
	validatedTemporarySession repository2.ValidatedTemporarySession
}

func (l *SingleFileRunRequest) Sanitize() {
	p := bluemonday.StrictPolicy()

	l.Uuid = p.Sanitize(l.Uuid)
}

func (l *PublicSingleFileRunRequest) Sanitize() {
	p := bluemonday.StrictPolicy()

	l.Uuid = p.Sanitize(l.Uuid)
}

func (l *SingleFileRunRequest) Validate() error {
	if err := validation.ValidateStruct(l,
		validation.Field(&l.Uuid, validation.Required, is.UUID),
	); err != nil {
		return err
	}

	blockExists := func(request interface{}) error {
		sessionUuid := request.(string)

		validatedSession, err := repository2.ValidateTemporarySession(sessionUuid)

		if err != nil {
			return errors.New("Code block does not exist")
		}

		sessionData, err := repository2.GetCodeBlock(validatedSession.Session, sessionUuid)

		if err != nil {
			return errors.New("Code block does not exist")
		}

		go repository2.InvalidateTemporarySession(sessionUuid)

		l.codeBlock = sessionData
		l.validatedTemporarySession = validatedSession

		return nil
	}

	if err := validation.Validate(map[string]interface{}{
		"blockExists": l.Uuid,
	}, validation.Map(
		validation.Key("blockExists", validation.By(blockExists)),
	)); err != nil {
		return err
	}

	return nil
}

func (l *PublicSingleFileRunRequest) Validate() error {
	if err := validation.ValidateStruct(l,
		validation.Field(&l.Uuid, validation.Required, is.UUID),
		validation.Field(&l.Text, validation.When(l.Text != "", validation.RuneLength(0, 5000))),
	); err != nil {
		return err
	}

	blockExists := func(request interface{}) error {
		sessionUuid := request.(string)

		validatedSession, err := repository2.ValidateTemporarySession(sessionUuid)

		if err != nil {
			return errors.New("Code block does not exist")
		}

		sessionData, err := repository2.GetAnonymousCodeBlock(sessionUuid)

		if err != nil {
			return errors.New("Code block does not exist")
		}

		sessionData.Text = l.Text

		go repository2.InvalidateTemporarySession(sessionUuid)

		l.codeBlock = sessionData
		l.validatedTemporarySession = validatedSession

		return nil
	}

	if err := validation.Validate(map[string]interface{}{
		"blockExists": l.Uuid,
	}, validation.Map(
		validation.Key("blockExists", validation.By(blockExists)),
	)); err != nil {
		return err
	}

	return nil
}

func (l *SnippetRequest) Sanitize() {
	p := bluemonday.StrictPolicy()

	l.Uuid = p.Sanitize(l.Uuid)
}

func (l *SnippetRequest) Validate() error {
	if err := validation.ValidateStruct(l,
		validation.Field(&l.Uuid, validation.Required, is.UUID),
	); err != nil {
		return err
	}

	blockExists := func(request interface{}) error {
		sessionUuid := request.(string)

		session, err := repository2.ValidateTemporarySession(sessionUuid)

		if err != nil {
			return errors.New("Snippet does not exist")
		}

		sessionData, err := repository2.GetAuthenticatedSnippet(session.Session, sessionUuid)

		if err != nil {
			return errors.New("Snippet does not exist")
		}

		go repository2.InvalidateTemporarySession(sessionUuid)

		l.snippet = sessionData
		l.validatedTemporarySession = session

		return nil
	}

	if err := validation.Validate(map[string]interface{}{
		"blockExists": l.Uuid,
	}, validation.Map(
		validation.Key("blockExists", validation.By(blockExists)),
	)); err != nil {
		return err
	}

	return nil
}

func (l *PublicSnippetRequest) Sanitize() {
	p := bluemonday.StrictPolicy()

	l.Uuid = p.Sanitize(l.Uuid)
}

func (l *PublicSnippetRequest) Validate() error {
	if err := validation.ValidateStruct(l,
		validation.Field(&l.Uuid, validation.Required, is.UUID),
	); err != nil {
		return err
	}

	blockExists := func(request interface{}) error {
		sessionUuid := request.(string)

		session, err := repository2.ValidateTemporarySession(sessionUuid)

		if err != nil {
			return errors.New("Snippet does not exist")
		}

		sessionData, err := repository2.GetAnonymousSnippet(sessionUuid)

		if err != nil {
			return errors.New("Snippet does not exist")
		}

		go repository2.InvalidateTemporarySession(sessionUuid)

		l.snippet = sessionData
		l.validatedTemporarySession = session

		return nil
	}

	if err := validation.Validate(map[string]interface{}{
		"blockExists": l.Uuid,
	}, validation.Map(
		validation.Key("blockExists", validation.By(blockExists)),
	)); err != nil {
		return err
	}

	return nil
}
