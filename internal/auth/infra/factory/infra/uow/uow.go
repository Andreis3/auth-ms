package uow

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/andreis3/auth-ms/internal/auth/domain/interfaces/adapter"
	"github.com/andreis3/auth-ms/internal/auth/infra/uow"
)

func NewUnitOfWorkFactory(pool *pgxpool.Pool, prometheus adapter.Prometheus, tracer adapter.Tracer) func(ctx context.Context) adapter.UnitOfWork {
	return func(ctx context.Context) adapter.UnitOfWork {
		return uow.NewUnitOfWork(pool, prometheus, tracer)
	}
}
