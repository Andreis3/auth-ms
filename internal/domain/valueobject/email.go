package valueobject

import (
	"regexp"
	"strings"

	"github.com/andreis3/customers-ms/internal/domain/validator"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type Email struct {
	value string
}

func NewEmail(email string) Email {
	return Email{value: email}
}

func (e *Email) Validate() *validator.Validator {
	var validate validator.Validator

	cleanedEmail := cleanEmail(e.value)
	validate.Assert(validator.NotBlank(cleanedEmail), "email", validator.ErrNotBlank)
	validate.Assert(isValidEmail(cleanedEmail), "email", "invalid email format")

	return &validate
}

func cleanEmail(email string) string {
	return strings.TrimSpace(email)
}

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func (e *Email) String() string {
	return e.value
}
