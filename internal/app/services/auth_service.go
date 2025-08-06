package services

import (
	"context"

	"github.com/andreis3/customers-ms/internal/domain/entity"
	"github.com/andreis3/customers-ms/internal/domain/errors"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/postgres"
	"github.com/andreis3/customers-ms/internal/domain/valueobject"
)

type AuthService struct {
	tokenService       adapter.TokenService
	customerRepository postgres.CustomerRepository
}

func NewAuthService(
	tokenService adapter.TokenService,
	customerRepository postgres.CustomerRepository,
) *AuthService {
	return &AuthService{
		tokenService:       tokenService,
		customerRepository: customerRepository,
	}
}

func (a *AuthService) GenerateToken(customer entity.Customer) (*valueobject.TokenClaims, *errors.Error) {
	return a.tokenService.CreateToken(customer)
}

func (a *AuthService) ValidateToken(tokenString string) (*valueobject.TokenClaims, *errors.Error) {
	return a.tokenService.ValidateToken(tokenString)
}

func (a *AuthService) RefreshToken(tokenString string) (*valueobject.TokenClaims, *errors.Error) {
	return a.tokenService.RefreshToken(tokenString)
}

func (a *AuthService) DecodeToken(ctx context.Context, tokenString string) (*valueobject.TokenClaims, *errors.Error) {
	tokenClaims, err := a.tokenService.DecodeToken(tokenString)
	if err != nil {
		return nil, err
	}

	customer, err := a.customerRepository.FindByID(ctx, tokenClaims.CustomerID)

	if err != nil {
		return nil, err
	}

	if customer == nil {
		return nil, errors.ErrorInvalidToken()
	}
	return tokenClaims, nil
}
