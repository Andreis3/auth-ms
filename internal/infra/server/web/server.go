package web

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

	"github.com/andreis3/customers-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/customers-ms/internal/infra/adapters/db"
	"github.com/andreis3/customers-ms/internal/infra/adapters/logger"
	"github.com/andreis3/customers-ms/internal/infra/adapters/observability"
	"github.com/andreis3/customers-ms/internal/infra/configs"
	"github.com/andreis3/customers-ms/internal/infra/server/web/routes"
	"github.com/andreis3/customers-ms/internal/util"
)

type Server struct {
	HTTPServer *http.Server
	Postgres   *db.Postgres
	Log        logger.Logger
	Prometheus *observability.Prometheus
	Tracer     adapter.Tracer
}

func NewServer(conf *configs.Configs, log logger.Logger) *Server {
	start := time.Now()

	prometheus := observability.NewPrometheus()
	pool := db.NewPoolConnections(conf, prometheus)

	redis := db.NewRedis(*conf)

	tracer, _ := observability.InitOtelTracer(context.Background(), "customers-ms")

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
