package services

import (
	"context"

	"github.com/andreis3/customers-ms/internal/domain/interfaces/postgres"
)

type CustomerService struct {
	customerRepository postgres.CustomerRepository
}

func NewCustomerService(customerRepository postgres.CustomerRepository) *CustomerService {
	return &CustomerService{
		customerRepository: customerRepository,
	}
}

func (s *CustomerService) ExistCustomerByEmail(ctx context.Context, email string) bool {
	customer, err := s.customerRepository.FindCustomerByEmail(ctx, email)
	return err == nil && customer != nil
}
