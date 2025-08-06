package valueobject

import (
	"slices"
	"strings"
	"unicode"

	"github.com/andreis3/customers-ms/internal/domain/validator"
)

const (
	CPFLength            = 11
	CPFFirstDigitIdx     = 9
	CPFSecondDigitIdx    = 10
	CPFModuleDivisor     = 11
	CPFBlacklistLength   = 10
	CPFASCIIZero         = '0'
	CPFFirstDigitWeight  = 10
	CPFSecondDigitWeight = 11
)

var blackListCPF = []string{
	"00000000000", "11111111111", "22222222222", "33333333333",
	"44444444444", "55555555555", "66666666666", "77777777777",
	"88888888888", "99999999999",
}

type CPF struct {
	value string
}

func NewCPF(cpf string) CPF {
	return CPF{value: cpf}
}

func (c *CPF) Validate() *validator.Validator {
	var validate validator.Validator
	cleanedCPF := cleanCPF(c.value)
	validate.Assert(validator.NotBlank(cleanedCPF), "cpf", validator.ErrNotBlank)
	validate.Assert(validateCPF(cleanedCPF), "cpf", "cpf invalid")

	return &validate
}

func cleanCPF(cpf string) string {
	var sb strings.Builder
	for _, r := range cpf {
		if unicode.IsDigit(r) {
			sb.WriteRune(r)
		}
	}

	return sb.String()
}

func validateCPF(cpf string) bool {
	if len(cpf) != CPFLength || slices.Contains(blackListCPF, cpf) {
		return false
	}
	return validateDigit(cpf, CPFFirstDigitIdx, CPFFirstDigitWeight) &&
		validateDigit(cpf, CPFSecondDigitIdx, CPFSecondDigitWeight)
}

func validateDigit(cpf string, position, startWeight int) bool {
	sum := 0
	for i, char := range cpf[:position] {
		sum += int(char-CPFASCIIZero) * (startWeight - i)
	}

	rest := (sum * 10) % CPFModuleDivisor
	if rest == CPFBlacklistLength {
		rest = 0
	}

	return rest == int(cpf[position]-CPFASCIIZero)
}

func (c *CPF) String() string {
	return c.value
}
