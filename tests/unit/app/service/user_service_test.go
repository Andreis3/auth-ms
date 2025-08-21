//go:build unit

package service_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/andreis3/auth-ms/internal/auth/domain/entity"
	errors2 "github.com/andreis3/auth-ms/internal/auth/domain/errors"
	"github.com/andreis3/auth-ms/internal/auth/domain/interfaces/adapter"
	"github.com/andreis3/auth-ms/tests/suts"
)

var _ = Describe("INTERNAL :: APP :: SERVICE :: USER_SERVICE", func() {
	Describe("#Execute", func() {
		Context("success cases", func() {
			It("should not return an error when not exist user with email equal", func() {
				ctx := context.Background()
				email := "test@test.com"

				sut := suts.MakeUserServiceSut()

				sut.Tracer.On("Start", ctx, "UserService.ValidateEmailAvailability").
					Return(ctx, adapter.Span(sut.Span))
				sut.Sc.On("TraceID").Return("trace-123")
				sut.Span.On("SpanContext").Return(adapter.SpanContext(sut.Sc))
				sut.Span.On("End").Return()

				sut.Repo.On("FindUserByEmail", ctx, email).Return(nil, nil)
				sut.Log.On("InfoJSON", mock.Anything, mock.Anything).Return()

				userService := sut.Build()

				err := userService.ValidateEmailAvailability(context.Background(), "test@test.com")

				Expect(err).To(BeNil())
				Expect(sut.Repo.AssertNumberOfCalls(GinkgoT(), "FindUserByEmail", 1)).To(BeTrue())
				Expect(sut.Tracer.AssertNumberOfCalls(GinkgoT(), "Start", 1)).To(BeTrue())
				Expect(sut.Span.AssertNumberOfCalls(GinkgoT(), "End", 1)).To(BeTrue())
				Expect(sut.Sc.AssertNumberOfCalls(GinkgoT(), "TraceID", 1)).To(BeTrue())
				Expect(sut.Log.AssertNumberOfCalls(GinkgoT(), "InfoJSON", 2)).To(BeTrue())
				Expect(sut.Log.AssertCalled(GinkgoT(), "InfoJSON", "Validating email availability", mock.MatchedBy(func(m map[string]any) bool {
					return m["trace_id"] == "trace-123" && m["email"] == email
				}))).To(BeTrue())
				Expect(sut.Log.AssertCalled(GinkgoT(), "InfoJSON", "Email is available", mock.MatchedBy(func(m map[string]any) bool {
					return m["trace_id"] == "trace-123" && m["email"] == email
				}))).To(BeTrue())
				Expect(sut.Repo.AssertCalled(GinkgoT(), "FindUserByEmail", ctx, email)).To(BeTrue())
				Expect(sut.Tracer.AssertCalled(GinkgoT(), "Start", ctx, "UserService.ValidateEmailAvailability")).To(BeTrue())
				Expect(sut.Span.AssertCalled(GinkgoT(), "SpanContext")).To(BeTrue())
				Expect(sut.Sc.AssertCalled(GinkgoT(), "TraceID")).To(BeTrue())

			})
		})

		Context("error cases", func() {
			It("should return an error when call userRepository", func() {
				ctx := context.Background()
				email := "test@test.com"

				sut := suts.MakeUserServiceSut()

				sut.Tracer.On("Start", ctx, "UserService.ValidateEmailAvailability").
					Return(ctx, adapter.Span(sut.Span))
				sut.Sc.On("TraceID").Return("trace-123")
				sut.Span.On("SpanContext").Return(adapter.SpanContext(sut.Sc))
				sut.Span.On("End").Return()
				sut.Span.On("RecordError", mock.Anything).Return()

				sut.Repo.On("FindUserByEmail", ctx, email).Return(nil, errors2.ErrorTransactionAlreadyExists())
				sut.Log.On("InfoJSON", mock.Anything, mock.Anything).Return()

				sut.Log.On("ErrorJSON", mock.Anything, mock.Anything).Return()
				sut.Log.On("ErrorJSON", mock.Anything, mock.Anything).Return()

				userService := sut.Build()

				err := userService.ValidateEmailAvailability(context.Background(), "test@test.com")

				Expect(err).ToNot(BeNil())
				Expect(err).To(Equal(errors2.ErrorTransactionAlreadyExists()))
				Expect(sut.Repo.AssertNumberOfCalls(GinkgoT(), "FindUserByEmail", 1)).To(BeTrue())
				Expect(sut.Tracer.AssertNumberOfCalls(GinkgoT(), "Start", 1)).To(BeTrue())
				Expect(sut.Span.AssertNumberOfCalls(GinkgoT(), "End", 1)).To(BeTrue())
				Expect(sut.Span.AssertNumberOfCalls(GinkgoT(), "RecordError", 1)).To(BeTrue())
				Expect(sut.Sc.AssertNumberOfCalls(GinkgoT(), "TraceID", 1)).To(BeTrue())
				Expect(sut.Log.AssertNumberOfCalls(GinkgoT(), "InfoJSON", 1)).To(BeTrue())
				Expect(sut.Log.AssertCalled(GinkgoT(), "InfoJSON", "Validating email availability", mock.MatchedBy(func(m map[string]any) bool {
					return m["trace_id"] == "trace-123" && m["email"] == email
				}))).To(BeTrue())
				Expect(sut.Log.AssertNumberOfCalls(GinkgoT(), "ErrorJSON", 1)).To(BeTrue())
				Expect(sut.Log.AssertCalled(GinkgoT(), "ErrorJSON", "Error finding user by email", mock.MatchedBy(func(m map[string]any) bool {
					return m["trace_id"] == "trace-123" && m["email"] == email && m["error"] == err.Error()
				}))).To(BeTrue())
				Expect(sut.Repo.AssertCalled(GinkgoT(), "FindUserByEmail", ctx, email)).To(BeTrue())
				Expect(sut.Tracer.AssertCalled(GinkgoT(), "Start", ctx, "UserService.ValidateEmailAvailability")).To(BeTrue())
				Expect(sut.Span.AssertCalled(GinkgoT(), "SpanContext")).To(BeTrue())
				Expect(sut.Sc.AssertCalled(GinkgoT(), "TraceID")).To(BeTrue())
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

				//repo := new(mrepository.UserRepositoryMock)
				//tracer := new(madapters.TracerMock)
				//span := new(madapters.SpanMock)
				//sc := new(madapters.SpanContextMock)
				//log := new(madapters.LoggerMock)

				sut := suts.MakeUserServiceSut()

				sut.Tracer.On("Start", ctx, "UserService.ValidateEmailAvailability").
					Return(ctx, adapter.Span(sut.Span))
				sut.Sc.On("TraceID").Return("trace-123")
				sut.Span.On("SpanContext").Return(adapter.SpanContext(sut.Sc))
				sut.Span.On("End").Return()
				sut.Span.On("RecordError", mock.Anything).Return()
				sut.Repo.On("FindUserByEmail", ctx, email).Return(&user, nil)

				sut.Log.On("InfoJSON", mock.Anything, mock.Anything).Return()
				sut.Log.On("ErrorJSON", mock.Anything, mock.Anything).Return()

				userService := sut.Build()

				err := userService.ValidateEmailAvailability(context.Background(), "test@test.com")

				Expect(err).ToNot(BeNil())
				Expect(err).To(Equal(errors2.ErrorAlreadyExists(user.PublicID())))
				Expect(sut.Repo.AssertNumberOfCalls(GinkgoT(), "FindUserByEmail", 1)).To(BeTrue())
				Expect(sut.Tracer.AssertNumberOfCalls(GinkgoT(), "Start", 1)).To(BeTrue())
				Expect(sut.Span.AssertNumberOfCalls(GinkgoT(), "End", 1)).To(BeTrue())
				Expect(sut.Span.AssertNumberOfCalls(GinkgoT(), "RecordError", 1)).To(BeTrue())
				Expect(sut.Sc.AssertNumberOfCalls(GinkgoT(), "TraceID", 1)).To(BeTrue())
				Expect(sut.Log.AssertNumberOfCalls(GinkgoT(), "InfoJSON", 1)).To(BeTrue())
				Expect(sut.Log.AssertCalled(GinkgoT(), "InfoJSON", "Validating email availability", mock.MatchedBy(func(m map[string]any) bool {
					return m["trace_id"] == "trace-123" && m["email"] == email
				}))).To(BeTrue())
				Expect(sut.Log.AssertNumberOfCalls(GinkgoT(), "ErrorJSON", 1)).To(BeTrue())
				Expect(sut.Log.AssertCalled(GinkgoT(), "ErrorJSON", "Email already exists", mock.MatchedBy(func(m map[string]any) bool {
					return m["trace_id"] == "trace-123" && m["email"] == email && m["public_id"] == user.PublicID()
				}))).To(BeTrue())
				Expect(sut.Repo.AssertCalled(GinkgoT(), "FindUserByEmail", ctx, email)).To(BeTrue())
				Expect(sut.Tracer.AssertCalled(GinkgoT(), "Start", ctx, "UserService.ValidateEmailAvailability")).To(BeTrue())
				Expect(sut.Span.AssertCalled(GinkgoT(), "SpanContext")).To(BeTrue())
				Expect(sut.Sc.AssertCalled(GinkgoT(), "TraceID")).To(BeTrue())
			})
		})
	})
})
