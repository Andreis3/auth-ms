//go:build unit

package service_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/andreis3/auth-ms/internal/app/service"
	"github.com/andreis3/auth-ms/internal/domain/entity"
	"github.com/andreis3/auth-ms/internal/domain/errors"
	"github.com/andreis3/auth-ms/internal/domain/interfaces/iadapter"
	"github.com/andreis3/auth-ms/tests/mocks/infra/madapters"
	"github.com/andreis3/auth-ms/tests/mocks/infra/mrepository"
)

var _ = Describe("INTERNAL :: APP :: SERVICE :: USER_SERVICE", func() {
	Describe("#Execute", func() {
		Context("success cases", func() {
			It("should not return an error when not exist user with email equal", func() {
				ctx := context.Background()
				email := "test@test.com"

				repo := new(mrepository.UserRepositoryMock)
				tracer := new(madapters.TracerMock)
				span := new(madapters.SpanMock)
				sc := new(madapters.SpanContextMock)
				log := new(madapters.LoggerMock)

				tracer.On("Start", ctx, "UserService.ValidateEmailAvailability").
					Return(ctx, iadapter.Span(span))
				sc.On("TraceID").Return("trace-123")
				span.On("SpanContext").Return(iadapter.SpanContext(sc))
				span.On("End").Return()

				repo.On("FindUserByEmail", ctx, email).Return(nil, nil)
				log.On("InfoJSON", mock.Anything, mock.Anything).Return()

				userService := service.NewCustomerService(repo, tracer, log)

				err := userService.ValidateEmailAvailability(context.Background(), "test@test.com")

				Expect(err).To(BeNil())
				Expect(repo.AssertNumberOfCalls(GinkgoT(), "FindUserByEmail", 1)).To(BeTrue())
				Expect(tracer.AssertNumberOfCalls(GinkgoT(), "Start", 1)).To(BeTrue())
				Expect(span.AssertNumberOfCalls(GinkgoT(), "End", 1)).To(BeTrue())
				Expect(sc.AssertNumberOfCalls(GinkgoT(), "TraceID", 1)).To(BeTrue())
				Expect(log.AssertNumberOfCalls(GinkgoT(), "InfoJSON", 2)).To(BeTrue())
				Expect(log.AssertCalled(GinkgoT(), "InfoJSON", "Validating email availability", mock.MatchedBy(func(m map[string]any) bool {
					return m["trace_id"] == "trace-123" && m["email"] == email
				}))).To(BeTrue())
				Expect(log.AssertCalled(GinkgoT(), "InfoJSON", "Email is available", mock.MatchedBy(func(m map[string]any) bool {
					return m["trace_id"] == "trace-123" && m["email"] == email
				}))).To(BeTrue())
				Expect(repo.AssertCalled(GinkgoT(), "FindUserByEmail", ctx, email)).To(BeTrue())
				Expect(tracer.AssertCalled(GinkgoT(), "Start", ctx, "UserService.ValidateEmailAvailability")).To(BeTrue())
				Expect(span.AssertCalled(GinkgoT(), "SpanContext")).To(BeTrue())
				Expect(sc.AssertCalled(GinkgoT(), "TraceID")).To(BeTrue())

			})
		})

		Context("error cases", func() {
			It("should return an error when call userRepository", func() {
				ctx := context.Background()
				email := "test@test.com"

				repo := new(mrepository.UserRepositoryMock)
				tracer := new(madapters.TracerMock)
				span := new(madapters.SpanMock)
				sc := new(madapters.SpanContextMock)
				log := new(madapters.LoggerMock)

				tracer.On("Start", ctx, "UserService.ValidateEmailAvailability").
					Return(ctx, iadapter.Span(span))
				sc.On("TraceID").Return("trace-123")
				span.On("SpanContext").Return(iadapter.SpanContext(sc))
				span.On("End").Return()
				span.On("RecordError", mock.Anything).Return()

				repo.On("FindUserByEmail", ctx, email).Return(nil, errors.ErrorTransactionAlreadyExists())
				log.On("InfoJSON", mock.Anything, mock.Anything).Return()

				log.On("ErrorJSON", mock.Anything, mock.Anything).Return()
				log.On("ErrorJSON", mock.Anything, mock.Anything).Return()

				userService := service.NewCustomerService(repo, tracer, log)

				err := userService.ValidateEmailAvailability(context.Background(), "test@test.com")

				Expect(err).ToNot(BeNil())
				Expect(err).To(Equal(errors.ErrorTransactionAlreadyExists()))
				Expect(repo.AssertNumberOfCalls(GinkgoT(), "FindUserByEmail", 1)).To(BeTrue())
				Expect(tracer.AssertNumberOfCalls(GinkgoT(), "Start", 1)).To(BeTrue())
				Expect(span.AssertNumberOfCalls(GinkgoT(), "End", 1)).To(BeTrue())
				Expect(span.AssertNumberOfCalls(GinkgoT(), "RecordError", 1)).To(BeTrue())
				Expect(sc.AssertNumberOfCalls(GinkgoT(), "TraceID", 1)).To(BeTrue())
				Expect(log.AssertNumberOfCalls(GinkgoT(), "InfoJSON", 1)).To(BeTrue())
				Expect(log.AssertCalled(GinkgoT(), "InfoJSON", "Validating email availability", mock.MatchedBy(func(m map[string]any) bool {
					return m["trace_id"] == "trace-123" && m["email"] == email
				}))).To(BeTrue())
				Expect(log.AssertNumberOfCalls(GinkgoT(), "ErrorJSON", 1)).To(BeTrue())
				Expect(log.AssertCalled(GinkgoT(), "ErrorJSON", "Error finding user by email", mock.MatchedBy(func(m map[string]any) bool {
					return m["trace_id"] == "trace-123" && m["email"] == email && m["error"] == err.Error()
				}))).To(BeTrue())
				Expect(repo.AssertCalled(GinkgoT(), "FindUserByEmail", ctx, email)).To(BeTrue())
				Expect(tracer.AssertCalled(GinkgoT(), "Start", ctx, "UserService.ValidateEmailAvailability")).To(BeTrue())
				Expect(span.AssertCalled(GinkgoT(), "SpanContext")).To(BeTrue())
				Expect(sc.AssertCalled(GinkgoT(), "TraceID")).To(BeTrue())
			})

			It("should return an error when user with email already exists", func() {
				ctx := context.Background()
				email := "test@test.com"

				user := entity.BuilderUser().
					WithPublicID("123e4567-e89b-12d3-a456-426614174000").
					WithEmail("test@test.com").
					WithPassword("SenhaSegura537!").
					WithName("Test User").
					WithRole("admin").
					WithCreateAT(time.Now()).
					WithUpdateAT(time.Now()).
					Build()
				user.AssignID(1)
				user.AssignPasswordHash("hashedpassword123")

				repo := new(mrepository.UserRepositoryMock)
				tracer := new(madapters.TracerMock)
				span := new(madapters.SpanMock)
				sc := new(madapters.SpanContextMock)
				log := new(madapters.LoggerMock)

				tracer.On("Start", ctx, "UserService.ValidateEmailAvailability").
					Return(ctx, iadapter.Span(span))
				sc.On("TraceID").Return("trace-123")
				span.On("SpanContext").Return(iadapter.SpanContext(sc))
				span.On("End").Return()
				span.On("RecordError", mock.Anything).Return()
				repo.On("FindUserByEmail", ctx, email).Return(&user, nil)

				log.On("InfoJSON", mock.Anything, mock.Anything).Return()
				log.On("ErrorJSON", mock.Anything, mock.Anything).Return()

				userService := service.NewCustomerService(repo, tracer, log)

				err := userService.ValidateEmailAvailability(context.Background(), "test@test.com")

				Expect(err).ToNot(BeNil())
				Expect(err).To(Equal(errors.ErrorAlreadyExists(user.PublicID())))
				Expect(repo.AssertNumberOfCalls(GinkgoT(), "FindUserByEmail", 1)).To(BeTrue())
				Expect(tracer.AssertNumberOfCalls(GinkgoT(), "Start", 1)).To(BeTrue())
				Expect(span.AssertNumberOfCalls(GinkgoT(), "End", 1)).To(BeTrue())
				Expect(span.AssertNumberOfCalls(GinkgoT(), "RecordError", 1)).To(BeTrue())
				Expect(sc.AssertNumberOfCalls(GinkgoT(), "TraceID", 1)).To(BeTrue())
				Expect(log.AssertNumberOfCalls(GinkgoT(), "InfoJSON", 1)).To(BeTrue())
				Expect(log.AssertCalled(GinkgoT(), "InfoJSON", "Validating email availability", mock.MatchedBy(func(m map[string]any) bool {
					return m["trace_id"] == "trace-123" && m["email"] == email
				}))).To(BeTrue())
				Expect(log.AssertNumberOfCalls(GinkgoT(), "ErrorJSON", 1)).To(BeTrue())
				Expect(log.AssertCalled(GinkgoT(), "ErrorJSON", "Email already exists", mock.MatchedBy(func(m map[string]any) bool {
					return m["trace_id"] == "trace-123" && m["email"] == email && m["public_id"] == user.PublicID()
				}))).To(BeTrue())
				Expect(repo.AssertCalled(GinkgoT(), "FindUserByEmail", ctx, email)).To(BeTrue())
				Expect(tracer.AssertCalled(GinkgoT(), "Start", ctx, "UserService.ValidateEmailAvailability")).To(BeTrue())
				Expect(span.AssertCalled(GinkgoT(), "SpanContext")).To(BeTrue())
				Expect(sc.AssertCalled(GinkgoT(), "TraceID")).To(BeTrue())
			})
		})
	})
})
