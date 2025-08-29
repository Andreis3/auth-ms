package uow

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	adapter2 "github.com/andreis3/auth-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/auth-ms/internal/infra/uow"
)

func NewUnitOfWorkFactory(pool *pgxpool.Pool, prometheus adapter2.Prometheus, tracer adapter2.Tracer) func(ctx context.Context) adapter2.UnitOfWork {
	return func(ctx context.Context) adapter2.UnitOfWork {
		return uow.NewUnitOfWork(pool, prometheus, tracer)
	}
}
