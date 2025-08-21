package service

import (
	"context"

	errors2 "github.com/andreis3/auth-ms/internal/auth/domain/errors"
	"github.com/andreis3/auth-ms/internal/auth/domain/interfaces/adapter"
	"github.com/andreis3/auth-ms/internal/auth/domain/port"
)

type UserService struct {
	repository port.UserRepository
	tracer     adapter.Tracer
	log        adapter.Logger
}

func NewCustomerService(repository port.UserRepository, trace adapter.Tracer, log adapter.Logger) *UserService {
	return &UserService{
		repository: repository,
		tracer:     trace,
		log:        log,
	}
}

func (s *UserService) ValidateEmailAvailability(ctx context.Context, email string) *errors2.Error {
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
		err := errors2.ErrorAlreadyExists(user.PublicID())
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
