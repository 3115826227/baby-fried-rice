package redis

import (
	"github.com/3115826227/baby-fried-rice/module/gateway/config"
	"github.com/garyburd/redigo/redis"
	"time"
)

var rdsCache *redis.Pool

func init() {
	rdsCache = newClient()
}

func newClient() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     config.Config.Redis.MaxIdleConnections,
		IdleTimeout: config.Config.Redis.IdleTimeout,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			return dialWithDB("tcp", config.Config.Redis.URL, config.Config.Redis.Password, config.Config.Redis.Db)
		},
	}
}

func dialWithDB(network, address, password, DB string) (redis.Conn, error) {
	c, err := dial(network, address, password)
	if err != nil {
		return nil, err
	}
	if _, err := c.Do("SELECT", DB); err != nil {
		c.Close()
		return nil, err
	}
	return c, err
}

func dial(network, address, password string) (redis.Conn, error) {
	c, err := redis.Dial(network, address)
	if err != nil {
		return nil, err
	}
	if password != "" {
		if _, err := c.Do("AUTH", password); err != nil {
			c.Close()
			return nil, err
		}
	}
	return c, err
}

func Set(key string, value string, exp string) (bool, error) {
	conn := rdsCache.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value, "EX", exp)
	if err != nil {
		return false, err
	}
	return true, nil
}

func Get(key string) (string, error) {
	conn := rdsCache.Get()
	defer conn.Close()
	reply, err := conn.Do("GET", key)
	sVal, err := redis.String(reply, err)
	if err != nil {
		return "", err
	}
	return sVal, nil
}

func Delete(key string) bool {
	conn := rdsCache.Get()
	defer conn.Close()
	reply, err := conn.Do("DEL", key)
	_, err = redis.String(reply, err)
	if err != nil {
		return false
	}
	return true
}
