package redis

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"github.com/3115826227/baby-fried-rice/module/crawler/log"
	"github.com/3115826227/baby-fried-rice/module/crawler/config"
	"time"
	"strconv"
)

var rds *redis.Client

const (
	DefaultIncr = 1
)

func init() {
	rds = redis.NewClient(&redis.Options{
		Addr:     config.Config.RedisUrl,
		Password: config.Config.RedisPassword,
		DB:       config.Config.RedisDB,
		PoolSize: 5,
	})

	pong, err := rds.Ping().Result()
	if err != nil {
		log.Logger.Error("REDIS", zap.String("e", err.Error()))
	} else {
		log.Logger.Info("REDIS", zap.String("url", config.Config.RedisUrl), zap.String("ping", pong))
	}
}

func Incr(key string) {
	rds.Incr(key)
}

func IncrGet(key string) (int64, error) {
	str, err := rds.Get(key).Result()
	if err != nil {
		return -1, err
	}
	incr, err := strconv.Atoi(str)
	if err != nil {
		log.Logger.Warn(err.Error())
		return -1, err
	}
	return int64(incr), nil
}

func HashIncr(field, key string) {
	rds.HIncrBy(field, key, DefaultIncr)
}

func HashIncrGet(field, key string) (int64, error) {
	str, err := rds.HGet(field, key).Result()
	if err != nil {
		return -1, err
	}
	incr, err := strconv.Atoi(str)
	if err != nil {
		log.Logger.Warn(err.Error())
		return -1, err
	}
	return int64(incr), nil
}

func Add(key, value string) {
	if Exist(key) {
		return
	}
	rds.Append(key, value)
	rds.Expire(key, 3*time.Hour)
}

func Exist(key string) bool {
	count, err := rds.Exists(key).Result()
	if err != nil || count != 1 {
		return false
	}
	return true
}

func Get(key string) string {
	str, _ := rds.Get(key).Result()
	return str
}

func HashAdd(field, key, value string) {
	rds.HSet(field, key, value)
	rds.Expire(field, 3*time.Hour)
}

func HashExist(field, key string) bool {
	count, err := rds.HExists(field, key).Result()
	if err != nil || !count {
		return false
	}
	return true
}

func HashGet(field, key string) string {
	str, _ := rds.HGet(field, key).Result()
	return str
}

func HashGetAll(field string) (result map[string]string) {
	var err error
	result, err = rds.HGetAll(field).Result()
	if err != nil {
		return nil
	}
	return result
}

func SAdd(key, value string) {
	rds.SAdd(key, value)
}

func SMem(key string) []string {
	res, err := rds.SMembers(key).Result()
	if err != nil {
		return nil
	}
	return res
}

func LPush(key string, value string) {
	rds.LPush(key, value)
}

func BRPop(key string) (result []string) {
	return rds.BLPop(time.Second, key).Val()
}
