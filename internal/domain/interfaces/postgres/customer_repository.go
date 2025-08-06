package postgres

import (
	"context"

	"github.com/andreis3/customers-ms/internal/domain/entity"
	"github.com/andreis3/customers-ms/internal/domain/errors"
)

type CustomerRepository interface {
	InsertCustomer(ctx context.Context, data entity.Customer) (*entity.Customer, *errors.Error)
	FindCustomerByEmail(ctx context.Context, email string) (*entity.Customer, *errors.Error)

	FindByID(ctx context.Context, id int64) (*entity.Customer, *errors.Error)
}
