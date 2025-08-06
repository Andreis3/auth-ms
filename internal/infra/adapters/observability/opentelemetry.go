// internal/infra/adapters/observability/otel_tracer.go
package observability

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/andreis3/customers-ms/internal/domain/interfaces/adapter"
)

type otelTracer struct {
	tracer trace.Tracer
}

func (o *otelTracer) Start(ctx context.Context, spanName string) (context.Context, adapter.Span) {
	ctx, s := o.tracer.Start(ctx, spanName)
	return ctx, &otelSpan{s}
}

type otelSpan struct {
	trace.Span
}

func (s *otelSpan) End() {
	s.Span.End()
}

func (s *otelSpan) RecordError(err error) {
	s.Span.RecordError(err)
}

func (s *otelSpan) SpanContext() adapter.SpanContext {
	return &otelSpanContext{s.Span.SpanContext()} // ✅ chama o SpanContext do campo do OTEL
}

type otelSpanContext struct {
	trace.SpanContext
}

func (sc *otelSpanContext) TraceID() string {
	return sc.SpanContext.TraceID().String()
}

// ✅ Esta função inicializa o OpenTelemetry por completo e retorna o adapter.Tracer
func InitOtelTracer(ctx context.Context, serviceName string) (adapter.Tracer, func(context.Context) error) {
	exporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
	)
	if err != nil {
		log.Fatalf("failed to create OTLP exporter: %v", err)
	}

	resource, err := sdkresource.Merge(
		sdkresource.Default(),
		sdkresource.NewWithAttributes(
			"", //semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.1)),
	)

	otel.SetTracerProvider(provider)

	// Retorna o adapter.Tracer e a função de shutdown
	return &otelTracer{tracer: provider.Tracer(serviceName)}, provider.Shutdown
}
