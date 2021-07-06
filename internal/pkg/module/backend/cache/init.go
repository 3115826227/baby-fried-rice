package cache

import (
	"baby-fried-rice/internal/pkg/kit/cache"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/log"
)

var (
	c interfaces.Cache
)

func GetCache() interfaces.Cache {
	return c
}

func InitCache(addr, passwd string, db int, lc log.Logging) (err error) {
	c, err = cache.InitCache(addr, passwd, db, lc)
	return
}
