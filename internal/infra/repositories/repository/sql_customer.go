package repository

import (
	"context"
	"time"

	"github.com/andreis3/customers-ms/internal/domain/entity"
	"github.com/andreis3/customers-ms/internal/domain/errors"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/customers-ms/internal/infra/adapters/db"
	"github.com/andreis3/customers-ms/internal/infra/repositories/model"
)

type CustomerRepository struct {
	DB      adapter.Postgres
	metrics adapter.Prometheus
	model.Customer
	tracer adapter.Tracer
}

func NewCustomerRepository(
	db adapter.Postgres,
	metrics adapter.Prometheus,
	tracer adapter.Tracer,
) *CustomerRepository {
	return &CustomerRepository{
		DB:      db,
		metrics: metrics,
		tracer:  tracer,
	}
}

func (c *CustomerRepository) InsertCustomer(ctx context.Context, data entity.Customer) (*entity.Customer, *errors.Error) {
	ctx, span := c.tracer.Start(ctx, "CustomerRepository.InsertCustomer")
	start := time.Now()
	defer func() {
		end := time.Since(start)
		c.metrics.ObserveInstructionDBDuration("postgres", "customers", "insert", float64(end.Milliseconds()))
		span.End()
	}()
	modelCustomer := c.FromModel(data)
	const query = `
	INSERT INTO customers 
	(email, password, first_name, last_name, cpf, date_of_birth, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
	RETURNING id`

	var id int64

	db := c.resolveDB(ctx)

	err := db.QueryRow(ctx, query,
		modelCustomer.Email,
		modelCustomer.Password,
		modelCustomer.FirstName,
		modelCustomer.LastName,
		modelCustomer.CPF,
		modelCustomer.DateOfBirth,
		modelCustomer.CreatedAT,
		modelCustomer.UpdatedAT).Scan(&id)
	if err != nil {
		return nil, errors.ErrorSaveCustomer(err)
	}

	modelCustomer.ID = &id
	result := modelCustomer.ToEntity()

	// TODO: create SetAttributes in interface otel
	//span.SetAttributes(
	//	attribute.Int64("customer_id", id),
	//)

	return &result, nil
}

func (c *CustomerRepository) FindCustomerByEmail(ctx context.Context, email string) (*entity.Customer, *errors.Error) {
	ctx, span := c.tracer.Start(ctx, "CustomerRepository.FindCustomerByEmail")
	start := time.Now()

	defer func() {
		end := time.Since(start)
		c.metrics.ObserveInstructionDBDuration("postgres", "customers", "select", float64(end.Milliseconds()))
		span.End()
	}()

	const query = `
	SELECT id, email, password, first_name, last_name, cpf, date_of_birth, created_at, updated_at
	FROM customers
	WHERE email = $1`

	var modelCustomer model.Customer
	db := c.resolveDB(ctx)

	rows, err := db.Query(ctx, query, email)
	if err != nil {
		return nil, errors.ErrorFindCustomerByEmail(err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}
	err = rows.Scan(
		&modelCustomer.ID,
		&modelCustomer.Email,
		&modelCustomer.Password,
		&modelCustomer.FirstName,
		&modelCustomer.LastName,
		&modelCustomer.CPF,
		&modelCustomer.DateOfBirth,
		&modelCustomer.CreatedAT,
		&modelCustomer.UpdatedAT,
	)
	if err != nil {
		return nil, errors.ErrorFindCustomerByEmail(err)
	}

	result := modelCustomer.ToEntity()

	// TODO: create SetAttributes in interface otel
	//span.SetAttributes(
	//	attribute.Int64("customer_id", *modelCustomer.CustomerID),
	//)

	return &result, nil
}

func (c *CustomerRepository) FindByID(ctx context.Context, id int64) (*entity.Customer, *errors.Error) {
	ctx, span := c.tracer.Start(ctx, "CustomerRepository.FindByID")
	start := time.Now()

	defer func() {
		end := time.Since(start)
		c.metrics.ObserveInstructionDBDuration("postgres", "customers", "select", float64(end.Milliseconds()))
		span.End()
	}()

	var modelCustomer model.Customer
	db := c.resolveDB(ctx)

	const query = `
	SELECT id, email, password, first_name, last_name, cpf, date_of_birth, created_at, updated_at
	FROM customers
	WHERE id = $1`

	rows, err := db.Query(ctx, query, id)
	if err != nil {
		return nil, errors.ErrorFindByID(err)
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}
	err = rows.Scan(
		&modelCustomer.ID,
		&modelCustomer.Email,
		&modelCustomer.Password,
		&modelCustomer.FirstName,
		&modelCustomer.LastName,
		&modelCustomer.CPF,
		&modelCustomer.DateOfBirth,
		&modelCustomer.CreatedAT,
		&modelCustomer.UpdatedAT,
	)
	if err != nil {
		return nil, errors.ErrorFindByID(err)
	}

	result := modelCustomer.ToEntity()

	// TODO: create SetAttributes in interface otel
	//span.SetAttributes(
	//	attribute.Int64("customer_id", *modelCustomer.CustomerID),
	//)

	return &result, nil
}

func (c *CustomerRepository) resolveDB(ctx context.Context) adapter.Postgres {
	if tx, ok := db.TxFromContext(ctx); ok {
		return tx
	}
	return c.DB
}
