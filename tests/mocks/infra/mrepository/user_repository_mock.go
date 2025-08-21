package mrepository

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/andreis3/auth-ms/internal/auth/domain/entity"
	"github.com/andreis3/auth-ms/internal/auth/domain/errors"
)

type UserRepositoryMock struct{ mock.Mock }

func (r *UserRepositoryMock) CreateUser(ctx context.Context, user entity.User) (*entity.User, *errors.Error) {
	args := r.Called(ctx, user)

	var u *entity.User
	if v := args.Get(0); v != nil {
		u = v.(*entity.User)
	}

	var e *errors.Error
	if v := args.Get(1); v != nil {
		e = v.(*errors.Error)
	}

	return u, e
}

func (r *UserRepositoryMock) FindUserByEmail(ctx context.Context, email string) (*entity.User, *errors.Error) {
	args := r.Called(ctx, email)

	var u *entity.User
	if v := args.Get(0); v != nil {
		u = v.(*entity.User)
	}

	var e *errors.Error
	if v := args.Get(1); v != nil {
		e = v.(*errors.Error)
	}

	return u, e
}
