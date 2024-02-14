package storage

import (
	"context"
	"fmt"
	"log"
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

func (c *RedisCache) SetCalendarAsync(cal types.Calendar, err_chan chan error) {
	cal_json, err := cal.ToMap()
	if err != nil {
		log.Printf("err when mapping events : %v\n", err)
		err_chan <- err
		return
	}
	cal_stringified, err := types.JsonCalendarToString(cal_json)
	if err != nil {
		log.Printf("err when stringyfing events : %v\n", err)
		err_chan <- err
		return
	}
	c.Set(cal.Url, cal_stringified, time.Hour)
	err_chan <- nil
}
