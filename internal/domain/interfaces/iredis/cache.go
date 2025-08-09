package iredis

import (
	"context"

	"github.com/andreis3/auth-ms/internal/domain/errors"
)

type ICache interface {
	Get(ctx context.Context, key string, target any) (bool, *errors.Error)
	Set(ctx context.Context, key string, value any, ttlSeconds int) *errors.Error
}
