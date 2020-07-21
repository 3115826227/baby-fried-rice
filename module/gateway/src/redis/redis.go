package redis

import (
	"github.com/3115826227/baby-fried-rice/module/gateway/src/config"
	"github.com/3115826227/baby-fried-rice/module/gateway/src/log"
	"github.com/go-redis/redis"
	"time"
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

func Add(key string, value interface{}) {
	rds.Set(key, value, 60*time.Minute)
}

func Get(key string) (string, error) {
	str, err := rds.Get(key).Result()
	if err != nil {
		log.Logger.Warn(err.Error())
		return "", err
	}
	return str, nil
}
