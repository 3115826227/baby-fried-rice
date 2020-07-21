package config

import (
	"github.com/jinzhu/configor"
	"os"
)

var Config = struct {
	MysqlUrl string `env:"MYSQL_URL" required:"true"`
}{}

var Key = struct {
	Key []struct {
		Id        int    `env:"ID" required:"true"`
		AccessKey string `env:"ACCESS_KEY" required:"true"`
		SecretKey string `env:"SECRET_KEY" required:"true"`
	} `env:"KEY" required:"true"`
}{}

const (
	TimeLayout = "2006-01-02 15:04:05"
	DateLayout = "2006-01-02"
)

var Root = os.Getenv("GOPATH") + "/src/github.com/3115826227/baby-fried-rice/module/file"

func init() {
	var err error
	if err = configor.Load(&Config, "etc/config.yaml"); err != nil {
		panic(err)
	}
	if err = configor.Load(&Key, "etc/key.yaml"); err != nil {
		panic(err)
	}
}
