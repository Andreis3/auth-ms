package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/andreis3/customers-ms/internal/domain/errors"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/service"
	"github.com/andreis3/customers-ms/internal/presentation/http/helpers"
)

type ValidateCustomer struct {
	authService service.Auth
	logger      adapter.Logger
	tracer      adapter.Tracer
}

func NewValidateCustomerMiddleware(authService service.Auth, logger adapter.Logger, tracer adapter.Tracer) *ValidateCustomer {
	return &ValidateCustomer{
		authService: authService,
		logger:      logger,
		tracer:      tracer,
	}
}

func (v *ValidateCustomer) ValidateCustomer() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, span := v.tracer.Start(r.Context(), "ValidateCustomer")
			defer span.End()
			traceID := span.SpanContext().TraceID()
			v.logger.InfoJSON("middleware validate token started",
				slog.String("trace_id", traceID),
				slog.String("path", r.URL.Path))

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				v.logger.ErrorJSON("missing authorization header",
					slog.String("trace_id", traceID),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path))
				helpers.ResponseError(w, errors.ErrorInvalidToken())
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				err := errors.ErrorInvalidToken()
				v.logger.ErrorJSON("invalid authorization header",
					slog.String("trace_id", traceID),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Any("error", err))
				helpers.ResponseError(w, err)
				return
			}
			token := parts[1]
			v.logger.InfoJSON("validating token",
				slog.String("trace_id", traceID),
				slog.String("token", token))
			claims, err := v.authService.DecodeToken(ctx, token)
			if err != nil {
				v.logger.ErrorJSON("failed validate token",
					slog.String("trace_id", traceID),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Any("error", err))
				helpers.ResponseError(w, errors.ErrorInvalidToken())
				return
			}

			ctx = context.WithValue(ctx, "customer_id", claims.CustomerID)
			ctx = context.WithValue(ctx, "email", claims.Email)

			v.logger.InfoJSON("token validated successfully",
				slog.String("trace_id", traceID))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
