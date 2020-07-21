package config

import (
	"github.com/jinzhu/configor"
	"os"
)

var Config = struct {
	RedisUrl      string `env:"REDIS_URL" required:"true"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB" default:"9"`
	TokenSecret   string

	AccountDaoUrl string `env:"ACCOUNT_DAO_URL" required:"true"`
	ImUrl         string `env:"IM_URL" required:"true"`
}{}

const (
	TimeLayout  = "2006-01-02 15:04:05"
	DateLayout  = "2006-01-02"
	SuccessCode = 0
)

var Root = os.Getenv("GOPATH") + "/src/github.com/3115826227/baby-fried-rice/module/root-account"

func init() {
	var err error
	if err = configor.Load(&Config, "etc/config.yaml"); err != nil {
		panic(err)
	}
	var ok bool
	Config.TokenSecret, ok = os.LookupEnv("TOKEN_SECRET")
	if !ok {
		Config.TokenSecret = "baby"
	}
}
