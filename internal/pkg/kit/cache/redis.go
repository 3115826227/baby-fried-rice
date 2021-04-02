package cache

import (
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/log"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"time"
)

const (
	DefaultExpiration = 3 * time.Hour
)

type RedisCache struct {
	lc  log.Logging
	rds *redis.Client
}

var (
	rc *RedisCache
)

func InitCache(addr, passwd string, db int, lc log.Logging) (interfaces.Cache, error) {
	rds, err := newRedis(addr, passwd, db)
	if err != nil {
		err = errors.Wrap(err, "failed to new redis")
		return nil, err
	}
	rc = &RedisCache{
		lc:  lc,
		rds: rds,
	}
	return rc, nil
}

func newRedis(addr, passwd string, db int) (rds *redis.Client, err error) {
	rds = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		PoolSize: 20,
		DB:       db,
	})
	if err = rds.Ping().Err(); err != nil {
		return
	}
	return
}

func (c *RedisCache) Add(key string, value string) error {
	return c.rds.Set(key, value, DefaultExpiration).Err()
}

func (c *RedisCache) Get(key string) (string, error) {
	return c.rds.Get(key).Result()
}

func (c *RedisCache) Del(key string) error {
	return c.rds.Del(key).Err()
}

func (c *RedisCache) HSet(key, field string, value interface{}) error {
	return c.rds.HSet(key, field, value).Err()
}

func (c *RedisCache) HGet(key, field string) (string, error) {
	return c.rds.HGet(key, field).Result()
}

func (c *RedisCache) HGetAll(key string) (map[string]string, error) {
	return c.rds.HGetAll(key).Result()
}
