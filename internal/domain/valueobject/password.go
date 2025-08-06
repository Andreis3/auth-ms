package valueobject

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/andreis3/customers-ms/internal/domain/validator"
)

type Password struct {
	value     string
	encrypted bool
}

func NewPassword(raw string) Password {
	return Password{
		value:     strings.TrimSpace(raw),
		encrypted: false,
	}
}

func NewEncryptedPassword(encrypted string) Password {
	return Password{
		value:     strings.TrimSpace(encrypted),
		encrypted: true,
	}
}

func (p *Password) Validate() *validator.Validator {
	var validate validator.Validator
	trimmed := strings.TrimSpace(p.value)

	validate.Assert(validator.NotBlank(trimmed), "password", validator.ErrNotBlank)
	validate.Assert(validator.MinChars(trimmed, 8), "password", "must be at least 8 characters")

	validate.Assert(hasUppercase(trimmed), "password", "must contain at least one uppercase letter")
	validate.Assert(hasLowercase(trimmed), "password", "must contain at least one lowercase letter")
	validate.Assert(hasNumber(trimmed), "password", "must contain at least one number")
	validate.Assert(hasSpecialChar(trimmed), "password", "must contain at least one special character")
	validate.Assert(!hasSequentialNumbers(trimmed), "password", "must not contain sequential numbers (e.g. 123)")
	validate.Assert(!hasSequentialLetters(trimmed), "password", "must not contain sequential letters (e.g. abc)")

	return &validate
}

func hasUppercase(value string) bool {
	for _, char := range value {
		if unicode.IsUpper(char) {
			return true
		}
	}

	return false
}

func hasLowercase(value string) bool {
	for _, char := range value {
		if unicode.IsLower(char) {
			return true
		}
	}

	return false
}

func hasNumber(value string) bool {
	for _, char := range value {
		if unicode.IsDigit(char) {
			return true
		}
	}
	return false
}

func hasSpecialChar(value string) bool {
	re := regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]{};':"\\|,.<>\/?]`)
	return re.MatchString(value)
}

func hasSequentialLetters(value string) bool {
	parserLowercase := strings.ToLower(value)
	sequences := []string{
		"abc", "bcd", "cde",
		"def", "efg", "fgh",
		"ghi", "hij", "ijk",
		"jkl", "klm", "lmn",
		"mno", "nop", "opq",
		"pqr", "qrs", "rst",
		"stu", "tuv", "uvw",
		"vwx", "wxy", "xyz",
	}

	for _, seq := range sequences {
		if strings.Contains(parserLowercase, seq) {
			return true
		}
	}

	return false
}

func hasSequentialNumbers(value string) bool {
	sequences := []string{"012", "123", "234", "345", "456", "567", "678", "789"}

	for _, seq := range sequences {
		if strings.Contains(value, seq) {
			return true
		}
	}

	return false
}

func (p *Password) String() string {
	return p.value
}
