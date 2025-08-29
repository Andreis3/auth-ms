package port

import (
	"context"

	"github.com/andreis3/auth-ms/internal/domain/entity"
	"github.com/andreis3/auth-ms/internal/domain/errors"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user entity.User) (*entity.User, *errors.Error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, *errors.Error)
}
