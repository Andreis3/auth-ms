package uow_factory

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/andreis3/auth-ms/internal/domain/interfaces/iadapter"
	iuow "github.com/andreis3/auth-ms/internal/domain/interfaces/iuow"
	"github.com/andreis3/auth-ms/internal/infra/uow"
)

func NewUnitOfWorkFactory(pool *pgxpool.Pool, prometheus iadapter.IPrometheus, tracer iadapter.ITracer) func(ctx context.Context) iuow.IUnitOfWork {
	return func(ctx context.Context) iuow.IUnitOfWork {
		return uow.NewUnitOfWork(pool, prometheus, tracer)
	}
}
