package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/andreis3/auth-ms/internal/auth/adapter/input/http/helpers"
	"github.com/andreis3/auth-ms/internal/auth/app/dto"
	"github.com/andreis3/auth-ms/internal/auth/app/port/command"
	"github.com/andreis3/auth-ms/internal/auth/domain/interfaces/adapter"
)

type CreateAuthUserHandler struct {
	command    command.CreateAuthUser
	log        adapter.Logger
	prometheus adapter.Prometheus
	tracer     adapter.Tracer
}

func NewCreateAuthUserHandler(
	cmd command.CreateAuthUser,
	prometheus adapter.Prometheus,
	log adapter.Logger,
	tracer adapter.Tracer,
) *CreateAuthUserHandler {
	return &CreateAuthUserHandler{
		command:    cmd,
		log:        log,
		prometheus: prometheus,
		tracer:     tracer,
	}
}

func (h *CreateAuthUserHandler) Handle(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var end time.Duration
	ctx, span := h.tracer.Start(r.Context(), "CreateAuthUserHandler.Handle")
	traceID := span.SpanContext().TraceID()
	defer func() {
		end = time.Since(start)
		h.log.InfoJSON(
			"end request",
			slog.String("trace_id", traceID),
			slog.Float64("duration", float64(end.Milliseconds())))
		span.End()
	}()

	input, err := helpers.RequestDecoder[dto.CreateAuthUserInput](r)
	if err != nil {
		span.RecordError(err)
		h.log.ErrorJSON("failed decode request body",
			slog.String("trace_id", traceID),
			slog.Any("error", err))
		helpers.ResponseError(w, err)
		h.prometheus.ObserveRequestDuration("/auth/signup", "http", http.StatusCreated, "error", float64(end.Milliseconds()))
		return
	}

	res, err := h.command.Execute(ctx, input)
	if err != nil {
		helpers.ResponseError(w, err)
		h.prometheus.ObserveRequestDuration("/auth/signup", "http", http.StatusCreated, "error", float64(end.Milliseconds()))
		return
	}

	helpers.ResponseSuccess(w, http.StatusCreated, res)
	h.prometheus.ObserveRequestDuration("/auth/signup", "http", http.StatusCreated, "success", float64(end.Milliseconds()))
}
