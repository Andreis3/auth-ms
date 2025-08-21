package mapper

import (
	"github.com/andreis3/auth-ms/internal/auth/application/dto"
	"github.com/andreis3/auth-ms/internal/auth/domain/entity"
)

func ToUser(input dto.CreateAuthUserInput) entity.User {
	return entity.BuilderUser().
		WithName(input.Name).
		WithPassword(input.Password).
		WithEmail(input.Email).
		Build()
}

func ToCreateAuthUserOutput(user *entity.User) *dto.CreateAuthUserOutput {
	return &dto.CreateAuthUserOutput{
		ID:       user.ID(),
		PublicID: user.PublicID(),
		Name:     user.Name(),
		Email:    user.Email(),
	}
}
