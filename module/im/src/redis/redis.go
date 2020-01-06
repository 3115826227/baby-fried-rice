package redis

import (
	"github.com/3115826227/baby-fried-rice/module/gateway/src/config"
	"github.com/3115826227/baby-fried-rice/module/gateway/src/log"
	"github.com/go-redis/redis"
)

var rds *redis.Client

func init() {
	rds = redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.URL,
		Password: config.Config.Redis.Password,
		PoolSize: 20,
		DB:       config.Config.Redis.Db,
	})
	if err := rds.Ping().Err(); err != nil {
		log.Logger.Warn(err.Error())
		return
	}
}

func Add(key, value string) {
	rds.Append(key, value)
}

func Get(key string) (string, error) {
	str, err := rds.Get(key).Result()
	if err != nil {
		return "", err
	}
	return str, nil
}
