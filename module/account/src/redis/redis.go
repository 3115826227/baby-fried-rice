package redis

import (
	"github.com/go-redis/redis"
	"github.com/3115826227/baby-fried-rice/module/account/src/log"
	"time"
	"github.com/3115826227/baby-fried-rice/module/account/src/config"
)

var rds *redis.Client

func init() {
	rds = redis.NewClient(&redis.Options{
		Addr:     config.Config.RedisUrl,
		Password: config.Config.RedisPassword,
		PoolSize: 20,
		DB:       config.Config.RedisDB,
	})
	if err := rds.Ping().Err(); err != nil {
		log.Logger.Warn(err.Error())
		return
	}
}

func AddAccountToken(key, value string) {
	Add(key, value, 3*time.Hour)
}

func Add(key, value string, expire time.Duration) {
	rds.Del(key)
	rds.Append(key, value)
	rds.Expire(key, expire)
}
