package errors

import (
	"fmt"

	"github.com/andreis3/auth-ms/internal/domain/validator"
)

func InvalidEntity(validate *validator.Validator, messages string) *Error {
	return &Error{
		Code:            BadRequestCode,
		Map:             validate.FieldErrorsGrouped(),
		Errors:          validate.Errors(),
		OriginFunc:      fmt.Sprintf("Invalid entity %s", messages),
		FriendlyMessage: "Validation failed for the provided input.",
	}
}
