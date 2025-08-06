package handler_factory

import (
	"github.com/andreis3/customers-ms/internal/app/commands"
	"github.com/andreis3/customers-ms/internal/app/services"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/command"
	"github.com/andreis3/customers-ms/internal/infra/adapters/crypto"
	"github.com/andreis3/customers-ms/internal/infra/adapters/db"
	"github.com/andreis3/customers-ms/internal/infra/configs"
	"github.com/andreis3/customers-ms/internal/infra/factories/infra/uow-factory"
	"github.com/andreis3/customers-ms/internal/infra/repositories/repository"
	"github.com/andreis3/customers-ms/internal/presentation/http/handler"
)

type CreateCustomer struct {
	db      *db.Postgres
	redis   *db.Redis
	log     adapter.Logger
	metrics adapter.Prometheus
	tracer  adapter.Tracer
	conf    *configs.Configs
}

func NewCreateCustomer(db *db.Postgres, redis *db.Redis, log adapter.Logger, metrics adapter.Prometheus, tracer adapter.Tracer, conf *configs.Configs) *CreateCustomer {
	return &CreateCustomer{db, redis, log, metrics, tracer, conf}
}

func (f *CreateCustomer) NewCreateCustomer() *handler.CreateCustomerHandler {
	crypto := crypto.NewBcrypt()
	command := newCreateCustomer(f.db, crypto, f.log, f.tracer, f.metrics)
	return handler.NewCreateCustomerHandler(command, f.metrics, f.log, f.tracer)
}

func newCreateCustomer(
	db *db.Postgres,
	crypto adapter.Bcrypt,
	log adapter.Logger,
	tracer adapter.Tracer,
	metrics adapter.Prometheus,
) command.CreateCustomer {
	customerRepository := repository.NewCustomerRepository(db, metrics, tracer)
	addressRepository := repository.NewAddressRepository(db, metrics, tracer)
	uowFactory := uow_factory.NewUnitOfWorkFactory(db.Pool, metrics, tracer)
	customerService := services.NewCustomerService(customerRepository)
	return commands.NewCreateCustomer(uowFactory, crypto, log, customerService, customerRepository, addressRepository, tracer)
}
