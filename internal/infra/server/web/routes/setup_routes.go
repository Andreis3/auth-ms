package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/andreis3/auth-ms/internal/domain/interfaces/iadapter"
	"github.com/andreis3/auth-ms/internal/infra/adapters/db"
	"github.com/andreis3/auth-ms/internal/infra/configs"
	"github.com/andreis3/auth-ms/internal/infra/factories/presentation/router-factory"
	"github.com/andreis3/auth-ms/internal/presentation/http/routes"
)

type RegisterRoutesDeps struct {
	Mux        *chi.Mux
	PostgresDB *db.Postgres
	Redis      *db.Redis
	Log        iadapter.ILogger
	Prometheus iadapter.IPrometheus
	Conf       *configs.Configs
	Tracer     iadapter.ITracer
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
		router_factory.MakeCreateAuthUserRouter(deps.PostgresDB, deps.Redis, deps.Log, deps.Prometheus, deps.Tracer, deps.Conf),
	}
}
