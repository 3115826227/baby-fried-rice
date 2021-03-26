package cache

import (
	"baby-fried-rice/internal/pkg/kit/cache"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/log"
)

var (
	Cache interfaces.Cache
)

func InitCache(addr, passwd string, db int, lc log.Logging) (err error) {
	Cache, err = cache.InitCache(addr, passwd, db, lc)
	return
}
