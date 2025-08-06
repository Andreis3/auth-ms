package dto

type CreateCustomerInput struct {
	Email       *string               `json:"email"`
	Password    *string               `json:"password"`
	FirstName   *string               `json:"first_name"`
	LastName    *string               `json:"last_name"`
	CPF         *string               `json:"cpf"`
	DateOfBirth *string               `json:"date_of_birth"` // depois você pode converter para time.Time
	Addresses   *[]CreateAddressInput `json:"addresses"`
}

type CreateAddressInput struct {
	Street     *string `json:"street"`
	Number     *string `json:"number"`
	Complement *string `json:"complement"`
	City       *string `json:"city"`
	State      *string `json:"state"`
	PostalCode *string `json:"postal_code"`
	Country    *string `json:"country"`
}

type CreatedCustomerOutput struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	DateOfBirth string `json:"date_of_birth"`
}
