package entity

import "C"
import (
	"fmt"
	"time"

	"github.com/andreis3/customers-ms/internal/domain/validator"
	"github.com/andreis3/customers-ms/internal/domain/valueobject"
)

type Customer struct {
	ID          int64
	Email       valueobject.Email
	Password    valueobject.Password
	FirstName   string
	LastName    string
	CPF         valueobject.CPF
	DateOfBirth time.Time
	Addresses   []Address
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func BuilderCustomer() *Customer {
	return &Customer{}
}

func (c *Customer) Build() Customer {
	return *c
}

func (c *Customer) AssignID(id int64) {
	c.ID = id
}

func (c *Customer) AssignHashedPassword(hashed string) {
	c.Password = valueobject.NewEncryptedPassword(hashed)
}

func (c *Customer) FullName() string {
	return fmt.Sprintf("%s %s", c.FirstName, c.LastName)
}

func (c *Customer) Validate() *validator.Validator {
	v := validator.New()
	v.Assert(validator.NotBlank(c.FirstName), "first_name", validator.ErrNotBlank)
	v.Assert(validator.NotBlank(c.LastName), "last_name", validator.ErrNotBlank)
	v.Assert(!c.DateOfBirth.IsZero(), "date_of_birth", validator.ErrNotBlank)

	v.Merge(c.Email.Validate())
	v.Merge(c.CPF.Validate())
	v.Merge(c.Password.Validate())
	return v
}

func (c *Customer) WithID(id int64) *Customer {
	c.ID = id
	return c
}

func (c *Customer) WithEmail(email string) *Customer {
	c.Email = valueobject.NewEmail(email)
	return c
}

func (c *Customer) WithPassword(password string) *Customer {
	c.Password = valueobject.NewPassword(password)
	return c
}

func (c *Customer) WithFirstName(firstName string) *Customer {
	c.FirstName = firstName
	return c
}

func (c *Customer) WithLastName(lastName string) *Customer {
	c.LastName = lastName
	return c
}

func (c *Customer) WithCPF(cpf string) *Customer {
	c.CPF = valueobject.NewCPF(cpf)
	return c
}

func (c *Customer) WithDateOfBirth(dateOfBirth time.Time) *Customer {
	c.DateOfBirth = dateOfBirth
	return c
}

func (c *Customer) WithAddresses(addresses []Address) *Customer {
	c.Addresses = addresses
	return c
}

func (c *Customer) WithCreatedAt(createdAt time.Time) *Customer {
	c.CreatedAt = createdAt
	return c
}

func (c *Customer) WithUpdatedAt(updatedAt time.Time) *Customer {
	c.UpdatedAt = updatedAt
	return c
}
