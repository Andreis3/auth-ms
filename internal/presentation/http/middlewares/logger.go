package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/andreis3/auth-ms/internal/domain/interfaces/iadapter"
)

type Logging struct {
	logger iadapter.ILogger
	tracer iadapter.ITracer
}

func NewLoggingMiddleware(logger iadapter.ILogger, tracer iadapter.ITracer) *Logging {
	return &Logging{
		logger: logger,
		tracer: tracer,
	}
}

func (l *Logging) LoggingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, span := l.tracer.Start(r.Context(), "HTTP "+r.Method+" "+r.URL.Path)
			defer span.End()
			log := l.logger.WithTrace(ctx)

			log.Info("new request received",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
