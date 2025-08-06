package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/andreis3/customers-ms/internal/domain/errors"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/adapter"
)

type Cache struct {
	client  *redis.Client
	metrics adapter.Prometheus
	tracer  adapter.Tracer
}

func NewCache(client *redis.Client, metrics adapter.Prometheus, tracer adapter.Tracer) *Cache {
	return &Cache{
		client:  client,
		metrics: metrics,
		tracer:  tracer,
	}
}

func (c *Cache) Get(ctx context.Context, key string, target any) (bool, *errors.Error) {
	ctx, span := c.tracer.Start(ctx, "Cache.Get")
	start := time.Now()
	defer func() {
		end := time.Since(start)
		c.metrics.ObserveInstructionDBDuration("redis", "cache", "get", float64(end.Milliseconds()))
		span.End()
	}()
	result, err := c.client.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return false, nil
	}

	if err := json.Unmarshal([]byte(result), target); err != nil {
		return false, errors.ErrorGetCache(err)
	}

	return true, nil
}

func (c *Cache) Set(ctx context.Context, key string, value any, ttlSeconds int) *errors.Error {
	ctx, span := c.tracer.Start(ctx, "Cache.Set")
	start := time.Now()
	defer func() {
		end := time.Since(start)
		c.metrics.ObserveInstructionDBDuration("redis", "cache", "set", float64(end.Milliseconds()))
		span.End()
	}()
	bytes, err := json.Marshal(value)

	if err != nil {
		return errors.ErrorSetCache(err)
	}

	err = c.client.Set(ctx, key, string(bytes), time.Duration(ttlSeconds)*time.Second).Err()

	if err != nil {
		return errors.ErrorSetCache(err)
	}

	return nil
}
