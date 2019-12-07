package config

import (
	"github.com/jinzhu/configor"
	"os"
)

const (
	DayLayout  = "2006-01-02"
	TimeLayout = "2006-01-02 15:04:05"
)

const (
	TrainSeatTriggerStatus = "train:seat:trigger:status"
	TrainMetaTriggerStatus = "train:meta:trigger:status"

	TrainMetaStationTriggerPrefix = "train:meta:station:trigger"
	TrainMetaCodeTriggerPrefix    = "train:meta:code:trigger"

	TrainSeatTriggerPrefix = "train:seat:trigger"

	StationNameCodeKey = "station:name:code"
	StationCodeNameKey = "station:code:name"

	TrainStationCityUpdateKey = "train:station:city:update"

	RunningStatus = "run"
	SuccessStatus = "success"
	FailStatus    = "fail"
)

var Config = struct {
	PostgresUrl    string `env:"POSTGRES_URL" required:"true"`
	RedisUrl       string `env:"REDIS_URL" required:"true"`
	RedisPassword  string `env:"REDIS_PASSWORD" default:""`
	RedisDB        int    `env:"REDIS_DB" required:"true"`
	ConsumerNumber int    `env:"CONSUMER_NUMBER" default:"1"`
}{}

var BatchLoadAmount = 100

var Root = os.Getenv("GOPATH") + "/src/github.com/3115826227/baby-fried-rice/module/crawler"

func init() {
	var err error
	if err = configor.Load(&Config, Root+"/etc/config.yaml"); err != nil {
		panic(err.Error())
	}

}
