package db

import (
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/andreis3/auth-ms/internal/auth/infra/config"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(conf config.Configs) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", conf.RedisHost, conf.RedisPort),
		Password:     conf.RedisPassword,
		DB:           conf.RedisDB,
		PoolSize:     100,
		MinIdleConns: 10,
	})

	return &Redis{
		client: client,
	}
}

func (r *Redis) Close() {
	r.client.Close()
}

func (r *Redis) Client() *redis.Client {
	return r.client
}
