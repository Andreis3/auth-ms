//go:build unit

package command_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/andreis3/auth-ms/internal/app/dto"
	"github.com/andreis3/auth-ms/internal/app/mapper"
	"github.com/andreis3/auth-ms/internal/domain/entity"
	"github.com/andreis3/auth-ms/internal/domain/errors"
	"github.com/andreis3/auth-ms/internal/domain/interfaces/adapter"
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
		})
	})
})
