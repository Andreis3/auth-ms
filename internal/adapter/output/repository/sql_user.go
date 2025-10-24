package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/andreis3/auth-ms/internal/adapter/output/model"
	"github.com/andreis3/auth-ms/internal/domain/entity"
	"github.com/andreis3/auth-ms/internal/domain/errors"
	"github.com/andreis3/auth-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/auth-ms/internal/infra/db"
	"github.com/andreis3/auth-ms/internal/util"
)

type User struct {
	DB      adapter.Postgres
	metrics adapter.Prometheus
	tracer  adapter.Tracer
	model.User
}

func NewUserRepository(db adapter.Postgres, metrics adapter.Prometheus, tracer adapter.Tracer) *User {
	return &User{
		DB:      db,
		metrics: metrics,
		tracer:  tracer,
	}
}

func (u *User) CreateUser(ctx context.Context, user entity.User) (*entity.User, *errors.Error) {
	start := time.Now()
	ctx, span := u.tracer.Start(ctx, "UserRepository.CreateUser")

	defer func() {
		end := time.Since(start)
		u.metrics.ObserveInstructionDBDuration("postgres", "users", "insert", float64(end.Milliseconds()))
		span.End()
	}()

	modelUser := u.ToModel(user)

	const query = `
	INSERT INTO users (public_id, email, password_hash, name, role, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7) 
	RETURNING id`

	var id int64

	db := u.resolveDB(ctx)

	err := db.QueryRow(ctx, query,
		modelUser.PublicID,
		modelUser.Email,
		modelUser.Password,
		modelUser.Name,
		modelUser.Role,
		modelUser.CreatedAt,
		modelUser.UpdatedAt).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, errors.ErrorAlreadyExistsUser(err)
		}
		return nil, errors.CreateUserError(err)
	}

	user.AssignID(id)
	user.AssignCreateAT(util.ToTime(modelUser.CreatedAt))
	user.AssignUpdateAT(util.ToTime(modelUser.UpdatedAt))
	return &user, nil
}

func (u *User) FindUserByEmail(ctx context.Context, email string) (*entity.User, *errors.Error) {
	ctx, span := u.tracer.Start(ctx, "CustomerRepository.FindCustomerByEmail")
	start := time.Now()

	defer func() {
		end := time.Since(start)
		u.metrics.ObserveInstructionDBDuration("postgres", "customers", "select", float64(end.Milliseconds()))
		span.End()
	}()

	const query = `
	SELECT id, public_id, email, password_hash, name, role, created_at, updated_at, deleted_at
	FROM users
	WHERE email = $1`

	var model model.User
	db := u.resolveDB(ctx)

	rows, err := db.Query(ctx, query, email)
	if err != nil {
		return nil, errors.ErrorFindUserByEmail(err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}
	err = rows.Scan(
		&model.ID,
		&model.PublicID,
		&model.Email,
		&model.Password,
		&model.Name,
		&model.Role,
		&model.CreatedAt,
		&model.UpdatedAt,
		&model.DeletedAt,
	)
	if err != nil {
		return nil, errors.ErrorFindUserByEmail(err)
	}

	result := model.ToEntity()

	// TODO: create SetAttributes in interface otel
	//span.SetAttributes(
	//	attribute.Int64("customer_id", *modelCustomer.CustomerID),
	//)

	return &result, nil
}

func (u *User) resolveDB(ctx context.Context) adapter.Postgres {
	if tx, ok := db.TxFromContext(ctx); ok {
		return tx
	}
	return u.DB
}
