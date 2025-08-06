package model

import (
	"time"

	"github.com/andreis3/customers-ms/internal/domain/entity"
	"github.com/andreis3/customers-ms/internal/util"
)

type Customer struct {
	ID          *int64
	Email       *string
	Password    *string
	FirstName   *string
	LastName    *string
	CPF         *string
	DateOfBirth *time.Time
	CreatedAT   *time.Time
	UpdatedAT   *time.Time
}

func (c Customer) ToEntity() entity.Customer {
	return entity.BuilderCustomer().
		WithID(util.ToInt64(c.ID)).
		WithEmail(util.ToString(c.Email)).
		WithPassword(util.ToString(c.Password)).
		WithFirstName(util.ToString(c.FirstName)).
		WithLastName(util.ToString(c.LastName)).
		WithCPF(util.ToString(c.CPF)).
		WithDateOfBirth(util.ToTime(c.DateOfBirth)).
		WithCreatedAt(util.ToTime(c.CreatedAT)).
		WithUpdatedAt(util.ToTime(c.UpdatedAT)).
		Build()
}

func (c Customer) FromModel(customer entity.Customer) *Customer {
	dateNow := time.Now()
	return &Customer{
		Email:       util.ToStringPointer(customer.Email.String()),
		Password:    util.ToStringPointer(customer.Password.String()),
		FirstName:   util.ToStringPointer(customer.FirstName),
		LastName:    util.ToStringPointer(customer.LastName),
		CPF:         util.ToStringPointer(customer.CPF.String()),
		DateOfBirth: util.ToTimePointer(customer.DateOfBirth),
		CreatedAT:   util.ToTimePointer(dateNow),
		UpdatedAT:   util.ToTimePointer(dateNow),
	}
}
