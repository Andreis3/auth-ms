package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/andreis3/auth-ms/internal/adapter/input/http/routes"
	"github.com/andreis3/auth-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/auth-ms/internal/infra/config"
	"github.com/andreis3/auth-ms/internal/infra/db"
	"github.com/andreis3/auth-ms/internal/infra/factory/http/router"
)

type RegisterRoutesDeps struct {
	Mux        *chi.Mux
	PostgresDB *db.Postgres
	Redis      *db.Redis
	Log        adapter.Logger
	Prometheus adapter.Prometheus
	Conf       *config.Configs
	Tracer     adapter.Tracer
}

func Setup(deps *RegisterRoutesDeps) {
	registerRoutes := NewRegisterRoutes(
		deps.Mux,
		deps.Log,
		BuildRoutes(deps),
	)

	registerRoutes.Register()
}

func BuildRoutes(deps *RegisterRoutesDeps) []ModuleRoutes {
	return []ModuleRoutes{
		routes.NewHealthCheck(),
		routes.NewMetrics(),
		router.MakeCreateAuthUserRouter(deps.PostgresDB, deps.Redis, deps.Log, deps.Prometheus, deps.Tracer, deps.Conf),
	}
}
