package queries

import (
	"context"
	"log/slog"

	"github.com/andreis3/customers-ms/internal/app/dto"
	"github.com/andreis3/customers-ms/internal/app/mapper"
	"github.com/andreis3/customers-ms/internal/domain/errors"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/postgres"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/service"
	"github.com/andreis3/customers-ms/internal/infra/repositories/criteria"
	"github.com/andreis3/customers-ms/internal/util"
)

type GetCustomerAddresses struct {
	log               adapter.Logger
	addressRepository postgres.AddressesSearchRepository
	customerService   service.CustomerService
	tracer            adapter.Tracer
}

func NewGetCustomerAddresses(
	log adapter.Logger,
	addressRepository postgres.AddressesSearchRepository,
	customerService service.CustomerService,
	tracer adapter.Tracer,
) *GetCustomerAddresses {
	return &GetCustomerAddresses{
		log:               log,
		addressRepository: addressRepository,
		customerService:   customerService,
		tracer:            tracer,
	}
}

func (q *GetCustomerAddresses) Execute(ctx context.Context, input dto.GetCustomerAddressesInput) (*[]dto.GetCustomerAddressesOutput, *errors.Error) {
	ctx, span := q.tracer.Start(ctx, "GetAddress.Execute")
	defer span.End()
	traceID := span.SpanContext().TraceID()
	q.log.InfoJSON("Received input to get address",
		slog.String("trace_id", traceID),
		slog.Any("input", input))

	customerAlreadyExists := q.customerService.ExistCustomerByEmail(ctx, input.Email)

	if !customerAlreadyExists {
		err := errors.ErrCustomerNotFound()
		span.RecordError(err)
		q.log.ErrorJSON("Customer not found",
			slog.String("trace_id", traceID),
			slog.Any("error", err))
		return nil, err
	}

	customerAddresses, err := q.addressRepository.SearchAddresses(ctx, criteria.AddressSearchCriteria{
		CustomerID: util.ToInt64Pointer(input.CustomerID),
	})

	if err != nil {
		span.RecordError(err)
		q.log.ErrorJSON("Failed to get address",
			slog.String("trace_id", traceID),
			slog.Any("error", err))

		return nil, err
	}

	return mapper.ToGetCustomerAddressesOutput(customerAddresses), nil
}
