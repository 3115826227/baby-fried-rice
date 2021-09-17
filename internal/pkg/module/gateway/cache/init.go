package cache

import (
	"baby-fried-rice/internal/pkg/kit/cache"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/log"
	"baby-fried-rice/internal/pkg/kit/models"
)

var (
	c interfaces.Cache
)

func GetCache() interfaces.Cache {
	return c
}

func InitCache(redisConf models.Redis, lc log.Logging) (err error) {
	c, err = cache.InitCache(redisConf.Addr, redisConf.Password, redisConf.DB, lc)
	return
}
