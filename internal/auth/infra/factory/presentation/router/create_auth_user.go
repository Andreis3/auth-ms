package router

import (
	"github.com/andreis3/auth-ms/internal/auth/adapter/input/http/middlewares"
	"github.com/andreis3/auth-ms/internal/auth/adapter/input/http/routes"
	"github.com/andreis3/auth-ms/internal/auth/domain/interfaces/adapter"
	"github.com/andreis3/auth-ms/internal/auth/infra/config"
	"github.com/andreis3/auth-ms/internal/auth/infra/db"
	"github.com/andreis3/auth-ms/internal/auth/infra/factory/presentation/handler"
)

func MakeCreateAuthUserRouter(
	postgres *db.Postgres,
	redis *db.Redis,
	log adapter.Logger,
	prometheus adapter.Prometheus,
	tracer adapter.Tracer,
	conf *config.Configs) *routes.User {

	loggingMiddleware := middlewares.NewLoggingMiddleware(log, tracer)

	createAuthUserHandler := handler.NewCreateAuthUser(postgres, redis, log, prometheus, tracer, conf)
	customerRoutes := routes.NewUser(
		createAuthUserHandler,
		loggingMiddleware,
	)
	return customerRoutes
}
