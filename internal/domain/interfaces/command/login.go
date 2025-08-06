package command

import (
	"context"

	"github.com/andreis3/customers-ms/internal/app/dto"
	"github.com/andreis3/customers-ms/internal/domain/errors"
)

type Login interface {
	Execute(ctx context.Context, input dto.LoginInput) (*dto.LoginOutput, *errors.Error)
}
