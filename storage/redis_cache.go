package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/Courtcircuits/HackTheCrous.api/util"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache() *RedisCache {
	return &RedisCache{
		client: redis.NewClient(&redis.Options{
			Addr:     util.Get("REDISHOST") + ":" + util.Get("REDISPORT"),
			Password: util.Get("REDISPASSWORD"),
			Username: util.Get("REDISUSER"),
		}),
	}
}

func (c *RedisCache) Get(key string) (string, error) {
	ctx := context.Background()
	return c.client.Get(ctx, key).Result()
}

func (c *RedisCache) Set(key string, val string, exp time.Duration) error {
	ctx := context.Background()
	jsp := c.client.Set(ctx, key, val, exp)
	if jsp != nil {
		return fmt.Errorf("redis error : %q", jsp)
	}
	return nil
}
