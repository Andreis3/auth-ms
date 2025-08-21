package adapter

import (
	"context"

	"github.com/andreis3/auth-ms/internal/auth/domain/errors"
)

type UnitOfWorkFactory func(ctx context.Context) UnitOfWork

type RepositoryFactory func(tx any) any

type UnitOfWork interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) *errors.Error) *errors.Error
}
