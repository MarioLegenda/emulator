package singleFileExecution

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/microcosm-cc/bluemonday"
	"strings"
	"therebelsource/emulator/repository"
)

type SingleFileRunRequest struct {
	PageUuid           string          `json:"pageUuid"`
	BlockUuid        string          `json:"blockUuid"`
	Type             string          `json:"type"`
	State            string			 `json:"state"`

	codeBlock *repository.CodeBlock
}

func (l *SingleFileRunRequest) Sanitize() {
	p := bluemonday.StrictPolicy()

	l.PageUuid = p.Sanitize(l.PageUuid)
	l.BlockUuid = p.Sanitize(l.BlockUuid)
	l.Type = p.Sanitize(l.Type)
	l.State = p.Sanitize(l.State)
}

func (l *SingleFileRunRequest) Validate() error {
	if err := validation.ValidateStruct(l,
		validation.Field(&l.PageUuid, validation.Required, is.UUID),
		validation.Field(&l.BlockUuid, validation.Required, is.UUID),
	); err != nil {
		return err
	}

	blockExists := func(request interface{}) error {
		data := request.(struct{
			pageUuid string
			blockUuid string
			ksType string
		})
		repo := repository.InitRepository()

		codeBlock, err := repo.GetCodeBlock(data.pageUuid, data.blockUuid, data.ksType)

		if err != nil {
			return errors.New(fmt.Sprintf("Code block %s to be executed does not exist", data.blockUuid))
		}

		l.codeBlock = codeBlock

		return nil
	}

	typeValid := func(request interface{}) error {
		t := request.(string)

		validTypes := []string{"blog", "documentation", "book"}

		for _, k := range validTypes {
			if k == t {
				return nil
			}
		}

		return errors.New(fmt.Sprintf("Invalid type. Valid types are: %s", strings.Join(validTypes, ",")))
	}

	stateValid := func(request interface{}) error {
		t := request.(string)

		validStates := []string{"dev", "prod", "session", "single_file"}

		for _, k := range validStates {
			if k == t {
				return nil
			}
		}

		return errors.New(fmt.Sprintf("Invalid state. Valid states are: %s", strings.Join(validStates, ",")))
	}

	if err := validation.Validate(map[string]interface{}{
		"blockExists": struct {
			pageUuid string
			blockUuid string
			ksType string
		}{
			pageUuid: l.PageUuid,
			blockUuid: l.BlockUuid,
			ksType: l.Type,
		},
		"typeValid": l.Type,
		"stateValid": l.State,
	}, validation.Map(
		validation.Key("blockExists", validation.By(blockExists)),
		validation.Key("typeValid", validation.By(typeValid)),
		validation.Key("stateValid", validation.By(stateValid)),
	)); err != nil {
		return err
	}

	return nil
}
