package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	errors2 "github.com/andreis3/auth-ms/internal/domain/errors"
	adapter2 "github.com/andreis3/auth-ms/internal/domain/interfaces/adapter"
)

type Cache struct {
	client  *redis.Client
	metrics adapter2.Prometheus
	tracer  adapter2.Tracer
}

func NewCache(client *redis.Client, metrics adapter2.Prometheus, tracer adapter2.Tracer) *Cache {
	return &Cache{
		client:  client,
		metrics: metrics,
		tracer:  tracer,
	}
}

func (c *Cache) Get(ctx context.Context, key string, target any) (bool, *errors2.Error) {
	ctx, span := c.tracer.Start(ctx, "Cache.Get")
	start := time.Now()
	defer func() {
		end := time.Since(start)
		c.metrics.ObserveInstructionDBDuration("redis", "cache", "get", float64(end.Milliseconds()))
		span.End()
	}()
	result, err := c.client.Get(ctx, key).Result()

	if errors2.Is(err, redis.Nil) {
		return false, nil
	}

	if err = json.Unmarshal([]byte(result), target); err != nil {
		return false, errors2.ErrorGetCache(err)
	}

	return true, nil
}

func (c *Cache) Set(ctx context.Context, key string, value any, ttlSeconds int) *errors2.Error {
	ctx, span := c.tracer.Start(ctx, "Cache.Set")
	start := time.Now()
	defer func() {
		end := time.Since(start)
		c.metrics.ObserveInstructionDBDuration("redis", "cache", "set", float64(end.Milliseconds()))
		span.End()
	}()
	bytes, err := json.Marshal(value)

	if err != nil {
		return errors2.ErrorSetCache(err)
	}

	err = c.client.Set(ctx, key, string(bytes), time.Duration(ttlSeconds)*time.Second).Err()

	if err != nil {
		return errors2.ErrorSetCache(err)
	}

	return nil
}
