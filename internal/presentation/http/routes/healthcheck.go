package routes

import (
	"net/http"

	"github.com/andreis3/customers-ms/internal/presentation/http/handler"
	"github.com/andreis3/customers-ms/internal/presentation/http/helpers"
)

type HealthCheck struct{}

func NewHealthCheck() *HealthCheck {
	return &HealthCheck{}
}

func (r *HealthCheck) Routes() helpers.RouteType {
	return helpers.RouteType{
		{
			Method:      http.MethodGet,
			Path:        "/health",
			Handler:     handler.HealthCheck(),
			Description: "Health Check",
			Middlewares: helpers.Middlewares{},
		},
	}
}
