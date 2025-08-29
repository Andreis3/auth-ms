package service

import (
	"context"

	"github.com/andreis3/auth-ms/internal/domain/errors"
)

type UserService interface {
	ValidateEmailAvailability(ctx context.Context, email string) *errors.Error
}
