package cache

import (
	"baby-fried-rice/internal/pkg/kit/cache"
	"github.com/go-redis/redis"
)

var (
	client *redis.Client
)

func InitCache(addr, passwd string, db int) (err error) {
	client, err = cache.NewRedis(addr, passwd, db)
	return
}
