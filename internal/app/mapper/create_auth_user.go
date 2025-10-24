package mapper

import (
	"github.com/andreis3/auth-ms/internal/app/dto"
	"github.com/andreis3/auth-ms/internal/domain/entity"
)

func ToUser(input dto.CreateAuthUserInput) entity.User {
	return entity.BuilderUser().
		WithName(input.Name).
		WithPassword(input.Password).
		WithEmail(input.Email).
		Build()
}

func ToCreateAuthUserOutput(user *entity.User) *dto.CreateAuthUserOutput {
	const layout = "2006-01-02T15:04:05.000000Z"
	return &dto.CreateAuthUserOutput{
		PublicID:  user.PublicID(),
		Name:      user.Name(),
		Email:     user.Email(),
		Role:      user.Role(),
		CreatedAt: user.CreateAT().Format(layout),
	}
}
