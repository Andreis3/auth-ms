package observability

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

const (
	MeterName    = "customers-ms"
	MeterVersion = "1.0.0"
)

var (
	exporterInstance      *prometheus.Exporter
	meterProviderInstance *metric.MeterProvider
	exporterOnce          sync.Once
)

type Prometheus struct {
	provider                     *metric.MeterProvider
	counterRequestStatusCode     api.Int64Counter
	histogramInstructionDuration api.Float64Histogram
	histogramRequestDuration     api.Float64Histogram
}

func NewPrometheus() *Prometheus {
	exporterOnce.Do(func() {
		exporterInstance, _ = prometheus.New()
		meterProviderInstance = metric.NewMeterProvider(metric.WithReader(exporterInstance))
	})
	meter := meterProviderInstance.Meter(MeterName, api.WithInstrumentationVersion(MeterVersion))

	counterRequestStatusCode, _ := meter.Int64Counter("proxy_requests_total",
		api.WithDescription("Total number of proxy requests"))

	histogramInstructionDuration, _ := meter.Float64Histogram("histogram_instruction_db",
		api.WithDescription("Histogram of instruction db"),
		api.WithExplicitBucketBoundaries(
			5, 10, 15, 20, 30, 50,
			100, 200, 300, 500, 1000,
			2000, 5000, 10000, 20000,
			30000, 50000, 100000))

	histogramRequestDuration, _ := meter.Float64Histogram("histogram_instruction_request",
		api.WithDescription("Histogram of instruction request"),
		api.WithExplicitBucketBoundaries(
			5, 10, 15, 20, 30, 50,
			100, 200, 300, 500, 1000,
			2000, 5000, 10000, 20000,
			30000, 50000, 100000))

	return &Prometheus{
		provider:                     meterProviderInstance,
		counterRequestStatusCode:     counterRequestStatusCode,
		histogramInstructionDuration: histogramInstructionDuration,
		histogramRequestDuration:     histogramRequestDuration,
	}
}

func (p *Prometheus) CounterRequestStatusCode(router, protocol string, statusCode int) {
	opt := api.WithAttributes(
		attribute.Key("router").String(router),
		attribute.Key("status_code").Int(statusCode),
		attribute.Key("protocol").String(protocol),
	)
	p.counterRequestStatusCode.Add(context.Background(), 1, opt)
}

func (p *Prometheus) ObserveInstructionDBDuration(database, table, method string, duration float64) {
	opt := api.WithAttributes(
		attribute.Key("database").String(database),
		attribute.Key("table").String(table),
		attribute.Key("method").String(method),
	)
	p.histogramInstructionDuration.Record(context.Background(), duration, opt)
}

func (p *Prometheus) ObserveRequestDuration(router, protocol string, statusCode int, duration float64) {
	opt := api.WithAttributes(
		attribute.Key("router").String(router),
		attribute.Key("status_code").Int(statusCode),
		attribute.Key("protocol").String(protocol),
	)
	p.histogramRequestDuration.Record(context.Background(), duration, opt)
}

func (p *Prometheus) Close() {
	_ = p.provider.Shutdown(context.Background())
}

func (p *Prometheus) MeterProvider() *metric.MeterProvider {
	return p.provider
}
