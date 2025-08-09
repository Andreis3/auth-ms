package iuow

import (
	"context"

	"github.com/andreis3/auth-ms/internal/domain/errors"
)

type RepositoryFactory func(tx any) any

type IUnitOfWork interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) *errors.Error) *errors.Error
}
