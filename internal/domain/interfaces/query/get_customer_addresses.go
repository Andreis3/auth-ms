package query

import (
	"context"

	"github.com/andreis3/customers-ms/internal/app/dto"
	"github.com/andreis3/customers-ms/internal/domain/errors"
)

type GetCustomerAddresses interface {
	Execute(ctx context.Context, input dto.GetCustomerAddressesInput) (*[]dto.GetCustomerAddressesOutput, *errors.Error)
}
