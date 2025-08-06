package commands

import (
	"context"
	"log/slog"

	"github.com/andreis3/customers-ms/internal/app/dto"
	"github.com/andreis3/customers-ms/internal/app/mapper"
	"github.com/andreis3/customers-ms/internal/domain/entity"
	"github.com/andreis3/customers-ms/internal/domain/errors"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/postgres"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/service"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/uow"
)

type CreateCustomerCommand struct {
	uow                uow.UnitOfWorkFactory
	bcrypt             adapter.Bcrypt
	log                adapter.Logger
	customerRepository postgres.CustomerRepository
	addressRepository  postgres.AddressRepository
	customerService    service.CustomerService
	tracer             adapter.Tracer
}

func NewCreateCustomer(
	uow uow.UnitOfWorkFactory,
	bcrypt adapter.Bcrypt,
	log adapter.Logger,
	customerService service.CustomerService,
	customerRepository postgres.CustomerRepository,
	addressRepository postgres.AddressRepository,
	tracer adapter.Tracer,
) *CreateCustomerCommand {
	return &CreateCustomerCommand{
		uow:                uow,
		bcrypt:             bcrypt,
		log:                log,
		customerService:    customerService,
		customerRepository: customerRepository,
		addressRepository:  addressRepository,
		tracer:             tracer,
	}
}

func (c *CreateCustomerCommand) Execute(ctx context.Context, input dto.CreateCustomerInput) (*entity.Customer, *errors.Error) {
	ctx, span := c.tracer.Start(ctx, "CreatedCustomer.Execute")
	customerProfile := mapper.ToCustomerProfile(input)
	defer span.End()
	traceID := span.SpanContext().TraceID()
	c.log.InfoJSON("Received input to create CreateCustomer",
		slog.String("trace_id", traceID),
		slog.Any("input", input))

	err := customerProfile.Validate()
	if err != nil {
		span.RecordError(err)
		c.log.ErrorJSON("Failed validate input to create CreateCustomer",
			slog.String("trace_id", traceID),
			slog.Any("error", err.Errors))
		return nil, err
	}

	customerAlreadyExists := c.customerService.ExistCustomerByEmail(ctx, customerProfile.Customer.Email.String())
	if customerAlreadyExists {
		err = errors.ErrCustomerAlreadyExists()
		span.RecordError(err)
		c.log.ErrorJSON("Customer already exists",
			slog.String("trace_id", traceID),
			slog.Any("error", err))
		return nil, err
	}

	var customerResult *entity.Customer
	uow := c.uow(ctx)

	err = uow.WithTransaction(ctx, func(ctxUow context.Context) *errors.Error {
		hash, errInternal := c.bcrypt.Hash(customerProfile.Customer.Password.String())
		if errInternal != nil {
			span.RecordError(err)
			c.log.ErrorJSON("Failed hash password",
				slog.String("trace_id", traceID),
				slog.Any("error", err))
			return err
		}

		customerProfile.Customer.AssignHashedPassword(hash)

		customerResult, errInternal = c.customerRepository.InsertCustomer(ctxUow, customerProfile.Customer)
		if errInternal != nil {
			span.RecordError(err)
			c.log.ErrorJSON("Failed insert customerResult",
				slog.String("trace_id", traceID),
				slog.Any("error", err))
			return err
		}

		if len(customerProfile.Addresses) > 0 {
			_, errInternal = c.addressRepository.InsertBatchAddress(ctxUow, customerResult.ID, customerProfile.Addresses)
			if errInternal != nil {
				span.RecordError(err)
				c.log.ErrorJSON("Failed insert address",
					slog.String("trace_id", traceID),
					slog.Any("error", err))
				return err
			}
		}
		return nil
	})
	if err != nil {
		span.RecordError(err)
		c.log.ErrorJSON("Failed insert customerResult",
			slog.String("trace_id", traceID),
			slog.Any("error", err))
		return nil, err
	}

	return customerResult, nil
}
