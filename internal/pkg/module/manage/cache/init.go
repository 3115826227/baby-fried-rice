package cache

import (
	"baby-fried-rice/internal/pkg/kit/cache"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/log"
	"baby-fried-rice/internal/pkg/kit/models"
	"github.com/go-redis/redis"
)

var (
	c   interfaces.Cache
	rds *redis.Client
)

func GetRedisClient() *redis.Client {
	return rds
}

func InitRedisClient(redisConf models.Redis) (err error) {
	rds, err = cache.NewRedisClient(redisConf.Addr, redisConf.Password, redisConf.DB)
	if err != nil {
		return
	}
	return
}

func GetCache() interfaces.Cache {
	return c
}

func InitCache(redisConf models.Redis, lc log.Logging) (err error) {
	c, err = cache.InitCache(redisConf.Addr, redisConf.Password, redisConf.DB, lc)
	return
}
