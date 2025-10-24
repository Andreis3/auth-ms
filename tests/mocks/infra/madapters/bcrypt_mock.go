package madapters

import (
	"github.com/stretchr/testify/mock"

	"github.com/andreis3/auth-ms/internal/domain/errors"
)

type BcryptMock struct{ mock.Mock }

func (b *BcryptMock) Hash(data string) (string, *errors.Error) {
	args := b.Called(data)

	var err *errors.Error
	if v := args.Get(1); v != nil {
		err = v.(*errors.Error)
	}

	return args.String(0), err
}

func (b *BcryptMock) CompareHash(hash, data string) bool {
	args := b.Called(hash, data)
	return args.Bool(0)
}
