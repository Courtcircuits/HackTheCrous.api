package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/Courtcircuits/HackTheCrous.api/types"
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

func (c *RedisCache) SetCalendarAsync(cal types.Calendar) {
	cal_json, err := cal.ToMap()
	if err != nil {
		panic(err)
	}
	cal_stringified, err := types.JsonCalendarToString(cal_json)
	if err != nil {
		panic(err) //can be enhanced by using channels to return errors
	}
	err = c.Set(cal.Url, cal_stringified, time.Hour)
	if err != nil {
		panic(err)
	}
}
