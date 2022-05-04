package singleFileExecution

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/microcosm-cc/bluemonday"
	"therebelsource/emulator/repository"
)

type SingleFileRunRequest struct {
	Uuid string `json:"uuid"`

	codeBlock                 *repository.CodeBlock
	validatedTemporarySession repository.ValidatedTemporarySession
}

type PublicSingleFileRunRequest struct {
	Uuid string `json:"uuid"`
	Text string `json:"text"`

	codeBlock                 *repository.CodeBlock
	validatedTemporarySession repository.ValidatedTemporarySession
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

		repo := repository.InitRepository()

		session, err := repo.ValidateTemporarySession(sessionUuid)

		if err != nil {
			return errors.New("Code block does not exist")
		}

		sessionData, err := repo.GetCodeBlock(sessionUuid)

		if err != nil {
			return errors.New("Code block does not exist")
		}

		go repo.InvalidateTemporarySession(sessionUuid)

		l.codeBlock = sessionData
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

func (l *PublicSingleFileRunRequest) Validate() error {
	if err := validation.ValidateStruct(l,
		validation.Field(&l.Uuid, validation.Required, is.UUID),
		validation.Field(&l.Text, validation.When(l.Text != "", validation.RuneLength(0, 5000))),
	); err != nil {
		return err
	}

	blockExists := func(request interface{}) error {
		sessionUuid := request.(string)

		repo := repository.InitRepository()

		session, err := repo.ValidateTemporarySession(sessionUuid)

		if err != nil {
			return errors.New("Code block does not exist")
		}

		sessionData, err := repo.GetCodeBlock(sessionUuid)

		if err != nil {
			return errors.New("Code block does not exist")
		}

		sessionData.Text = l.Text

		go repo.InvalidateTemporarySession(sessionUuid)

		l.codeBlock = sessionData
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
