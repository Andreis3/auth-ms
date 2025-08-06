//go:build unit
// +build unit

package entity_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/andreis3/customers-ms/internal/domain/entity/address"
	"github.com/andreis3/customers-ms/internal/domain/validator"
)

var _ = Describe("INTERNAL :: DOMAIN :: ENTITY :: ADDRESS", func() {
	Describe("#Validate", func() {
		Context("success cases", func() {
			It("should not return an error when build new address", func() {
				entity := address.NewBuilder().
					WithCustomerID(123).
					WithStreet("any_street").
					WithNumber("any_number").
					WithComplement("any_complement").
					WithCity("any_city").
					WithState("any_state").
					WithPostalCode("any_postal_code").
					WithCountry("any_country").
					WithCreatedAt(time.Now()).
					WithUpdatedAt(time.Now()).
					Build()

				validate := entity.Validate()

				Expect(validate.Errors()).To(BeEmpty())
			})
		})

		Context("error cases", func() {
			It("should return an error when address is empty", func() {
				entity := address.NewBuilder().Build()

				validate := entity.Validate()

				Expect(validate.Errors()).NotTo(BeNil())
				Expect(validate.Errors()).To(HaveLen(6))
				Expect(validate.Errors()).To(ContainElement(fmt.Sprintf("country: %s", validator.ErrNotBlank)))
				Expect(validate.Errors()).To(ContainElement(fmt.Sprintf("state: %s", validator.ErrNotBlank)))
			})
		})
	})
})
