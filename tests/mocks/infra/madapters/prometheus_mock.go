package madapters

import (
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/sdk/metric"

	"github.com/andreis3/auth-ms/internal/domain/interfaces/adapter"
)

type PrometheusMock struct{ mock.Mock }

func (p *PrometheusMock) CounterRequestStatusCode(router, protocol string, statusCode int) {
	p.Called(router, protocol, statusCode)
}

func (p *PrometheusMock) ObserveInstructionDBDuration(database, table, method string, duration float64) {
	p.Called(database, table, method, duration)
}

func (p *PrometheusMock) ObserveRequestDuration(router, protocol string, statusCode int, status string, duration float64) {
	p.Called(router, protocol, statusCode, status, duration)
}

func (p *PrometheusMock) Close() {}

func (p *PrometheusMock) MeterProvider() *metric.MeterProvider {
	args := p.Called()
	if v := args.Get(0); v != nil {
		return v.(*metric.MeterProvider)
	}

	return nil
}

var _ adapter.Prometheus = (*PrometheusMock)(nil)
