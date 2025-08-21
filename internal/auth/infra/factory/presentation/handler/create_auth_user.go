package handler

import (
	"github.com/andreis3/auth-ms/internal/auth/adapter/input/http/handler"
	"github.com/andreis3/auth-ms/internal/auth/adapter/output/repository"
	"github.com/andreis3/auth-ms/internal/auth/adapter/output/security"
	"github.com/andreis3/auth-ms/internal/auth/application/command"
	"github.com/andreis3/auth-ms/internal/auth/application/service"
	"github.com/andreis3/auth-ms/internal/auth/domain/interfaces/adapter"
	"github.com/andreis3/auth-ms/internal/auth/infra/config"
	"github.com/andreis3/auth-ms/internal/auth/infra/db"
	"github.com/andreis3/auth-ms/internal/auth/infra/shared"
)

type CreateAuthUser struct {
	db      *db.Postgres
	redis   *db.Redis
	log     adapter.Logger
	metrics adapter.Prometheus
	tracer  adapter.Tracer
	conf    *config.Configs
}

func NewCreateAuthUser(database *db.Postgres, redis *db.Redis, log adapter.Logger, metrics adapter.Prometheus, tracer adapter.Tracer, conf *config.Configs) *CreateAuthUser {
	return &CreateAuthUser{database, redis, log, metrics, tracer, conf}
}

func (f *CreateAuthUser) NewCreateAuthUser() *handler.CreateAuthUserHandler {
	crypto := security.NewBcrypt()
	cmd := newCreateAuthUser(f.db, crypto, f.log, f.tracer, f.metrics)
	return handler.NewCreateAuthUserHandler(cmd, f.metrics, f.log, f.tracer)
}

func newCreateAuthUser(
	db *db.Postgres,
	crypto adapter.Bcrypt,
	log adapter.Logger,
	tracer adapter.Tracer,
	metrics adapter.Prometheus,
) *command.CreateAuthUser {
	userRepository := repository.NewUserRepository(db, metrics, tracer)
	userService := service.NewCustomerService(userRepository, tracer, log)
	utils := shared.Utils{}
	return command.NewCreateAuthUser(
		userRepository,
		userService,
		crypto,
		log,
		tracer,
		utils,
	)
}
