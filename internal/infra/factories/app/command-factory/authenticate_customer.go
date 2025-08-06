package command_factory

import (
	"github.com/andreis3/customers-ms/internal/app/commands"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/command"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/postgres"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/service"
)

type authenticateCustomerFactory struct {
	log                adapter.Logger
	customerRepository postgres.CustomerRepository
	authService        service.Auth
	bcrypt             adapter.Bcrypt
	tracer             adapter.Tracer
}

func NewAuthenticateCustomerFactory(
	log adapter.Logger,
	customerRepository postgres.CustomerRepository,
	authService service.Auth,
	bcrypt adapter.Bcrypt,
	tracer adapter.Tracer,
) command.Login {
	return commands.NewLoginCustomer(log, customerRepository, authService, bcrypt, tracer)
}
