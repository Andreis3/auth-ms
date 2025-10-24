//go:build unit

package command_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/andreis3/auth-ms/internal/app/dto"
	"github.com/andreis3/auth-ms/internal/app/mapper"
	"github.com/andreis3/auth-ms/internal/domain/entity"
	"github.com/andreis3/auth-ms/internal/domain/errors"
	"github.com/andreis3/auth-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/auth-ms/internal/domain/validator"
	"github.com/andreis3/auth-ms/internal/infra/logger"
	"github.com/andreis3/auth-ms/tests/suts"
)

var _ = Describe("INTERNAL :: APP :: COMMAND :: CREATE_AUTH_USER", func() {
	Describe("#Execute", func() {
		Context("success cases", func() {
			It("should mask password fields when logging input body", func() {
				ctx := context.Background()
				input := dto.CreateAuthUserInput{
					Email:           "user@example.com",
					Password:        "Sup3r$ecretZ",
					PasswordConfirm: "Sup3r$ecretZ",
					Name:            "Test User",
				}

				sut := suts.MakeCreateAuthUserSut()

				sut.Tracer.On("Start", ctx, "CreateAuthUser.Execute").Return(ctx, adapter.Span(sut.Span))
				sut.Span.On("SpanContext").Return(adapter.SpanContext(sut.Sc))
				sut.Span.On("End").Return()
				sut.Sc.On("TraceID").Return("trace-123")

				sut.Log.On("InfoJSON", mock.Anything, mock.Anything).Return()

				sut.Utils.On("UUID").Return("generated-uuid")

				sut.Service.On("ValidateEmailAvailability", ctx, input.Email).Return((*errors.Error)(nil))

				sut.Bcrypt.On("Hash", input.Password).Return("hashed-password", (*errors.Error)(nil))

				createdUser := mapper.ToUser(input)
				createdUser.AssignPublicID("generated-uuid")
				createdUser.AssignRole(entity.RoleUser)
				createdUser.AssignPasswordHash("hashed-password")
				createdUser.AssignID(1)

				sut.Repo.On("CreateUser", ctx, mock.MatchedBy(func(user entity.User) bool {
					return user.PublicID() == "generated-uuid" &&
						user.Password() == input.Password &&
						user.PasswordHash() == "hashed-password" &&
						user.Email() == input.Email &&
						user.Name() == input.Name
				})).Return(&createdUser, (*errors.Error)(nil))

				command := sut.Build()

				output, err := command.Execute(ctx, input)

				Expect(err).To(BeNil())
				Expect(output).ToNot(BeNil())
				Expect(output.PublicID).To(Equal("generated-uuid"))

				Expect(sut.Log.AssertNumberOfCalls(GinkgoT(), "InfoJSON", 1)).To(BeTrue())
				Expect(sut.Log.AssertCalled(GinkgoT(), "InfoJSON", "Creating auth user", mock.MatchedBy(func(m map[string]any) bool {
					traceID, ok := m["trace_id"].(string)
					if !ok || traceID != "trace-123" {
						return false
					}

					body, ok := m["body"].(dto.CreateAuthUserInput)
					if !ok {
						return false
					}

					return body.Email == input.Email &&
						body.Name == input.Name &&
						body.Password == logger.Mask &&
						body.PasswordConfirm == logger.Mask
				}))).To(BeTrue())

				Expect(sut.Tracer.AssertCalled(GinkgoT(), "Start", ctx, "CreateAuthUser.Execute")).To(BeTrue())
				Expect(sut.Span.AssertCalled(GinkgoT(), "SpanContext")).To(BeTrue())
				Expect(sut.Sc.AssertCalled(GinkgoT(), "TraceID")).To(BeTrue())
				Expect(sut.Span.AssertCalled(GinkgoT(), "End")).To(BeTrue())
				Expect(sut.Utils.AssertCalled(GinkgoT(), "UUID")).To(BeTrue())
				Expect(sut.Service.AssertCalled(GinkgoT(), "ValidateEmailAvailability", ctx, input.Email)).To(BeTrue())
				Expect(sut.Bcrypt.AssertCalled(GinkgoT(), "Hash", input.Password)).To(BeTrue())
				Expect(sut.Repo.AssertCalled(GinkgoT(), "CreateUser", ctx, mock.AnythingOfType("entity.User"))).To(BeTrue())
			})

			Context("error cases", func() {
				It("should record span error when user validation fails", func() {
					ctx := context.Background()
					input := dto.CreateAuthUserInput{
						Email:           "user@example.com",
						Password:        "Sup3r$ecretZ",
						PasswordConfirm: "Sup3r$ecretZ",
						Name:            "",
					}

					sut := suts.MakeCreateAuthUserSut()

					sut.Tracer.On("Start", ctx, "CreateAuthUser.Execute").Return(ctx, adapter.Span(sut.Span))
					sut.Span.On("SpanContext").Return(adapter.SpanContext(sut.Sc))
					sut.Span.On("End").Return()
					sut.Sc.On("TraceID").Return("trace-123")

					sut.Log.On("InfoJSON", mock.Anything, mock.Anything).Return()

					sut.Utils.On("UUID").Return("generated-uuid")

					sut.Span.On("RecordError", mock.MatchedBy(func(err error) bool {
						domainErr, ok := err.(*errors.Error)
						if !ok {
							return false
						}
						_, hasName := domainErr.Fields["name"]
						return domainErr.Code == errors.ValidationCode && hasName
					})).Return()

					sut.Log.On("CriticalJSON", "User validation failed", mock.MatchedBy(func(m map[string]any) bool {
						traceID, ok := m["trace_id"].(string)
						if !ok || traceID != "trace-123" {
							return false
						}

						errs, ok := m["errors"].(map[string]any)
						if !ok {
							return false
						}

						value, ok := errs["name"].(string)
						return ok && value == validator.ErrNotBlank
					})).Return()

					command := sut.Build()

					output, err := command.Execute(ctx, input)

					Expect(output).To(BeNil())
					Expect(err).ToNot(BeNil())
					Expect(err.Code).To(Equal(errors.ValidationCode))
					Expect(err.Fields).To(HaveKeyWithValue("name", validator.ErrNotBlank))

					Expect(sut.Span.AssertCalled(GinkgoT(), "RecordError", err)).To(BeTrue())
					Expect(sut.Service.AssertNotCalled(GinkgoT(), "ValidateEmailAvailability", mock.Anything, mock.Anything)).To(BeTrue())
					Expect(sut.Bcrypt.AssertNotCalled(GinkgoT(), "Hash", mock.Anything)).To(BeTrue())
					Expect(sut.Repo.AssertNotCalled(GinkgoT(), "CreateUser", mock.Anything, mock.Anything)).To(BeTrue())
				})

				It("should record span error when email validation fails", func() {
					ctx := context.Background()
					input := dto.CreateAuthUserInput{
						Email:           "user@example.com",
						Password:        "Sup3r$ecretZ",
						PasswordConfirm: "Sup3r$ecretZ",
						Name:            "Test User",
					}

					sut := suts.MakeCreateAuthUserSut()

					sut.Tracer.On("Start", ctx, "CreateAuthUser.Execute").Return(ctx, adapter.Span(sut.Span))
					sut.Span.On("SpanContext").Return(adapter.SpanContext(sut.Sc))
					sut.Span.On("End").Return()
					sut.Sc.On("TraceID").Return("trace-123")

					sut.Log.On("InfoJSON", mock.Anything, mock.Anything).Return()

					sut.Utils.On("UUID").Return("generated-uuid")

					validationErr := errors.ErrorTransactionAlreadyExists()
					sut.Service.On("ValidateEmailAvailability", ctx, input.Email).Return(validationErr)

					sut.Span.On("RecordError", validationErr).Return()

					sut.Log.On("ErrorJSON", "Email validation failed", mock.MatchedBy(func(m map[string]any) bool {
						return m["trace_id"] == "trace-123" &&
							m["email"] == input.Email &&
							m["error"] == validationErr.Error()
					})).Return()

					command := sut.Build()

					output, err := command.Execute(ctx, input)

					Expect(err).To(Equal(validationErr))
					Expect(output).To(BeNil())

					Expect(sut.Span.AssertCalled(GinkgoT(), "RecordError", validationErr)).To(BeTrue())
					Expect(sut.Bcrypt.AssertNotCalled(GinkgoT(), "Hash", mock.Anything)).To(BeTrue())
					Expect(sut.Repo.AssertNotCalled(GinkgoT(), "CreateUser", mock.Anything, mock.Anything)).To(BeTrue())
				})

				It("should record span error when hashing password fails", func() {
					ctx := context.Background()
					input := dto.CreateAuthUserInput{
						Email:           "user@example.com",
						Password:        "Sup3r$ecretZ",
						PasswordConfirm: "Sup3r$ecretZ",
						Name:            "Test User",
					}

					sut := suts.MakeCreateAuthUserSut()

					sut.Tracer.On("Start", ctx, "CreateAuthUser.Execute").Return(ctx, adapter.Span(sut.Span))
					sut.Span.On("SpanContext").Return(adapter.SpanContext(sut.Sc))
					sut.Span.On("End").Return()
					sut.Sc.On("TraceID").Return("trace-123")

					sut.Log.On("InfoJSON", mock.Anything, mock.Anything).Return()

					sut.Utils.On("UUID").Return("generated-uuid")

					sut.Service.On("ValidateEmailAvailability", ctx, input.Email).Return((*errors.Error)(nil))

					hashErr := errors.ErrorHashPassword(assert.AnError)
					sut.Bcrypt.On("Hash", input.Password).Return("", hashErr)

					sut.Span.On("RecordError", hashErr).Return()

					sut.Log.On("CriticalJSON", "Error hashing password", mock.MatchedBy(func(m map[string]any) bool {
						return m["trace_id"] == "trace-123" &&
							m["error"] == hashErr.Error()
					})).Return()

					command := sut.Build()

					output, err := command.Execute(ctx, input)

					Expect(err).To(Equal(hashErr))
					Expect(output).To(BeNil())

					Expect(sut.Span.AssertCalled(GinkgoT(), "RecordError", hashErr)).To(BeTrue())
					Expect(sut.Repo.AssertNotCalled(GinkgoT(), "CreateUser", mock.Anything, mock.Anything)).To(BeTrue())
				})

				It("should record span error when repository create fails", func() {
					ctx := context.Background()
					input := dto.CreateAuthUserInput{
						Email:           "user@example.com",
						Password:        "Sup3r$ecretZ",
						PasswordConfirm: "Sup3r$ecretZ",
						Name:            "Test User",
					}

					sut := suts.MakeCreateAuthUserSut()

					sut.Tracer.On("Start", ctx, "CreateAuthUser.Execute").Return(ctx, adapter.Span(sut.Span))
					sut.Span.On("SpanContext").Return(adapter.SpanContext(sut.Sc))
					sut.Span.On("End").Return()
					sut.Sc.On("TraceID").Return("trace-123")

					sut.Log.On("InfoJSON", mock.Anything, mock.Anything).Return()

					sut.Utils.On("UUID").Return("generated-uuid")

					sut.Service.On("ValidateEmailAvailability", ctx, input.Email).Return((*errors.Error)(nil))

					sut.Bcrypt.On("Hash", input.Password).Return("hashed-password", (*errors.Error)(nil))

					repoErr := errors.New(errors.ErrInternal, "repository error")
					sut.Repo.On("CreateUser", ctx, mock.AnythingOfType("entity.User")).Return((*entity.User)(nil), repoErr)

					sut.Span.On("RecordError", repoErr).Return()

					sut.Log.On("ErrorJSON", "Error creating user", mock.MatchedBy(func(m map[string]any) bool {
						return m["trace_id"] == "trace-123" &&
							m["error"] == repoErr.Error()
					})).Return()

					command := sut.Build()

					output, err := command.Execute(ctx, input)

					Expect(err).To(Equal(repoErr))
					Expect(output).To(BeNil())

					Expect(sut.Span.AssertCalled(GinkgoT(), "RecordError", repoErr)).To(BeTrue())
				})
			})
		})
	})
})
