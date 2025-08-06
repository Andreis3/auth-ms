package uow_factory

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/andreis3/customers-ms/internal/domain/interfaces/adapter"
	iuow "github.com/andreis3/customers-ms/internal/domain/interfaces/uow"
	"github.com/andreis3/customers-ms/internal/infra/uow"
)

func NewUnitOfWorkFactory(pool *pgxpool.Pool, prometheus adapter.Prometheus, tracer adapter.Tracer) func(ctx context.Context) iuow.UnitOfWork {
	return func(ctx context.Context) iuow.UnitOfWork {
		return uow.NewUnitOfWork(pool, prometheus, tracer)
	}
}
