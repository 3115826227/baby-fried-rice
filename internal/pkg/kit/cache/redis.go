package cache

import (
	"github.com/go-redis/redis"
)

func NewRedis(addr, passwd string, db int) (rds *redis.Client, err error) {
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
