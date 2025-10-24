package madapters

import "github.com/stretchr/testify/mock"

type UtilsMock struct{ mock.Mock }

func (u *UtilsMock) UUID() string {
	args := u.Called()
	return args.String(0)
}
