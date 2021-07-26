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

type redisCache struct {
	lc  log.Logging
	rds *redis.Client
}

func NewRedisClient(addr, passwd string, db int) (rds *redis.Client, err error) {
	rds, err = newRedis(addr, passwd, db)
	if err != nil {
		err = errors.Wrap(err, "failed to new redis")
		return
	}
	return
}

func InitCache(addr, passwd string, db int, lc log.Logging) (interfaces.Cache, error) {
	rds, err := newRedis(addr, passwd, db)
	if err != nil {
		err = errors.Wrap(err, "failed to new redis")
		return nil, err
	}
	rc := &redisCache{
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

func (c *redisCache) Add(key string, value string) error {
	return c.rds.Set(key, value, DefaultExpiration).Err()
}

func (c *redisCache) Get(key string) (string, error) {
	return c.rds.Get(key).Result()
}

func (c *redisCache) Del(key string) error {
	return c.rds.Del(key).Err()
}

func (c *redisCache) Info() (string, error) {
	return c.rds.Info().Result()
}

func (c *redisCache) HSet(key, field string, value interface{}) error {
	return c.rds.HSet(key, field, value).Err()
}

func (c *redisCache) HGet(key, field string) (string, error) {
	return c.rds.HGet(key, field).Result()
}

func (c *redisCache) HGetAll(key string) (map[string]string, error) {
	return c.rds.HGetAll(key).Result()
}

func (c *redisCache) HDel(key string, field ...string) error {
	return c.rds.HDel(key, field...).Err()
}
