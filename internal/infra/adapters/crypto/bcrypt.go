package crypto

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/andreis3/customers-ms/internal/domain/errors"
)

type Bcrypt struct{}

func NewBcrypt() *Bcrypt {
	return &Bcrypt{}
}

func (b *Bcrypt) Hash(data string) (string, *errors.Error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(data), 5)
	if err != nil {
		return "", errors.ErrorHashPassword(err)
	}
	return string(bytes), nil
}

func (b *Bcrypt) CompareHash(data string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(data)) == nil
}
