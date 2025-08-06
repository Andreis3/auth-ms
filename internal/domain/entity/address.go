package entity

import (
	"time"

	"github.com/andreis3/customers-ms/internal/domain/validator"
)

type Address struct {
	ID         int64
	CustomerID int64
	Street     string
	Number     string
	Complement string
	City       string
	State      string
	PostalCode string
	Country    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func BuilderAddress() *Address {
	return &Address{}
}

func (a *Address) AssignID(id int64) {
	a.ID = id
}

func (a *Address) Validate() *validator.Validator {
	v := validator.New()
	v.Assert(validator.NotBlank(a.Street), "Street", validator.ErrNotBlank)
	v.Assert(validator.NotBlank(a.Number), "Number", validator.ErrNotBlank)
	v.Assert(validator.NotBlank(a.City), "City", validator.ErrNotBlank)
	v.Assert(validator.NotBlank(a.State), "State", validator.ErrNotBlank)
	v.Assert(validator.NotBlank(a.PostalCode), "postal_code", validator.ErrNotBlank)
	v.Assert(validator.NotBlank(a.Country), "Country", validator.ErrNotBlank)
	return v
}

func (a *Address) WithID(id int64) *Address {
	a.ID = id
	return a
}

func (a *Address) WithCustomerID(customerID int64) *Address {
	a.CustomerID = customerID
	return a
}

func (a *Address) WithStreet(street string) *Address {
	a.Street = street
	return a
}

func (a *Address) WithNumber(number string) *Address {
	a.Number = number
	return a
}

func (a *Address) WithComplement(complement string) *Address {
	a.Complement = complement
	return a
}

func (a *Address) WithCity(city string) *Address {
	a.City = city
	return a
}

func (a *Address) WithState(state string) *Address {
	a.State = state
	return a
}

func (a *Address) WithPostalCode(postalCode string) *Address {
	a.PostalCode = postalCode
	return a
}

func (a *Address) WithCountry(country string) *Address {
	a.Country = country
	return a
}

func (a *Address) WithCreatedAt(createdAt time.Time) *Address {
	a.CreatedAt = createdAt
	return a
}

func (a *Address) WithUpdatedAt(updatedAt time.Time) *Address {
	a.UpdatedAt = updatedAt
	return a
}

func (a *Address) Build() Address {
	return *a
}
