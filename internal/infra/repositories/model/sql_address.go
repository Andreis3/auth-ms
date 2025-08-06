package model

import (
	"time"

	"github.com/andreis3/customers-ms/internal/domain/entity"
	"github.com/andreis3/customers-ms/internal/util"
)

type Address struct {
	ID         *int64
	CustomerID *int64
	Street     *string
	Number     *string
	Complement *string
	City       *string
	State      *string
	PostalCode *string
	Country    *string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

func (a Address) ToEntity() entity.Address {

	return entity.BuilderAddress().
		WithID(util.ToInt64(a.ID)).
		WithCustomerID(util.ToInt64(a.CustomerID)).
		WithStreet(util.ToString(a.Street)).
		WithNumber(util.ToString(a.Number)).
		WithComplement(util.ToString(a.Complement)).
		WithCity(util.ToString(a.City)).
		WithState(util.ToString(a.State)).
		WithPostalCode(util.ToString(a.PostalCode)).
		WithCountry(util.ToString(a.Country)).
		WithCreatedAt(util.ToTime(a.CreatedAt)).
		WithUpdatedAt(util.ToTime(a.UpdatedAt)).
		Build()
}

func (a Address) FromModel(entity entity.Address) *Address {
	dateNow := time.Now()
	return &Address{
		CustomerID: util.ToInt64Pointer(entity.CustomerID),
		Street:     util.ToStringPointer(entity.Street),
		Number:     util.ToStringPointer(entity.Number),
		Complement: util.ToStringPointer(entity.Complement),
		City:       util.ToStringPointer(entity.City),
		State:      util.ToStringPointer(entity.State),
		PostalCode: util.ToStringPointer(entity.PostalCode),
		Country:    util.ToStringPointer(entity.Country),
		CreatedAt:  util.ToTimePointer(dateNow),
		UpdatedAt:  util.ToTimePointer(dateNow),
	}
}
