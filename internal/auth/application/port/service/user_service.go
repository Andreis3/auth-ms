package service

import (
	"context"

	"github.com/andreis3/auth-ms/internal/auth/domain/errors"
)

type UserService interface {
	ValidateEmailAvailability(ctx context.Context, email string) *errors.Error
}
