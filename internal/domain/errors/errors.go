package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Error struct {
	Code            Code
	Errors          []string
	Map             map[string]any
	OriginFunc      string
	Cause           string
	FriendlyMessage string
}

type InputError Error

func (e *Error) Error() string {
	joinErrors := strings.Join(e.Errors, " ")
	m, _ := json.Marshal(e.Map)
	return fmt.Sprintf("{"+
		"code: %s, "+
		"errors: %v, "+
		"map: %v, "+
		"originFunc: %s, "+
		"cause: %s, "+
		"friendlyMessage: %s"+
		"}",
		e.Code, joinErrors, string(m), e.OriginFunc, e.Cause, e.FriendlyMessage)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func New(inputError InputError) *Error {
	return &Error{
		Code:            inputError.Code,
		Errors:          inputError.Errors,
		OriginFunc:      inputError.OriginFunc,
		FriendlyMessage: inputError.FriendlyMessage,
	}
}
