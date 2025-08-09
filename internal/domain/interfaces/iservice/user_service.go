package iservice

import (
	"context"

	"github.com/andreis3/auth-ms/internal/domain/errors"
)

type IUserService interface {
	ValidateEmailAvailability(ctx context.Context, email string) *errors.Error
}
