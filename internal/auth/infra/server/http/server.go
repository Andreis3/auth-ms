package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/andreis3/auth-ms/internal/auth/adapter/output/logger"
	"github.com/andreis3/auth-ms/internal/auth/adapter/output/metrics"
	"github.com/andreis3/auth-ms/internal/auth/adapter/output/tracer"
	"github.com/andreis3/auth-ms/internal/auth/domain/interfaces/adapter"
	"github.com/andreis3/auth-ms/internal/auth/infra/config"
	db2 "github.com/andreis3/auth-ms/internal/auth/infra/db"
	"github.com/andreis3/auth-ms/internal/auth/infra/server/http/routes"
	"github.com/andreis3/auth-ms/internal/auth/util"
)

type Server struct {
	HTTPServer *http.Server
	Postgres   *db2.Postgres
	Log        logger.Logger
	Prometheus *metrics.Prometheus
	Tracer     adapter.Tracer
}

func NewServer(conf *config.Configs, log logger.Logger) *Server {
	start := time.Now()

	prometheus := metrics.NewPrometheus()
	pool := db2.NewPoolConnections(conf, prometheus)

	redis := db2.NewRedis(*conf)

	tracer, _ := tracer.InitOtelTracer(context.Background(), "customers-ms")

	mux := chi.NewRouter()

	// OpenTelemetry Middleware
	mux.Use(func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, "customers-ms")
	})

	setupRoutesInput := routes.RegisterRoutesDeps{
		Mux:        mux,
		PostgresDB: pool,
		Redis:      redis,
		Log:        &log,
		Prometheus: prometheus,
		Conf:       conf,
		Tracer:     tracer,
	}

	routes.Setup(&setupRoutesInput)

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", conf.ServerPort),
		Handler: mux,
	}

	log.InfoText("[Server] ", "SERVER_STARTED", fmt.Sprintf("Server started in %s", time.Since(start)))
	log.InfoText("[Server] ", "SERVER_STARTED", fmt.Sprintf("Server address http://localhost:%s", conf.ServerPort))

	return &Server{
		HTTPServer: server,
		Postgres:   pool,
		Log:        log,
		Prometheus: prometheus,
	}
}

func (s *Server) Start() {
	if err := s.HTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.Log.CriticalText("[Server] ", "SERVER_ERROR", err.Error())
		os.Exit(util.ExitFailure)
	}
}

func (s *Server) GracefulShutdown() {
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-shutdownSignal
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.Log.InfoText("[Server] ", "SERVER_SHUTDOWN", "Server is shutting down...")

	if err := s.HTTPServer.Shutdown(ctx); err != nil {
		s.Log.ErrorText("[Server] ", "SERVER_SHUTDOWN", err.Error())
	}
	s.Log.InfoText("Closing postgres connection...")
	s.Postgres.Close()
	s.Log.InfoText("Closing prometheus...")
	s.Prometheus.Close()
	s.Log.InfoText("Shutdown complete exit code 0...")

	os.Exit(util.ExitSuccess)
}
