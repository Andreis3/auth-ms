package middleware_factory

import (
	"net/http"

	"github.com/andreis3/customers-ms/internal/app/services"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/customers-ms/internal/infra/adapters/db"
	"github.com/andreis3/customers-ms/internal/infra/adapters/jwt"
	"github.com/andreis3/customers-ms/internal/infra/configs"
	"github.com/andreis3/customers-ms/internal/infra/repositories/repository"
	"github.com/andreis3/customers-ms/internal/presentation/http/middlewares"
)

type ValidateCustomerFactory struct {
	postgres   *db.Postgres
	prometheus adapter.Prometheus
	logger     adapter.Logger
	tracer     adapter.Tracer
	conf       *configs.Configs
}

func NewValidateCustomerFactory(
	postgres *db.Postgres,
	prometheus adapter.Prometheus,
	logger adapter.Logger,
	tracer adapter.Tracer,
	conf *configs.Configs) *ValidateCustomerFactory {
	return &ValidateCustomerFactory{
		postgres:   postgres,
		prometheus: prometheus,
		logger:     logger,
		tracer:     tracer,
		conf:       conf,
	}
}

func (f *ValidateCustomerFactory) NewValidateCustomer() *middlewares.ValidateCustomer {
	jwt := jwt.NewJWT(f.conf)
	customerRepository := repository.NewCustomerRepository(f.postgres, f.prometheus, f.tracer)
	authService := services.NewAuthService(jwt, customerRepository)
	return middlewares.NewValidateCustomerMiddleware(authService, f.logger, f.tracer)
}

func MakeValidateCustomerMiddleware(factory *ValidateCustomerFactory) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m := factory.NewValidateCustomer()
			m.ValidateCustomer()(next).ServeHTTP(w, r)
		})
	}
}
