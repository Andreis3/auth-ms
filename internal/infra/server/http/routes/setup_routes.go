package routes

import (
	"github.com/go-chi/chi/v5"

	routes2 "github.com/andreis3/auth-ms/internal/adapter/input/http/routes"
	adapter2 "github.com/andreis3/auth-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/auth-ms/internal/infra/config"
	db2 "github.com/andreis3/auth-ms/internal/infra/db"
	"github.com/andreis3/auth-ms/internal/infra/factory/http/router"
)

type RegisterRoutesDeps struct {
	Mux        *chi.Mux
	PostgresDB *db2.Postgres
	Redis      *db2.Redis
	Log        adapter2.Logger
	Prometheus adapter2.Prometheus
	Conf       *config.Configs
	Tracer     adapter2.Tracer
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
		routes2.NewHealthCheck(),
		routes2.NewMetrics(),
		router.MakeCreateAuthUserRouter(deps.PostgresDB, deps.Redis, deps.Log, deps.Prometheus, deps.Tracer, deps.Conf),
	}
}
