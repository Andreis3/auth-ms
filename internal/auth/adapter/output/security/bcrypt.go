package security

import (
	"golang.org/x/crypto/bcrypt"

	errors2 "github.com/andreis3/auth-ms/internal/auth/domain/errors"
)

type Bcrypt struct{}

func NewBcrypt() *Bcrypt {
	return &Bcrypt{}
}

func (b *Bcrypt) Hash(data string) (string, *errors2.Error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(data), 5)
	if err != nil {
		return "", errors2.ErrorHashPassword(err)
	}
	return string(bytes), nil
}

func (b *Bcrypt) CompareHash(data string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(data)) == nil
}
