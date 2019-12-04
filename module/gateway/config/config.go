package config

import (
	"fmt"
	"github.com/jinzhu/configor"
	"net/url"
	"os"
	"time"
)

var Config = struct {
	Redis struct {
		URL                string `env:"REDIS_URL" default:"127.0.0.1:6379"`
		Password           string `env:"REDIS_PASSWORD"`
		Db                 string `env:"REDIS_DB" default:"0"`
		MaxIdleConnections int    `env:"REDIS_MAX_IDLE_CONNECTIONS" default:"3"`
		IdleTimeoutSeconds int    `env:"REDIS_IDLE_TIMEOUT_SECONDS" default:"240"`
		IdleTimeout        time.Duration
	}

	UserUrl      string `env:"USER_URL" required:"true"`
	ParseUserUrl *url.URL
}{}

var Root = os.Getenv("GOPATH") + "/github.com/3115826227/baby-fried-rice/module/gateway"

func init() {
	var err error
	if err = configor.Load(&Config, Root+"/etc/config.yaml"); err != nil {
		fmt.Println(err.Error())
	}
	Config.ParseUserUrl, err = url.Parse(Config.UserUrl)
	if err != nil {
		panic("error user url")
	}

}
