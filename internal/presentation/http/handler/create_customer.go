package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/andreis3/customers-ms/internal/app/dto"
	"github.com/andreis3/customers-ms/internal/app/mapper"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/command"
	"github.com/andreis3/customers-ms/internal/presentation/http/helpers"
)

type CreateCustomerHandler struct {
	command    command.CreateCustomer
	log        adapter.Logger
	prometheus adapter.Prometheus
	tracer     adapter.Tracer
}

func NewCreateCustomerHandler(
	cmd command.CreateCustomer,
	prometheus adapter.Prometheus,
	log adapter.Logger,
	tracer adapter.Tracer,
) *CreateCustomerHandler {
	return &CreateCustomerHandler{
		command:    cmd,
		log:        log,
		prometheus: prometheus,
		tracer:     tracer,
	}
}

func (h *CreateCustomerHandler) Handle(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var end time.Duration
	ctx, span := h.tracer.Start(r.Context(), "CreateCustomerHandler.Handle")
	traceID := span.SpanContext().TraceID()
	defer func() {
		end = time.Since(start)
		h.log.InfoJSON(
			"end request",
			slog.String("trace_id", traceID),
			slog.Float64("duration", float64(end.Milliseconds())))
		span.End()
		h.prometheus.ObserveRequestDuration("/customers", "http", http.StatusCreated, float64(end.Milliseconds()))
	}()

	input, err := helpers.RequestDecoder[dto.CreateCustomerInput](r)
	if err != nil {
		span.RecordError(err)
		h.log.ErrorJSON("failed decode request body",
			slog.String("trace_id", traceID),
			slog.Any("error", err))
		helpers.ResponseError(w, err)
		return
	}

	res, err := h.command.Execute(ctx, input)

	if err != nil {
		span.RecordError(err)
		h.log.ErrorJSON("failed execute create customer command",
			slog.String("trace_id", traceID),
			slog.Any("error", err))
		helpers.ResponseError(w, err)
		return
	}

	helpers.ResponseSuccess(w, http.StatusCreated, mapper.CustomerOutput(*res))
}
