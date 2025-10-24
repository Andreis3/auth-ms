//go:build unit

package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/andreis3/auth-ms/internal/adapter/input/http/handler"
	"github.com/andreis3/auth-ms/internal/app/dto"
	"github.com/andreis3/auth-ms/internal/domain/errors"
	"github.com/andreis3/auth-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/auth-ms/tests/mocks/app/mcommand"
	"github.com/andreis3/auth-ms/tests/mocks/infra/madapters"
)

var _ = Describe("INTERNAL :: ADAPTER :: INPUT :: HTTP :: HANDLER :: CREATE_AUTH_USER", func() {
	const route = "/auth/signup"

	var (
		tracer      *madapters.TracerMock
		span        *madapters.SpanMock
		spanCtx     *madapters.SpanContextMock
		logger      *madapters.LoggerMock
		prometheus  *madapters.PrometheusMock
		commandMock *mcommand.CreateAuthUserCommandMock
	)

	BeforeEach(func() {
		tracer = &madapters.TracerMock{}
		span = &madapters.SpanMock{}
		spanCtx = &madapters.SpanContextMock{}
		logger = &madapters.LoggerMock{}
		prometheus = &madapters.PrometheusMock{}
		commandMock = &mcommand.CreateAuthUserCommandMock{}
	})

	newHandler := func() *handler.CreateAuthUserHandler {
		return handler.NewCreateAuthUserHandler(commandMock, prometheus, logger, tracer)
	}

	Context("success cases", func() {
		It("should record request duration metrics with success status", func() {
			body := `{"email":"user@example.com","password":"Sup3r$ecret","password_confirm":"Sup3r$ecret","name":"User"}`
			req := httptest.NewRequest(http.MethodPost, route, strings.NewReader(body))
			w := httptest.NewRecorder()

			tracer.On("Start", req.Context(), "CreateAuthUserHandler.Handle").Return(req.Context(), adapter.Span(span))
			span.On("SpanContext").Return(adapter.SpanContext(spanCtx))
			spanCtx.On("TraceID").Return("trace-123")
			span.On("End").Return()

			logger.On("InfoJSON", "end request", mock.AnythingOfType("slog.Attr"), mock.AnythingOfType("slog.Attr")).Return()

			output := &dto.CreateAuthUserOutput{PublicID: "user-1", Email: "user@example.com", Name: "User"}
			commandMock.On("Execute", req.Context(), mock.MatchedBy(func(input dto.CreateAuthUserInput) bool {
				return input.Email == "user@example.com" && input.Password == "Sup3r$ecret" && input.PasswordConfirm == "Sup3r$ecret" && input.Name == "User"
			})).Return(output, (*errors.Error)(nil))

			prometheus.On("ObserveRequestDuration", route, "http", http.StatusCreated, "success", mock.MatchedBy(func(value float64) bool {
				return value >= 0
			})).Return()

			newHandler().Handle(w, req)

			Expect(w.Code).To(Equal(http.StatusCreated))
			Expect(prometheus.AssertExpectations(GinkgoT())).To(BeTrue())
			Expect(commandMock.AssertExpectations(GinkgoT())).To(BeTrue())
			Expect(tracer.AssertExpectations(GinkgoT())).To(BeTrue())
			Expect(span.AssertExpectations(GinkgoT())).To(BeTrue())
			Expect(spanCtx.AssertExpectations(GinkgoT())).To(BeTrue())
			Expect(logger.AssertExpectations(GinkgoT())).To(BeTrue())
		})
	})

	Context("error cases", func() {
		It("should record request duration metrics with decoder error status", func() {
			req := httptest.NewRequest(http.MethodPost, route, strings.NewReader("invalid"))
			w := httptest.NewRecorder()

			tracer.On("Start", req.Context(), "CreateAuthUserHandler.Handle").Return(req.Context(), adapter.Span(span))
			span.On("SpanContext").Return(adapter.SpanContext(spanCtx))
			spanCtx.On("TraceID").Return("trace-err-1")
			span.On("RecordError", mock.Anything).Return()
			span.On("End").Return()

			logger.On("InfoJSON", "end request", mock.AnythingOfType("slog.Attr"), mock.AnythingOfType("slog.Attr")).Return()
			logger.On("ErrorJSON", "failed decode request body", mock.AnythingOfType("slog.Attr"), mock.AnythingOfType("slog.Attr")).Return()

			prometheus.On("ObserveRequestDuration", route, "http", http.StatusBadRequest, "error", mock.MatchedBy(func(value float64) bool {
				return value >= 0
			})).Return()

			newHandler().Handle(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
			Expect(prometheus.AssertExpectations(GinkgoT())).To(BeTrue())
			Expect(tracer.AssertExpectations(GinkgoT())).To(BeTrue())
			Expect(span.AssertExpectations(GinkgoT())).To(BeTrue())
			Expect(spanCtx.AssertExpectations(GinkgoT())).To(BeTrue())
			Expect(logger.AssertExpectations(GinkgoT())).To(BeTrue())
		})

		It("should record request duration metrics with command error status", func() {
			body := `{"email":"user@example.com","password":"Sup3r$ecret","password_confirm":"Sup3r$ecret","name":"User"}`
			req := httptest.NewRequest(http.MethodPost, route, strings.NewReader(body))
			w := httptest.NewRecorder()

			tracer.On("Start", req.Context(), "CreateAuthUserHandler.Handle").Return(req.Context(), adapter.Span(span))
			span.On("SpanContext").Return(adapter.SpanContext(spanCtx))
			spanCtx.On("TraceID").Return("trace-err-2")
			span.On("End").Return()

			logger.On("InfoJSON", "end request", mock.AnythingOfType("slog.Attr"), mock.AnythingOfType("slog.Attr")).Return()

			conflictErr := errors.New(errors.ErrConflict, "conflict")
			prometheus.On("ObserveRequestDuration", route, "http", http.StatusConflict, "error", mock.MatchedBy(func(value float64) bool {
				return value >= 0
			})).Return()

			commandMock.On("Execute", req.Context(), mock.AnythingOfType("dto.CreateAuthUserInput")).Return((*dto.CreateAuthUserOutput)(nil), conflictErr)

			newHandler().Handle(w, req)

			Expect(w.Code).To(Equal(http.StatusConflict))
			Expect(prometheus.AssertExpectations(GinkgoT())).To(BeTrue())
			Expect(commandMock.AssertExpectations(GinkgoT())).To(BeTrue())
			Expect(tracer.AssertExpectations(GinkgoT())).To(BeTrue())
			Expect(span.AssertExpectations(GinkgoT())).To(BeTrue())
			Expect(spanCtx.AssertExpectations(GinkgoT())).To(BeTrue())
			Expect(logger.AssertExpectations(GinkgoT())).To(BeTrue())
		})
	})
})
