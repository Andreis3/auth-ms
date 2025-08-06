package postgres

import (
	"context"

	"github.com/andreis3/customers-ms/internal/domain/entity"
	"github.com/andreis3/customers-ms/internal/domain/errors"
	"github.com/andreis3/customers-ms/internal/infra/repositories/criteria"
)

type AddressRepository interface {
	InsertBatchAddress(ctx context.Context, customerID int64, addresses []entity.Address) (*[]entity.Address, *errors.Error)
	FindAddressesByCustomerID(ctx context.Context, customerID int64) (*[]entity.Address, *errors.Error)
	SearchAddresses(ctx context.Context, params criteria.AddressSearchCriteria) (*[]entity.Address, *errors.Error)
}

type AddressesSearchRepository interface {
	SearchAddresses(ctx context.Context, params criteria.AddressSearchCriteria) (*[]entity.Address, *errors.Error)
}
