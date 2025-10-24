package command

import (
	"context"

	"github.com/andreis3/auth-ms/internal/app/dto"
	"github.com/andreis3/auth-ms/internal/app/mapper"
	"github.com/andreis3/auth-ms/internal/app/port/service"
	"github.com/andreis3/auth-ms/internal/domain/entity"
	"github.com/andreis3/auth-ms/internal/domain/errors"
	"github.com/andreis3/auth-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/auth-ms/internal/domain/port"
	"github.com/andreis3/auth-ms/internal/infra/logger"
)

type CreateAuthUser struct {
	userRepository port.UserRepository
	userService    service.UserService
	bcrypt         adapter.Bcrypt
	log            adapter.Logger
	tracer         adapter.Tracer
	utils          adapter.Utils
}

func NewCreateAuthUser(
	userRepository port.UserRepository,
	userService service.UserService,
	bcrypt adapter.Bcrypt,
	log adapter.Logger,
	tracer adapter.Tracer,
	utils adapter.Utils,
) *CreateAuthUser {
	return &CreateAuthUser{
		userRepository: userRepository,
		userService:    userService,
		bcrypt:         bcrypt,
		log:            log,
		tracer:         tracer,
		utils:          utils,
	}
}

func (c *CreateAuthUser) Execute(ctx context.Context, input dto.CreateAuthUserInput) (*dto.CreateAuthUserOutput, *errors.Error) {
	ctx, span := c.tracer.Start(ctx, "CreateAuthUser.Execute")
	defer span.End()
	tracerID := span.SpanContext().TraceID()
	c.log.InfoJSON("Creating auth user",
		map[string]any{
			"trace_id": tracerID,
			"body": logger.RedactStruct[dto.CreateAuthUserInput](input, "password",
				"password_confirm"),
		})

	user := mapper.ToUser(input)
	user.AssignPublicID(c.utils.UUID())
	user.AssignRole(entity.RoleUser)
	isValid := user.Validate()

	if isValid.HasErrors() {
		c.log.CriticalJSON("User validation failed",
			map[string]any{
				"trace_id": tracerID,
				"errors":   isValid.FieldErrorsFlat(),
			})
		span.RecordError(isValid)
		return nil, errors.InvalidEntity(isValid, "user")
	}

	err := c.userService.ValidateEmailAvailability(ctx, input.Email)
	if err != nil {
		c.log.ErrorJSON("Email validation failed",
			map[string]any{
				"trace_id": tracerID,
				"email":    input.Email,
				"error":    err.Error(),
			})
		return nil, err
	}

	hashedPassword, err := c.bcrypt.Hash(user.Password())
	if err != nil {
		c.log.CriticalJSON(
			"Error hashing password",
			map[string]any{
				"error": err.Error(),
			})
		return nil, err
	}
	user.AssignPasswordHash(hashedPassword)

	createUser, err := c.userRepository.CreateUser(ctx, user)
	if err != nil {
		c.log.ErrorJSON(
			"Error creating user",
			map[string]any{
				"error": err.Error(),
			})
		return nil, err
	}

	return mapper.ToCreateAuthUserOutput(createUser), nil

}
