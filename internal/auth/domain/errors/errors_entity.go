package errors

import (
	"fmt"

	"github.com/andreis3/auth-ms/internal/auth/domain/validator"
)

const ValidationCode Code = ErrBadRequest

func InvalidEntity(v *validator.Validator, subject string) *Error {
	return New(ValidationCode, "validation failed").
		WithFields(v.FieldErrorsGrouped()).
		WithOrigin(fmt.Sprintf("Invalid entity: %s", subject)).
		WithFriendly("Validation failed for the provided input.").
		// opcional: também anexar todas as mensagens planas do validador
		AddMsg(v.Error())
}
