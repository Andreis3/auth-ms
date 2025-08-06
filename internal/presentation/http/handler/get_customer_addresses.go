package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/andreis3/customers-ms/internal/app/dto"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/query"
	"github.com/andreis3/customers-ms/internal/presentation/http/helpers"
)

type GetCustomerAddressesHandler struct {
	query      query.GetCustomerAddresses
	log        adapter.Logger
	prometheus adapter.Prometheus
	tracer     adapter.Tracer
}

func NewGetAddressHandler(
	query query.GetCustomerAddresses,
	log adapter.Logger,
	prometheus adapter.Prometheus,
	tracer adapter.Tracer,
) *GetCustomerAddressesHandler {
	return &GetCustomerAddressesHandler{
		query:      query,
		log:        log,
		prometheus: prometheus,
		tracer:     tracer,
	}
}

func (h *GetCustomerAddressesHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "GetCustomerAddressesHandler.Handle")
	start := time.Now()
	var end time.Duration
	traceID := span.SpanContext().TraceID()

	defer func() {
		end = time.Since(start)
		h.log.InfoJSON(
			"end request",
			slog.String("trace_id", traceID),
			slog.Float64("duration", float64(end.Milliseconds())))
		span.End()
		h.prometheus.ObserveRequestDuration("/", "http", http.StatusOK, float64(end.Milliseconds()))
	}()

	customerID, _ := ctx.Value("customer_id").(int64)
	email, _ := ctx.Value("email").(string)

	input, err := h.query.Execute(ctx, dto.GetCustomerAddressesInput{
		CustomerID: customerID,
		Email:      email,
	})

	if err != nil {
		span.RecordError(err)
		h.log.ErrorJSON("failed execute get address query",
			slog.String("trace_id", traceID),
			slog.Any("error", err))
		helpers.ResponseError(w, err)
		return
	}

	helpers.ResponseSuccess[any](w, http.StatusOK, input)
}
