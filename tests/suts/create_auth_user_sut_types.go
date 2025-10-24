//go:build unit

package suts

import (
	"github.com/andreis3/auth-ms/internal/app/command"
	"github.com/andreis3/auth-ms/tests/mocks/app/mservice"
	"github.com/andreis3/auth-ms/tests/mocks/infra/madapters"
	"github.com/andreis3/auth-ms/tests/mocks/infra/mrepository"
)

type CreateAuthUserSut struct {
	Repo    *mrepository.UserRepositoryMock
	Service *mservice.UserServiceMock
	Bcrypt  *madapters.BcryptMock
	Log     *madapters.LoggerMock
	Tracer  *madapters.TracerMock
	Span    *madapters.SpanMock
	Sc      *madapters.SpanContextMock
	Utils   *madapters.UtilsMock
	Cmd     *command.CreateAuthUser
}

func MakeCreateAuthUserSut() *CreateAuthUserSut {
	return &CreateAuthUserSut{
		Repo:    new(mrepository.UserRepositoryMock),
		Service: new(mservice.UserServiceMock),
		Bcrypt:  new(madapters.BcryptMock),
		Log:     new(madapters.LoggerMock),
		Tracer:  new(madapters.TracerMock),
		Span:    new(madapters.SpanMock),
		Sc:      new(madapters.SpanContextMock),
		Utils:   new(madapters.UtilsMock),
	}
}

func (s *CreateAuthUserSut) Build() *command.CreateAuthUser {
	s.Cmd = command.NewCreateAuthUser(s.Repo, s.Service, s.Bcrypt, s.Log, s.Tracer, s.Utils)
	return s.Cmd
}
