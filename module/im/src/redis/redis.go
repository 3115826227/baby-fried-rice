package redis

import (
	"github.com/3115826227/baby-fried-rice/module/im/src/config"
	"github.com/3115826227/baby-fried-rice/module/im/src/log"
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

func Get(key string) (string, error) {
	return rds.Get(key).Result()
}

func SAdd(key string, value interface{}) {
	rds.SAdd(key, value)
}

func SMember(key string) ([]string, error) {
	return rds.SMembers(key).Result()
}

func HashGet(key string) (map[string]string, error) {
	return rds.HGetAll(key).Result()
}

func HashAdd(key, field, value string) {
	if err := rds.HSet(key, field, value).Err(); err != nil {
		log.Logger.Warn(err.Error())
	}
}

func BLPop(key string) ([]string, error) {
	return rds.BLPop(5*time.Second, key).Result()
}

func LPop(key string) (string, error) {
	return rds.LPop(key).Result()
}

func Push(key string, value interface{}) {
	err := rds.LPush(key, value).Err()
	if err != nil {
		log.Logger.Warn(err.Error())
	}
}
