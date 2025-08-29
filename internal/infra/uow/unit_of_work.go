package uow

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	errors2 "github.com/andreis3/auth-ms/internal/domain/errors"
	adapter2 "github.com/andreis3/auth-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/auth-ms/internal/infra/db"
)

type UnitOfWork struct {
	DB         *pgxpool.Pool
	TX         pgx.Tx
	prometheus adapter2.Prometheus
	tracer     adapter2.Tracer
}

func NewUnitOfWork(db *pgxpool.Pool, prometheus adapter2.Prometheus, tracer adapter2.Tracer) *UnitOfWork {
	return &UnitOfWork{
		DB:         db,
		prometheus: prometheus,
		tracer:     tracer,
	}
}

// WithTransaction handles transaction lifecycle safely.
func (u *UnitOfWork) WithTransaction(ctx context.Context, fn func(ctx context.Context) *errors2.Error) *errors2.Error {
	ctx, span := u.tracer.Start(ctx, "UnitOfWork.WithTransaction")
	defer func() {
		span.End()
		u.TX = nil
	}()

	if u.TX != nil {
		span.RecordError(errors2.ErrorTransactionAlreadyExists())
		return errors2.ErrorTransactionAlreadyExists()
	}

	tx, err := u.DB.Begin(ctx)
	if err != nil {
		span.RecordError(errors2.ErrorOpeningTransaction(err))
		return errors2.ErrorOpeningTransaction(err)
	}

	u.TX = tx
	ctxTx := db.WithTx(ctx, tx)

	if err := fn(ctxTx); err != nil {
		rollbackErr := u.TX.Rollback(ctx)
		if rollbackErr != nil {
			span.RecordError(errors2.ErrorExecuteRollback(rollbackErr))
			return errors2.ErrorExecuteRollback(rollbackErr)
		}
		span.RecordError(err)
		return err
	}

	if err := u.TX.Commit(ctx); err != nil {
		span.RecordError(errors2.ErrorCommitOrRollback(err))
		return errors2.ErrorCommitOrRollback(err)
	}

	return nil
}
