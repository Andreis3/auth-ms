package handler

import (
	"log/slog"
	"net/http"
	"time"

	helpers2 "github.com/andreis3/auth-ms/internal/adapter/input/http/helpers"
	"github.com/andreis3/auth-ms/internal/app/dto"
	"github.com/andreis3/auth-ms/internal/app/port/command"
	adapter2 "github.com/andreis3/auth-ms/internal/domain/interfaces/adapter"
)

type CreateAuthUserHandler struct {
	command    command.CreateAuthUser
	log        adapter2.Logger
	prometheus adapter2.Prometheus
	tracer     adapter2.Tracer
}

func NewCreateAuthUserHandler(
	cmd command.CreateAuthUser,
	prometheus adapter2.Prometheus,
	log adapter2.Logger,
	tracer adapter2.Tracer,
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
	ctx, span := h.tracer.Start(r.Context(), "CreateAuthUserHandler.Handle")
	traceID := span.SpanContext().TraceID()
	defer func() {
		end := time.Since(start)
		h.log.InfoJSON(
			"end request",
			slog.String("trace_id", traceID),
			slog.Float64("duration", float64(end.Milliseconds())))
		span.End()
	}()

	input, err := helpers2.RequestDecoder[dto.CreateAuthUserInput](r)
	if err != nil {
		span.RecordError(err)
		h.log.ErrorJSON("failed decode request body",
			slog.String("trace_id", traceID),
			slog.Any("error", err))
		status := helpers2.ResponseError(w, err)
		duration := time.Since(start)
		h.prometheus.ObserveRequestDuration("/auth/signup", "http", status, "error", float64(duration.Milliseconds()))
		return
	}

	res, err := h.command.Execute(ctx, input)
	if err != nil {
		status := helpers2.ResponseError(w, err)
		duration := time.Since(start)
		h.prometheus.ObserveRequestDuration("/auth/signup", "http", status, "error", float64(duration.Milliseconds()))
		return
	}

	helpers2.ResponseSuccess(w, http.StatusCreated, res)
	duration := time.Since(start)
	h.prometheus.ObserveRequestDuration("/auth/signup", "http", http.StatusCreated, "success", float64(duration.Milliseconds()))
}
