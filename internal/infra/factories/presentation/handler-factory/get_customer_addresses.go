package handler_factory

import (
	"github.com/andreis3/customers-ms/internal/app/queries"
	"github.com/andreis3/customers-ms/internal/app/services"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/query"
	"github.com/andreis3/customers-ms/internal/infra/adapters/db"
	"github.com/andreis3/customers-ms/internal/infra/configs"
	"github.com/andreis3/customers-ms/internal/infra/decorator"
	"github.com/andreis3/customers-ms/internal/infra/repositories/repository"
	"github.com/andreis3/customers-ms/internal/presentation/http/handler"
)

type GetCustomerAddresses struct {
	db      *db.Postgres
	cache   *db.Redis
	log     adapter.Logger
	tracer  adapter.Tracer
	metrics adapter.Prometheus
	conf    *configs.Configs
}

func NewGetCustomerAddresses(
	db *db.Postgres,
	redis *db.Redis,
	log adapter.Logger,
	metrics adapter.Prometheus,
	tracer adapter.Tracer,
	conf *configs.Configs) *GetCustomerAddresses {
	return &GetCustomerAddresses{db, redis, log, tracer, metrics, conf}
}

func (f *GetCustomerAddresses) NewGetCustomerAddresses() *handler.GetCustomerAddressesHandler {
	query := newGetCustomerAddressesFactory(f.db, f.cache, f.log, f.tracer, f.metrics)
	return handler.NewGetAddressHandler(query, f.log, f.metrics, f.tracer)
}

func newGetCustomerAddressesFactory(
	db *db.Postgres,
	cache *db.Redis,
	log adapter.Logger,
	tracer adapter.Tracer,
	metrics adapter.Prometheus,
) query.GetCustomerAddresses {
	customerRepository := repository.NewCustomerRepository(db, metrics, tracer)
	addressRepository := repository.NewAddressRepository(db, metrics, tracer)
	redisCache := repository.NewCache(cache.Client(), metrics, tracer)
	addressesDecorator := decorator.NewCachedAddressesRepository(addressRepository, redisCache, log, tracer)
	customerService := services.NewCustomerService(customerRepository)
	return queries.NewGetCustomerAddresses(log, addressesDecorator, customerService, tracer)
}
