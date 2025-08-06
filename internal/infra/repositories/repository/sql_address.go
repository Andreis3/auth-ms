package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/andreis3/customers-ms/internal/domain/entity"
	"github.com/andreis3/customers-ms/internal/domain/errors"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/customers-ms/internal/infra/adapters/db"
	"github.com/andreis3/customers-ms/internal/infra/repositories/criteria"
	"github.com/andreis3/customers-ms/internal/infra/repositories/model"
)

type AddressRepository struct {
	DB      adapter.Postgres
	metrics adapter.Prometheus
	model.Address
	tracer adapter.Tracer
}

func NewAddressRepository(
	db adapter.Postgres,
	metrics adapter.Prometheus,
	tracer adapter.Tracer,
) *AddressRepository {
	return &AddressRepository{
		DB:      db,
		metrics: metrics,
		tracer:  tracer,
	}
}

func (c *AddressRepository) InsertBatchAddress(ctx context.Context, customerID int64, addresses []entity.Address) (*[]entity.Address, *errors.Error) {
	ctx, span := c.tracer.Start(ctx, "AddressRepository.InsertBatchAddress")
	start := time.Now()

	defer func() {
		end := time.Since(start)
		c.metrics.ObserveInstructionDBDuration("postgres", "addresses", "insert", float64(end.Milliseconds()))
		span.End()
	}()

	batch := &pgx.Batch{}

	query := `
	INSERT INTO addresses 
	(customer_id, street, number, complement, city, state, postal_code, country, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
	RETURNING id`

	for _, address := range addresses {
		modelAddress := c.FromModel(address)
		batch.Queue(query,
			customerID,
			modelAddress.Street,
			modelAddress.Number,
			modelAddress.Complement,
			modelAddress.City,
			modelAddress.State,
			modelAddress.PostalCode,
			modelAddress.Country,
			modelAddress.CreatedAt,
			modelAddress.UpdatedAt,
		)
	}

	db := c.resolveDB(ctx)

	br := db.SendBatch(ctx, batch)
	defer br.Close()

	addressesResult := make([]entity.Address, 0, len(addresses))

	for _, address := range addresses {

		var id int64

		_ = br.QueryRow().Scan(&id)

		addressCopy := address
		addressCopy.AssignID(id)

		addressesResult = append(addressesResult, addressCopy)
	}

	return &addressesResult, nil
}

func (c *AddressRepository) FindAddressesByCustomerID(ctx context.Context, customerID int64) (*[]entity.Address, *errors.Error) {
	ctx, span := c.tracer.Start(ctx, "AddressRepository.FindByCustomerID")
	start := time.Now()

	defer func() {
		end := time.Since(start)
		c.metrics.ObserveInstructionDBDuration("postgres", "addresses", "select", float64(end.Milliseconds()))
		span.End()
	}()

	const query = `
	SELECT id, customer_id, street, number, complement, city, state, postal_code, country, created_at, updated_at
	FROM addresses
	WHERE customer_id = $1`

	db := c.resolveDB(ctx)

	rows, err := db.Query(ctx, query, customerID)
	if err != nil {
		return nil, errors.ErrorFindByCustomerID(err)
	}
	defer rows.Close()

	addresses := make([]entity.Address, 0)
	for rows.Next() {
		var modelAddress model.Address
		err = rows.Scan(
			&modelAddress.ID,
			&modelAddress.CustomerID,
			&modelAddress.Street,
			&modelAddress.Number,
			&modelAddress.Complement,
			&modelAddress.City,
			&modelAddress.State,
			&modelAddress.PostalCode,
			&modelAddress.Country,
			&modelAddress.CreatedAt,
			&modelAddress.UpdatedAt,
		)
		if err != nil {
			return nil, errors.ErrorFindByCustomerID(err)
		}
		addresses = append(addresses, modelAddress.ToEntity())
	}
	return &addresses, nil
}

func (c *AddressRepository) SearchAddresses(ctx context.Context, params criteria.AddressSearchCriteria) (*[]entity.Address, *errors.Error) {
	var conditions []string
	var args []interface{}
	i := 1

	if params.CustomerID != nil {
		conditions = append(conditions, fmt.Sprintf("customer_id = $%d", i))
		args = append(args, *params.CustomerID)
		i++
	}

	if params.Email != nil {
		conditions = append(conditions, fmt.Sprintf("email = $%d", i))
		args = append(args, *params.Email)
		i++
	}

	query := `
	SELECT id, street, number, complement, city, state, postal_code, country, created_at, updated_at
	FROM addresses`

	if len(conditions) == 0 {
		return nil, nil
	}

	query += fmt.Sprintf(" WHERE %s", strings.Join(conditions, " AND "))

	db := c.resolveDB(ctx)

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.ErrorFindByCustomerID(err)
	}
	defer rows.Close()

	addresses := make([]entity.Address, 0)
	for rows.Next() {
		var modelAddress model.Address
		err = rows.Scan(
			&modelAddress.ID,
			&modelAddress.Street,
			&modelAddress.Number,
			&modelAddress.Complement,
			&modelAddress.City,
			&modelAddress.State,
			&modelAddress.PostalCode,
			&modelAddress.Country,
			&modelAddress.CreatedAt,
			&modelAddress.UpdatedAt,
		)
		if err != nil {
			return nil, errors.ErrorFindByCustomerID(err)
		}
		addresses = append(addresses, modelAddress.ToEntity())
	}

	return &addresses, nil
}

func (c *AddressRepository) resolveDB(ctx context.Context) adapter.Postgres {
	if tx, ok := db.TxFromContext(ctx); ok {
		return tx
	}
	return c.DB
}
