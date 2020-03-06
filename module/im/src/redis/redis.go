package redis

import (
	"github.com/3115826227/baby-fried-rice/module/im/src/config"
	"github.com/3115826227/baby-fried-rice/module/im/src/log"
	"github.com/garyburd/redigo/redis"
)

var rds redis.Conn

func init() {
	var err error
	rds, err = redis.Dial("tcp", config.Config.RedisUrl)
	if err != nil {
		log.Logger.Warn(err.Error())
		panic(err)
	}
}

func Get(key string) (string, error) {
	str, err := redis.String(rds.Do("get",key))
	if err != nil {
		return "", err
	}
	return str, nil
}
