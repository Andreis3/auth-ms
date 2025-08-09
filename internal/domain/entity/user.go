package entity

import (
	"time"

	"github.com/andreis3/auth-ms/internal/domain/validator"
	"github.com/andreis3/auth-ms/internal/domain/valueobject"
)

type RoleTypes string

const (
	RoleUser RoleTypes = "user"
)

type User struct {
	id           int64
	publicID     string
	email        valueobject.Email
	password     valueobject.Password
	passwordHash string
	name         string
	role         RoleTypes
	createAT     time.Time
	updateAT     time.Time
	deletedAt    *time.Time
}

func BuilderUser() *User {
	return &User{}
}

func (u *User) Build() User {
	return *u
}

func (u *User) WithID(id int64) *User {
	u.id = id
	return u
}

func (u *User) WithPublicID(publicID string) *User {
	u.publicID = publicID
	return u
}

func (u *User) WithEmail(email string) *User {
	u.email = valueobject.NewEmail(email)
	return u
}

func (u *User) WithPassword(password string) *User {
	u.password = valueobject.NewPassword(password)
	return u
}

func (u *User) WithName(name string) *User {
	u.name = name
	return u
}

func (u *User) WithRole(role RoleTypes) *User {
	u.role = role
	return u
}

func (u *User) WithCreateAT(createAT time.Time) *User {
	u.createAT = createAT
	return u
}

func (u *User) WithUpdateAT(updateAT time.Time) *User {
	u.updateAT = updateAT
	return u
}

func (u *User) WithDeletedAt(deletedAt *time.Time) *User {
	u.deletedAt = deletedAt
	return u
}

func (u *User) Validate() *validator.Validator {
	v := validator.New()
	v.Assert(validator.NotBlank(u.name), "name", validator.ErrNotBlank)
	v.Assert(validator.NotBlank(string(u.role)), "role", validator.ErrNotBlank)
	v.Assert(validator.NotBlank(u.publicID), "public_id", validator.ErrNotBlank)
	v.Merge(u.email.Validate())
	v.Merge(u.password.Validate())
	return v
}

func (u *User) AssignID(id int64) *User {
	u.id = id
	return u
}

func (u *User) AssignPasswordHash(passwordHash string) *User {
	u.passwordHash = passwordHash
	return u
}

func (u *User) AssignPublicID(publicID string) *User {
	u.publicID = publicID
	return u
}

func (u *User) AssignRole(role RoleTypes) *User {
	u.role = role
	return u
}

func (u *User) ID() int64 {
	return u.id
}
func (u *User) PublicID() string {
	return u.publicID
}
func (u *User) Email() string {
	return u.email.String()
}
func (u *User) Password() string {
	return u.password.String()
}
func (u *User) PasswordHash() string {
	return u.passwordHash
}
func (u *User) Name() string {
	return u.name
}
func (u *User) Role() string {
	return string(u.role)
}
func (u *User) CreateAT() time.Time {
	return u.createAT
}
func (u *User) UpdateAT() time.Time {
	return u.updateAT
}
func (u *User) DeletedAt() *time.Time {
	return u.deletedAt
}
