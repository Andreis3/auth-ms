package decorator

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log/slog"

	"github.com/andreis3/customers-ms/internal/domain/entity"
	"github.com/andreis3/customers-ms/internal/domain/errors"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/adapter"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/postgres"
	"github.com/andreis3/customers-ms/internal/domain/interfaces/redis"
	"github.com/andreis3/customers-ms/internal/infra/repositories/criteria"
)

type CachedAddressesRepository struct {
	inner  postgres.AddressRepository
	cache  redis.Cache
	log    adapter.Logger
	tracer adapter.Tracer
}

func NewCachedAddressesRepository(inner postgres.AddressRepository, cache redis.Cache, log adapter.Logger, tracer adapter.Tracer) postgres.AddressesSearchRepository {
	return &CachedAddressesRepository{
		inner:  inner,
		cache:  cache,
		log:    log,
		tracer: tracer,
	}
}

func (c *CachedAddressesRepository) SearchAddresses(ctx context.Context, params criteria.AddressSearchCriteria) (*[]entity.Address, *errors.Error) {
	ctx, span := c.tracer.Start(ctx, "CachedAddressesRepository.SearchAddresses")
	defer func() {
		span.End()
	}()
	cachedKey, err := c.generateCacheKey(params)
	if err != nil {
		return nil, errors.ErrorGenerateCacheKey(err)
	}

	var cached []entity.Address

	found, err := c.cache.Get(ctx, cachedKey, &cached)
	if err != nil {
		c.log.ErrorJSON("Failed to get cache",
			slog.String("trace_id", span.SpanContext().TraceID()),
			slog.Any("error", err))
	}

	if found {
		c.log.InfoJSON("Found addresses in cache",
			slog.String("trace_id", span.SpanContext().TraceID()),
			slog.Any("addresses", cached))
		return &cached, nil
	}

	addresses, err := c.inner.SearchAddresses(ctx, params)
	if err != nil {
		c.log.ErrorJSON("Failed to get addresses",
			slog.String("trace_id", span.SpanContext().TraceID()),
			slog.Any("error", err))
		return nil, errors.ErrorSearchAddresses(err)
	}

	const ttl = 60
	err = c.cache.Set(ctx, cachedKey, addresses, ttl)
	if err != nil {
		c.log.ErrorJSON("Failed to set cache",
			slog.String("trace_id", span.SpanContext().TraceID()),
			slog.Any("error", err))
	}
	c.log.InfoJSON("Set addresses in cache",
		slog.String("trace_id", span.SpanContext().TraceID()),
		slog.Int("ttl", ttl),
		slog.Any("addresses", addresses))
	return addresses, nil
}

func (c *CachedAddressesRepository) generateCacheKey(params criteria.AddressSearchCriteria) (string, *errors.Error) {
	bytes, err := json.Marshal(params)
	if err != nil {
		return "", errors.ErrorGenerateCacheKey(err)
	}
	sum := sha256.Sum256(bytes)
	return "customers:search:" + hex.EncodeToString(sum[:]), nil
}
