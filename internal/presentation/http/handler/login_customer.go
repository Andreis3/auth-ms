package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/andreis3/customers-ms/internal/app/dto"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/command"
	"github.com/andreis3/customers-ms/internal/presentation/http/helpers"
)

type LoginCustomer struct {
	log        adapter.Logger
	prometheus adapter.Prometheus
	cmd        command.Login
	tracer     adapter.Tracer
}

func NewLoginCustomer(
	log adapter.Logger,
	prometheus adapter.Prometheus,
	tracer adapter.Tracer,
	cmd command.Login,
) LoginCustomer {
	return LoginCustomer{
		log:        log,
		prometheus: prometheus,
		cmd:        cmd,
		tracer:     tracer,
	}
}

func (h *LoginCustomer) Handle(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "LoginCustomer.Handle")
	start := time.Now()
	traceID := span.SpanContext().TraceID()
	defer span.End()

	input, err := helpers.RequestDecoder[dto.LoginInput](r)
	if err != nil {
		span.RecordError(err)
		h.log.ErrorJSON("failed decode request body",
			slog.String("trace_id", traceID),
			slog.Any("error", err))
		helpers.ResponseError(w, err)
		return
	}

	res, err := h.cmd.Execute(ctx, input)
	end := time.Since(start)
	if err != nil {
		span.RecordError(err)
		h.log.ErrorJSON("failed execute create customer command",
			slog.String("trace_id", traceID),
			slog.Any("error", err))
		h.log.InfoJSON("end request", slog.String("trace_id", traceID), slog.Float64("duration", float64(end.Milliseconds())))
		helpers.ResponseError(w, err)
		return
	}

	h.prometheus.ObserveRequestDuration("/token", "http", http.StatusCreated, float64(end.Milliseconds()))
	h.log.InfoJSON("end request", slog.String("trace_id", traceID), slog.Float64("duration", float64(end.Milliseconds())))
	helpers.ResponseSuccess(w, http.StatusCreated, res)
}
