package iadapter

import "go.opentelemetry.io/otel/sdk/metric"

type IPrometheus interface {
	CounterRequestStatusCode(router, protocol string, statusCode int)
	ObserveInstructionDBDuration(database, table, method string, duration float64)
	ObserveRequestDuration(router, protocol string, statusCode int, status string, duration float64)
	Close()
	MeterProvider() *metric.MeterProvider
}
