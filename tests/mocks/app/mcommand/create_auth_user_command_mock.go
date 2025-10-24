package mcommand

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/andreis3/auth-ms/internal/app/dto"
	"github.com/andreis3/auth-ms/internal/domain/errors"
)

type CreateAuthUserCommandMock struct{ mock.Mock }

func (c *CreateAuthUserCommandMock) Execute(ctx context.Context, input dto.CreateAuthUserInput) (*dto.CreateAuthUserOutput, *errors.Error) {
	args := c.Called(ctx, input)

	var output *dto.CreateAuthUserOutput
	if v := args.Get(0); v != nil {
		output = v.(*dto.CreateAuthUserOutput)
	}

	if err := args.Get(1); err != nil {
		return output, err.(*errors.Error)
	}

	return output, nil
}
