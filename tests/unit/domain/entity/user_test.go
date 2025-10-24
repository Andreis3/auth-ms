//go:build unit

package entity_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/andreis3/auth-ms/internal/domain/entity"
	"github.com/andreis3/auth-ms/internal/domain/validator"
)

var _ = Describe("INTERNAL :: DOMAIN :: ENTITY :: USER", func() {
	Describe("#Validate", func() {
		Context("success cases", func() {
			It("should not return an error when build new user", func() {
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

				err := user.Validate()

				Expect(err.Errors()).To(BeEmpty())
				Expect(user.PublicID()).To(Equal("123e4567-e89b-12d3-a456-426614174000"))
				Expect(user.Email()).To(Equal("test@test.com"))
				Expect(user.PasswordHash()).To(Equal("hashedpassword123"))
				Expect(user.Name()).To(Equal("Test User"))
				Expect(user.Role()).To(Equal("admin"))
				Expect(user.CreateAT()).To(BeTemporally("~", time.Now(), 1*time.Second))
				Expect(user.UpdateAT()).To(BeTemporally("~", time.Now(), 1*time.Second))
			})
		})

		Context("error cases", func() {
			var (
				validationErr  *validator.Validator
				expectedErrors []string
			)

			BeforeEach(func() {
				user := entity.BuilderUser().Build()

				validationErr = user.Validate()
				expectedErrors = []string{
					fmt.Sprintf("email: %s", validator.ErrNotBlank),
					"email: invalid email format",
					fmt.Sprintf("name: %s", validator.ErrNotBlank),
					fmt.Sprintf("password: %s", validator.ErrNotBlank),
					"password: must be at least 8 characters",
					"password: must contain at least one uppercase letter",
					"password: must contain at least one lowercase letter",
					"password: must contain at least one number",
					"password: must contain at least one special character",
					fmt.Sprintf("public_id: %s", validator.ErrNotBlank),
					fmt.Sprintf("role: %s", validator.ErrNotBlank),
				}
			})

			It("should return all expected validation errors when user is empty", func() {
				errs := validationErr.Errors()

				Expect(errs).NotTo(BeNil())
				Expect(errs).To(HaveLen(len(expectedErrors)))
				for _, expected := range expectedErrors {
					Expect(errs).To(ContainElement(expected))
				}
			})

			It("should not include unexpected validation errors when user is empty", func() {
				Expect(validationErr.Errors()).To(Equal(expectedErrors))
			})
		})
	})
})
