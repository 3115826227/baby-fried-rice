package config

import (
	"os"

	"github.com/jinzhu/configor"
)

var Config = struct {
	MysqlUrl      string `env:"MYSQL_URL" required:"true"`
	RedisUrl      string `env:"REDIS_URL" required:"true"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB" default:"9"`
	Kafka         string `env:"KAFKA" required:"true"`
}{}

const (
	ChatTopic          = "chat"
	ChatNewMessageKey  = "chat:new:message"
	ChatReadMessageKey = "chat:read:message"
	DefaultMessageSize = 10
	DefaultPartition   = 0
)

var Root = os.Getenv("GOPATH") + "/src/github.com/3115826227/baby-fried-rice/module/im"

func init() {
	var err error
	if err = configor.Load(&Config, "/etc/config.yaml"); err != nil {
		panic(err)
	}
}
