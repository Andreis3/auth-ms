package mapper

import (
	"github.com/andreis3/customers-ms/internal/app/dto"
)

func TokenOutput(token string, expiresAt int64) *dto.LoginOutput {
	return &dto.LoginOutput{
		Token:     token,
		ExpiresAt: expiresAt,
	}
}
