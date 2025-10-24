package uow

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/andreis3/auth-ms/internal/domain/errors"
	"github.com/andreis3/auth-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/auth-ms/internal/infra/db"
)

type UnitOfWork struct {
	DB         *pgxpool.Pool
	prometheus adapter.Prometheus
	tracer     adapter.Tracer
}

func NewUnitOfWork(db *pgxpool.Pool, prometheus adapter.Prometheus, tracer adapter.Tracer) *UnitOfWork {
	return &UnitOfWork{
		DB:         db,
		prometheus: prometheus,
		tracer:     tracer,
	}
}

// WithTransaction handles transaction lifecycle safely.
func (u *UnitOfWork) WithTransaction(ctx context.Context, fn func(ctx context.Context) *errors.Error) *errors.Error {
	ctx, span := u.tracer.Start(ctx, "UnitOfWork.WithTransaction")
	defer func() {
		span.End()
	}()

	start := time.Now()
	status := "success"
	defer func() {
		u.prometheus.ObserveInstructionDBDuration("postgres", "transaction", "with_transaction_"+status, float64(time.Since(start).Milliseconds()))
	}()

	if _, ok := db.TxFromContext(ctx); ok {
		status = "error"
		span.RecordError(errors.ErrorTransactionAlreadyExists())
		return errors.ErrorTransactionAlreadyExists()
	}

	tx, err := u.DB.Begin(ctx)
	if err != nil {
		status = "error"
		span.RecordError(errors.ErrorOpeningTransaction(err))
		return errors.ErrorOpeningTransaction(err)
	}

	ctxTx := db.WithTx(ctx, tx)

	if err := fn(ctxTx); err != nil {
		status = "error"
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			joinedErr := errors.Join(err, rollbackErr)
			rollbackError := errors.ErrorExecuteRollback(joinedErr)
			span.RecordError(rollbackError)
			return rollbackError
		}
		span.RecordError(err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		status = "error"
		span.RecordError(errors.ErrorCommitOrRollback(err))
		return errors.ErrorCommitOrRollback(err)
	}

	return nil
}
