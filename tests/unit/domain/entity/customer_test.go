//go:build unit
// +build unit

package entity_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/andreis3/customers-ms/internal/domain/entity"
	"github.com/andreis3/customers-ms/internal/domain/validator"
)

var _ = Describe("INTERNAL :: DOMAIN :: ENTITY :: CUSTOMER", func() {
	Describe("#Validate", func() {
		Context("success cases", func() {
			It("should not return an error when build new customer", func() {
				entity := entity.BuilderCustomer().
					WithID(123).
					WithEmail("any_email").
					WithPassword("any_password").
					WithFirstName("any_first_name").
					WithLastName("any_last_name").
					WithCPF("any_cpf").
					WithDateOfBirth(time.Now()).
					WithCreatedAt(time.Now()).
					WithUpdatedAt(time.Now()).
					Build()

				validate := entity.Validate()

				Expect(validate.Errors()).To(BeEmpty())
			})
		})

		Context("error cases", func() {
			It("should return an error when customer is empty", func() {
				entity := entity.BuilderCustomer().Build()

				validate := entity.Validate()

				Expect(validate.Errors()).NotTo(BeNil())
				Expect(validate.Errors()).To(HaveLen(4))
				Expect(validate.Errors()).To(ContainElement(fmt.Sprintf("email: %s", validator.ErrNotBlank)))
				Expect(validate.Errors()).To(ContainElement(fmt.Sprintf("first_name: %s", validator.ErrNotBlank)))
			})
		})
	})
})
