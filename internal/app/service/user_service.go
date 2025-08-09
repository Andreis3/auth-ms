package service

import (
	"context"

	"github.com/andreis3/auth-ms/internal/domain/errors"
	"github.com/andreis3/auth-ms/internal/domain/interfaces/iadapter"
	"github.com/andreis3/auth-ms/internal/domain/interfaces/ipostgres"
)

type UserService struct {
	repository ipostgres.IUserRepository
	tracer     iadapter.ITracer
	log        iadapter.ILogger
}

func NewCustomerService(repository ipostgres.IUserRepository, trace iadapter.ITracer, log iadapter.ILogger) *UserService {
	return &UserService{
		repository: repository,
		tracer:     trace,
		log:        log,
	}
}

func (s *UserService) ValidateEmailAvailability(ctx context.Context, email string) *errors.Error {
	ctx, span := s.tracer.Start(ctx, "UserService.ValidateEmailAvailability")
	defer span.End()
	traceID := span.SpanContext().TraceID()
	s.log.InfoJSON("Validating email availability",
		map[string]any{
			"trace_id": traceID,
			"email":    email,
		})

	user, err := s.repository.FindUserByEmail(ctx, email)
	if err != nil {
		span.RecordError(err)
		s.log.ErrorJSON("Error finding user by email",
			map[string]any{
				"trace_id": traceID,
				"email":    email,
				"error":    err.Error(),
			})
		return err
	}

	if user != nil {
		err := errors.ErrorAlreadyExists(user.PublicID())
		span.RecordError(err)
		s.log.ErrorJSON("Email already exists",
			map[string]any{
				"trace_id":  traceID,
				"email":     email,
				"public_id": user.PublicID(),
			})
		return err
	}

	s.log.InfoJSON("Email is available",
		map[string]any{
			"trace_id": traceID,
			"email":    email,
		})
	return nil
}
