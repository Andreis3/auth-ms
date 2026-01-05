package router

import (
	"github.com/andreis3/auth-ms/internal/adapter/input/http/middlewares"
	"github.com/andreis3/auth-ms/internal/adapter/input/http/routes"
	"github.com/andreis3/auth-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/auth-ms/internal/infra/config"
	"github.com/andreis3/auth-ms/internal/infra/db"
	"github.com/andreis3/auth-ms/internal/infra/factory/http/handler"
)

func MakeCreateAuthUserRouter(
	postgres *db.Postgres,
	redis *db.Redis,
	log adapter.Logger,
	prometheus adapter.Prometheus,
	tracer adapter.Tracer,
	conf *config.Configs,
) *routes.User {
	loggingMiddleware := middlewares.NewLoggingMiddleware(log, tracer)

	createAuthUserHandler := handler.NewCreateAuthUser(postgres, redis, log, prometheus, tracer, conf)
	customerRoutes := routes.NewUser(
		createAuthUserHandler,
		loggingMiddleware,
	)
	return customerRoutes
}
