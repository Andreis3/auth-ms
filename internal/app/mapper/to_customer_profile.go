package mapper

import (
	"time"

	"github.com/andreis3/customers-ms/internal/app/dto"
	"github.com/andreis3/customers-ms/internal/domain/aggregate"
	"github.com/andreis3/customers-ms/internal/domain/entity"
	"github.com/andreis3/customers-ms/internal/util"
)

func ToCustomerProfile(input dto.CreateCustomerInput) aggregate.CustomerProfile {
	dob, err := time.Parse("02-01-2006", util.ToString(input.DateOfBirth))
	if err != nil {
		return aggregate.CustomerProfile{}
	}

	customer := entity.BuilderCustomer().
		WithEmail(util.ToString(input.Email)).
		WithPassword(util.ToString(input.Password)).
		WithFirstName(util.ToString(input.FirstName)).
		WithLastName(util.ToString(input.LastName)).
		WithCPF(util.ToString(input.CPF)).
		WithDateOfBirth(dob).
		Build()

	var addresses []entity.Address
	for _, a := range *input.Addresses {
		addresses = append(addresses, entity.BuilderAddress().
			WithStreet(util.ToString(a.Street)).
			WithNumber(util.ToString(a.Number)).
			WithComplement(util.ToString(a.Complement)).
			WithCity(util.ToString(a.City)).
			WithState(util.ToString(a.State)).
			WithPostalCode(util.ToString(a.PostalCode)).
			WithCountry(util.ToString(a.Country)).
			Build())
	}

	return *aggregate.NewCustomerProfile(customer, addresses)
}

func CustomerOutput(customer entity.Customer) dto.CreatedCustomerOutput {
	return dto.CreatedCustomerOutput{
		ID:          customer.ID,
		Email:       customer.Email.String(),
		Name:        customer.FullName(),
		DateOfBirth: customer.DateOfBirth.Format("02-01-2006"),
	}
}
