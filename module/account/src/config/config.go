package config

import (
	"github.com/jinzhu/configor"
	"os"
	"path/filepath"
)

var Config = struct {
	PostgresUrl   string `env:"POSTGRES_URL" required:"true"`
	RedisUrl      string `env:"REDIS_URL" required:"true"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB" default:"9"`
	TokenSecret   string
}{}

var Root = os.Getenv("GOPATH") + "/src/github.com/3115826227/baby-fried-rice/module/account"

func init() {
	var err error
	if err = configor.Load(&Config, filepath.Join(Root, "etc/config.yaml")); err != nil {
		panic(err)
	}
	var ok bool
	Config.TokenSecret, ok = os.LookupEnv("TOKEN_SECRET")
	if !ok {
		Config.TokenSecret = "baby"
	}
}
