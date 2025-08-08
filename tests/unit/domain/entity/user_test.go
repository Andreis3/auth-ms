//go:build unit

package entity_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/andreis3/customers-ms/internal/domain/entity"
	"github.com/andreis3/customers-ms/internal/domain/validator"
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
			It("should return an error when user is empty", func() {
				user := entity.BuilderUser().Build()

				err := user.Validate()
				Expect(err.Errors()).NotTo(BeNil())
				Expect(err.Errors()).To(HaveLen(11))
				Expect(err.Errors()).To(ContainElement(fmt.Sprintf("public_id: %s", validator.ErrNotBlank)))
				Expect(err.Errors()).To(ContainElement(fmt.Sprintf("email: %s", validator.ErrNotBlank)))
				Expect(err.Errors()).To(ContainElement(fmt.Sprintf("password: %s", validator.ErrNotBlank)))
				Expect(err.Errors()).To(ContainElement(fmt.Sprintf("name: %s", validator.ErrNotBlank)))
				Expect(err.Errors()).To(ContainElement(fmt.Sprintf("role: %s", validator.ErrNotBlank)))
			})
		})
	})
})
