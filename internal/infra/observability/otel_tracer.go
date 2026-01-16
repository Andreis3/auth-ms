package observability

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/andreis3/auth-ms/internal/domain/interfaces/adapter"
)

type Tracer struct {
	tracer trace.Tracer
}

func (o *Tracer) Start(ctx context.Context, spanName string) (context.Context, adapter.Span) {
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
	if err == nil {
		return
	}

	s.Span.RecordError(err)

	// Melhora visualização no Tempo/Jaeger (status != apenas evento de erro)
	s.SetStatus(codes.Error, err.Error())
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

type OtelConfig struct {
	ServiceName     string
	HostTracer      string
	Insecure        bool
	CompressionGzip bool
	SampleRatio     float64
	ShutdownTimeout time.Duration

	// opcionais (mas uteis)
	ServiceVersion string //  git sha / semver
	Environment    string // dev, staging, production, etc.
	SetGlobal      bool
}

// ✅ Esta função inicializa o OpenTelemetry por completo e retorna o iadapter.Tracer
func InitializeOtelTracer(ctx context.Context, cfg OtelConfig) (adapter.Tracer, func(context.Context) error, error) {
	if cfg.ServiceName == "" {
		cfg.ServiceName = "unknown_service"
	}

	if cfg.HostTracer == "" {
		cfg.HostTracer = "localhost:4318"
	}

	if cfg.SampleRatio <= 0 || cfg.SampleRatio > 1 {
		cfg.SampleRatio = 1.0
	}

	if cfg.ShutdownTimeout <= 0 {
		cfg.ShutdownTimeout = 5 * time.Second
	}

	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(cfg.HostTracer),
	}

	if cfg.Insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	if cfg.CompressionGzip {
		opts = append(opts, otlptracehttp.WithCompression(otlptracehttp.GzipCompression))
	}

	exporter, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return nil, nil, err
	}

	res, err := sdkresource.Merge(sdkresource.Default(),
		sdkresource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
			semconv.DeploymentEnvironmentName(cfg.Environment),
		),
	)
	if err != nil {
		return nil, nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.SampleRatio)),
	)

	if cfg.SetGlobal {
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		))
	}

	shutdown := func(ctx context.Context) error {
		shutdownCtx, cancel := context.WithTimeout(ctx, cfg.ShutdownTimeout)
		defer cancel()
		return tp.Shutdown(shutdownCtx)
	}

	// Se SetGlobal=true, usar otel.Tracer pega  o provider SetGlobal
	var tr trace.Tracer
	if cfg.SetGlobal {
		tr = otel.Tracer(cfg.ServiceName)
	} else {
		tr = tp.Tracer(cfg.ServiceName)
	}

	log.Printf("OpenTelemetry initialized for service: %s, endpoint: %s", cfg.ServiceName, cfg.HostTracer)

	return &Tracer{tracer: tr}, shutdown, nil
}
