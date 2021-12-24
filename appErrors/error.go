package appErrors

import (
	"fmt"
	"os"
)

/**
Master codes
*/

const ApplicationError = 1
const ServerError = 2

/**
Application codes
*/

const ApplicationRuntimeError = 1
const NotFoundError = 2
const FilesystemError = 3

type AppError map[string]string

func ConstructError(masterCode int, appCode int, message string) string {
	return fmt.Sprintf("An appErrors occurred with MasterCode: %d; ApplicationCode: %d; Message: %s", masterCode, appCode, message)
}

func TerminateWithMessage(message string) {
	fmt.Println(message)

	os.Exit(1)
}

type Error struct {
	MasterCode int
	Code       int
	Message    string
	Data       AppError
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) AddData(k string, m string) {
	e.Data[k] = m
}

func (e *Error) GetData() AppError {
	if len(e.Data) == 0 {
		return nil
	}

	return e.Data
}

func New(masterCode int, appCode int, msg string) *Error {
	return &Error{
		MasterCode: masterCode,
		Code:       appCode,
		Message:    ConstructError(masterCode, appCode, msg),
		Data:       make(map[string]string),
	}
}
