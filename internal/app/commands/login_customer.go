package commands

import (
	"context"
	"log/slog"

	"github.com/andreis3/customers-ms/internal/app/dto"
	"github.com/andreis3/customers-ms/internal/app/mapper"
	"github.com/andreis3/customers-ms/internal/domain/errors"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/postgres"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/service"
)

type Login struct {
	log                adapter.Logger
	customerRepository postgres.CustomerRepository
	authService        service.Auth
	bcrypt             adapter.Bcrypt
	tracer             adapter.Tracer
}

func NewLoginCustomer(
	log adapter.Logger,
	customerRepository postgres.CustomerRepository,
	authService service.Auth,
	bcrypt adapter.Bcrypt,
	tracer adapter.Tracer,
) *Login {
	return &Login{
		log:                log,
		customerRepository: customerRepository,
		authService:        authService,
		bcrypt:             bcrypt,
		tracer:             tracer,
	}
}

func (a *Login) Execute(ctx context.Context, input dto.LoginInput) (*dto.LoginOutput, *errors.Error) {
	ctx, span := a.tracer.Start(ctx, "Login.Execute")
	a.log.InfoJSON("Received input to authenticate customer",
		slog.String("email", input.Email),
		slog.String("password", "**************"))

	defer span.End()
	traceID := span.SpanContext().TraceID()

	customer, err := a.customerRepository.FindCustomerByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if customer == nil {
		return nil, errors.ErrorInvalidCredentials()
	}

	isValid := a.bcrypt.CompareHash(input.Password, customer.Password.String())
	if !isValid {
		return nil, errors.ErrorInvalidCredentials()
	}

	token, err := a.authService.GenerateToken(*customer)

	if err != nil {
		return nil, err
	}

	a.log.InfoJSON("end request",
		slog.String("trace_id", traceID),
		slog.String("token", token.Token))

	return mapper.TokenOutput(token.Token, token.ExpiresAt.Unix()), nil
}
