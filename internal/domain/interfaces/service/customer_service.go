package service

import "context"

type CustomerService interface {
	ExistCustomerByEmail(ctx context.Context, email string) bool
}
