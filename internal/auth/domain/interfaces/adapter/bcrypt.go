package adapter

import (
	"github.com/andreis3/auth-ms/internal/auth/domain/errors"
)

type Bcrypt interface {
	Hash(data string) (string, *errors.Error)
	CompareHash(hash, data string) bool
}
