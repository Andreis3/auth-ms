package dto

type GetCustomerAddressesInput struct {
	CustomerID int64  `json:"customer_id"`
	Email      string `json:"email"`
}

type GetCustomerAddressesOutput struct {
	ID         *int64  `json:"id"`
	Street     *string `json:"street"`
	Number     *string `json:"number"`
	Complement *string `json:"complement"`
	City       *string `json:"city"`
	State      *string `json:"state"`
	PostalCode *string `json:"postal_code"`
	Country    *string `json:"country"`
}
