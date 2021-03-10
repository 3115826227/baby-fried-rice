package redis

import (
	"github.com/go-redis/redis"
)

func NewRedis(addr, pwd string, db int) (rds *redis.Client, err error) {
	rds = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		PoolSize: 20,
		DB:       db,
	})
	if err = rds.Ping().Err(); err != nil {
		return
	}
	return
}
