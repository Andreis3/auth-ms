package mservice

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/andreis3/auth-ms/internal/domain/errors"
)

type UserServiceMock struct{ mock.Mock }

func (s *UserServiceMock) ValidateEmailAvailability(ctx context.Context, email string) *errors.Error {
	args := s.Called(ctx, email)

	if v := args.Get(0); v != nil {
		return v.(*errors.Error)
	}

	return nil
}
