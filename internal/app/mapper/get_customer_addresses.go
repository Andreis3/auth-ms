package mapper

import (
	"github.com/andreis3/customers-ms/internal/app/dto"
	"github.com/andreis3/customers-ms/internal/domain/entity"
	"github.com/andreis3/customers-ms/internal/util"
)

func ToGetCustomerAddressesOutput(addresses *[]entity.Address) *[]dto.GetCustomerAddressesOutput {
	var getCustomerAddressesOutput []dto.GetCustomerAddressesOutput
	for _, address := range *addresses {
		output := dto.GetCustomerAddressesOutput{
			ID:         util.ToInt64Pointer(address.ID),
			City:       util.ToStringPointer(address.City),
			Street:     util.ToStringPointer(address.Street),
			Number:     util.ToStringPointer(address.Number),
			Complement: util.ToStringPointer(address.Complement),
			Country:    util.ToStringPointer(address.Country),
			State:      util.ToStringPointer(address.State),
			PostalCode: util.ToStringPointer(address.PostalCode),
		}
		getCustomerAddressesOutput = append(getCustomerAddressesOutput, output)

	}

	return &getCustomerAddressesOutput
}
