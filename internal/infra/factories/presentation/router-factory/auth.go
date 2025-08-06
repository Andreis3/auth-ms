package router_factory

import (
	"github.com/andreis3/customers-ms/internal/app/services"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/customers-ms/internal/infra/adapters/crypto"
	"github.com/andreis3/customers-ms/internal/infra/adapters/jwt"
	"github.com/andreis3/customers-ms/internal/infra/configs"
	"github.com/andreis3/customers-ms/internal/infra/factories/app/command-factory"
	"github.com/andreis3/customers-ms/internal/infra/repositories/repository"
	"github.com/andreis3/customers-ms/internal/presentation/http/handler"
	"github.com/andreis3/customers-ms/internal/presentation/http/middlewares"
	"github.com/andreis3/customers-ms/internal/presentation/http/routes"
)

func MakeAuthRouter(postgres adapter.Postgres, log adapter.Logger, prometheus adapter.Prometheus, conf *configs.Configs, tracer adapter.Tracer) *routes.LoginRoutes {
	customerRepository := repository.NewCustomerRepository(postgres, prometheus, tracer)
	tokenService := jwt.NewJWT(conf)
	authService := services.NewAuthService(tokenService, customerRepository)
	bcrypt := crypto.NewBcrypt()
	commands := command_factory.NewAuthenticateCustomerFactory(log, customerRepository, authService, bcrypt, tracer)
	authHandler := handler.NewLoginCustomer(log, prometheus, tracer, commands)
	loggingMiddleware := middlewares.NewLoggingMiddleware(log, tracer)
	return routes.NewLoginRoutes(authHandler, loggingMiddleware)
}
