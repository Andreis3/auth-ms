package routes

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/andreis3/customers-ms/internal/presentation/http/helpers"
)

type Metrics struct{}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) Routes() helpers.RouteType {
	return helpers.RouteType{
		{
			Method:      http.MethodGet,
			Path:        "/metrics",
			Handler:     promhttp.Handler(),
			Description: "Metrics Prometheus",
			Middlewares: helpers.Middlewares{},
		},
	}
}
