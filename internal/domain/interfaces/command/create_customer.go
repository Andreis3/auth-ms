package command

import (
	"context"

	"github.com/andreis3/customers-ms/internal/app/dto"
	"github.com/andreis3/customers-ms/internal/domain/entity"
	"github.com/andreis3/customers-ms/internal/domain/errors"
)

type CreateCustomer interface {
	Execute(ctx context.Context, input dto.CreateCustomerInput) (*entity.Customer, *errors.Error)
}
