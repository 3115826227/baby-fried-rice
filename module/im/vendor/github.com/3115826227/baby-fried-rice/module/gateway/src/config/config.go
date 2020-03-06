package config

import (
	"github.com/jinzhu/configor"
	"net/url"
	"os"
	"time"
)

var Config = struct {
	Redis struct {
		URL                string `env:"REDIS_URL" default:"127.0.0.1:6379"`
		Password           string `env:"REDIS_PASSWORD"`
		Db                 int    `env:"REDIS_DB" default:"9"`
		MaxIdleConnections int    `env:"REDIS_MAX_IDLE_CONNECTIONS" default:"3"`
		IdleTimeoutSeconds int    `env:"REDIS_IDLE_TIMEOUT_SECONDS" default:"240"`
		IdleTimeout        time.Duration
	}

	AccountUrl       string `env:"ACCOUNT_URL" required:"true"`
	PublicUrl        string `env:"PUBLIC_URL" required:"true"`
	SquareUrl        string `env:"SQUARE_URL" required:"true"`
	ImUrl            string `env:"IM_URL" required:"true"`
	ParserAccountUrl *url.URL
	ParserPublicUrl  *url.URL
	ParserSquareUrl  *url.URL
	ParserImUrl      *url.URL
}{}

var Root = os.Getenv("GOPATH") + "/src/github.com/3115826227/baby-fried-rice/module/gateway"

func init() {
	var err error
	if err = configor.Load(&Config, "etc/config.yaml"); err != nil {
		panic(err)
	}
	Config.ParserAccountUrl, err = url.Parse(Config.AccountUrl)
	if err != nil {
		panic("error user url")
	}
	Config.ParserPublicUrl, err = url.Parse(Config.PublicUrl)
	if err != nil {
		panic("error user url")
	}
	Config.ParserSquareUrl, err = url.Parse(Config.SquareUrl)
	if err != nil {
		panic("error user url")
	}
	Config.ParserImUrl, err = url.Parse(Config.ImUrl)
	if err != nil {
		panic("error user url")
	}

}
