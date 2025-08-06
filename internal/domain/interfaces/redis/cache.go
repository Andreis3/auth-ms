package redis

import (
	"context"

	"github.com/andreis3/customers-ms/internal/domain/errors"
)

type Cache interface {
	Get(ctx context.Context, key string, target any) (bool, *errors.Error)
	Set(ctx context.Context, key string, value any, ttlSeconds int) *errors.Error
}
