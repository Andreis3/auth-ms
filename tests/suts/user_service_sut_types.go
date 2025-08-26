//go:build unit

package suts

import (
	"github.com/andreis3/auth-ms/internal/auth/app/service"
	"github.com/andreis3/auth-ms/tests/mocks/infra/madapters"
	"github.com/andreis3/auth-ms/tests/mocks/infra/mrepository"
)

type UserServiceSut struct {
	Repo    *mrepository.UserRepositoryMock
	Tracer  *madapters.TracerMock
	Span    *madapters.SpanMock
	Sc      *madapters.SpanContextMock
	Log     *madapters.LoggerMock
	Service *service.UserService
}

func MakeUserServiceSut() *UserServiceSut {
	return &UserServiceSut{
		Repo:   new(mrepository.UserRepositoryMock),
		Tracer: new(madapters.TracerMock),
		Span:   new(madapters.SpanMock),
		Sc:     new(madapters.SpanContextMock),
		Log:    new(madapters.LoggerMock),
	}
}

func (s *UserServiceSut) Build() *service.UserService {
	s.Service = service.NewCustomerService(s.Repo, s.Tracer, s.Log)
	return s.Service
}
