package config

import (
	"fmt"
	"github.com/jinzhu/configor"
	"os"
)

const (
	DayLayout = "2006-01-02"

	TutorEmergency = "紧急"
	TutorGeneral = "一般"
)

var Config = struct {
	MysqlUrl string `env:"MYSQL_URL" required:"true"`
	TokenSecret string
}{}

var Root = os.Getenv("GOPATH") + "/src/github.com/3115826227/baby-fried-rice/module/public"

func init() {
	var err error
	if err = configor.Load(&Config, Root+"etc/config.yaml"); err != nil {
		fmt.Println(err.Error())
	}
	var ok bool
	Config.TokenSecret, ok = os.LookupEnv("TOKEN_SECRET")
	if !ok {
		Config.TokenSecret = "baby"
	}
}
