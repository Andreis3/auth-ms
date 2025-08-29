package router

import (
	"github.com/andreis3/auth-ms/internal/adapter/input/http/middlewares"
	"github.com/andreis3/auth-ms/internal/adapter/input/http/routes"
	adapter2 "github.com/andreis3/auth-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/auth-ms/internal/infra/config"
	db2 "github.com/andreis3/auth-ms/internal/infra/db"
	"github.com/andreis3/auth-ms/internal/infra/factory/http/handler"
)

func MakeCreateAuthUserRouter(
	postgres *db2.Postgres,
	redis *db2.Redis,
	log adapter2.Logger,
	prometheus adapter2.Prometheus,
	tracer adapter2.Tracer,
	conf *config.Configs) *routes.User {

	loggingMiddleware := middlewares.NewLoggingMiddleware(log, tracer)

	createAuthUserHandler := handler.NewCreateAuthUser(postgres, redis, log, prometheus, tracer, conf)
	customerRoutes := routes.NewUser(
		createAuthUserHandler,
		loggingMiddleware,
	)
	return customerRoutes
}
