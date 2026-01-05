package handler

import (
	"github.com/andreis3/auth-ms/internal/adapter/input/http/handler"
	"github.com/andreis3/auth-ms/internal/adapter/output/repository"
	"github.com/andreis3/auth-ms/internal/adapter/output/security"
	"github.com/andreis3/auth-ms/internal/app/command"
	"github.com/andreis3/auth-ms/internal/app/service"
	"github.com/andreis3/auth-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/auth-ms/internal/infra/config"
	"github.com/andreis3/auth-ms/internal/infra/db"
	"github.com/andreis3/auth-ms/internal/infra/shared"
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
	userService := service.NewUserService(userRepository, tracer, log)
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
