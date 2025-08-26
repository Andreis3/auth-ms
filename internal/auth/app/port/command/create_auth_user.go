package command

import (
	"context"

	"github.com/andreis3/auth-ms/internal/auth/app/dto"
	"github.com/andreis3/auth-ms/internal/auth/domain/errors"
)

type CreateAuthUser interface {
	Execute(ctx context.Context, input dto.CreateAuthUserInput) (*dto.CreateAuthUserOutput, *errors.Error)
}
