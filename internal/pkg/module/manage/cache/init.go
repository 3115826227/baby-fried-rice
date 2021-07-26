package cache

import (
	"baby-fried-rice/internal/pkg/kit/cache"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/log"
	"github.com/go-redis/redis"
)

var (
	c   interfaces.Cache
	rds *redis.Client
)

func GetRedisClient() *redis.Client {
	return rds
}

func InitRedisClient(addr, passwd string, db int) (err error) {
	rds, err = cache.NewRedisClient(addr, passwd, db)
	if err != nil {
		return
	}
	return
}

func GetCache() interfaces.Cache {
	return c
}

func InitCache(addr, passwd string, db int, lc log.Logging) (err error) {
	c, err = cache.InitCache(addr, passwd, db, lc)
	return
}
