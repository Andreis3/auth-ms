package model

import (
	"time"

	"github.com/andreis3/auth-ms/internal/domain/entity"
	"github.com/andreis3/auth-ms/internal/util"
)

type User struct {
	ID        *int64     `db:"id"`
	PublicID  *string    `db:"public_id"`
	Email     *string    `db:"email"`
	Password  *string    `db:"password"`
	Name      *string    `db:"name"`
	Role      *string    `db:"role"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) ToEntity() entity.User {
	var roleType entity.RoleTypes
	if u.Role != nil {
		roleType = entity.RoleTypes(*u.Role)
	}
	return entity.BuilderUser().
		WithID(util.ToInt64(u.ID)).
		WithPublicID(util.ToString(u.PublicID)).
		WithEmail(util.ToString(u.Email)).
		WithPassword(util.ToString(u.Password)).
		WithName(util.ToString(u.Name)).
		WithRole(roleType).
		WithCreateAT(util.ToTime(u.CreatedAt)).
		WithUpdateAT(util.ToTime(u.UpdatedAt)).
		WithDeletedAt(u.DeletedAt).
		Build()
}

func (u *User) ToModel(user entity.User) *User {
	dateNow := time.Now().UTC()
	return &User{
		PublicID:  util.ToStringPointer(user.PublicID()),
		Email:     util.ToStringPointer(user.Email()),
		Password:  util.ToStringPointer(user.PasswordHash()),
		Name:      util.ToStringPointer(user.Name()),
		Role:      util.ToStringPointer(user.Role()),
		CreatedAt: util.ToTimePointer(dateNow),
		UpdatedAt: util.ToTimePointer(dateNow),
	}
}
