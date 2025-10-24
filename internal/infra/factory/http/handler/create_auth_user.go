package handler

import (
	"github.com/andreis3/auth-ms/internal/adapter/input/http/handler"
	"github.com/andreis3/auth-ms/internal/adapter/output/repository"
	"github.com/andreis3/auth-ms/internal/adapter/output/security"
	"github.com/andreis3/auth-ms/internal/app/command"
	"github.com/andreis3/auth-ms/internal/app/service"
	adapter2 "github.com/andreis3/auth-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/auth-ms/internal/infra/config"
	db2 "github.com/andreis3/auth-ms/internal/infra/db"
	"github.com/andreis3/auth-ms/internal/infra/shared"
)

type CreateAuthUser struct {
	db      *db2.Postgres
	redis   *db2.Redis
	log     adapter2.Logger
	metrics adapter2.Prometheus
	tracer  adapter2.Tracer
	conf    *config.Configs
}

func NewCreateAuthUser(database *db2.Postgres, redis *db2.Redis, log adapter2.Logger, metrics adapter2.Prometheus, tracer adapter2.Tracer, conf *config.Configs) *CreateAuthUser {
	return &CreateAuthUser{database, redis, log, metrics, tracer, conf}
}

func (f *CreateAuthUser) NewCreateAuthUser() *handler.CreateAuthUserHandler {
	crypto := security.NewBcrypt()
	cmd := newCreateAuthUser(f.db, crypto, f.log, f.tracer, f.metrics)
	return handler.NewCreateAuthUserHandler(cmd, f.metrics, f.log, f.tracer)
}

func newCreateAuthUser(
	db *db2.Postgres,
	crypto adapter2.Bcrypt,
	log adapter2.Logger,
	tracer adapter2.Tracer,
	metrics adapter2.Prometheus,
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
