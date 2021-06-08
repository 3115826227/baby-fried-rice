package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Conf struct {
	Server struct {
		Addr string `json:"addr"`
		Port int    `json:"port"`
	} `json:"server"`

	Redis struct {
		RedisUrl      string `json:"redis_url"`
		RedisPassword string `json:"redis_password"`
		RedisDB       int    `json:"redis_db"`
	} `json:"cache"`

	TokenSecret string `json:"token_secret"`

	AccountDaoUrl string `json:"account_dao_url"`
}

const (
	TimeLayout  = "2006-01-02 15:04:05"
	DateLayout  = "2006-01-02"
	SuccessCode = 0
)

var (
	config Conf
)

func GetConfig() Conf {
	return config
}

func readConfig() (err error) {
	viper.SetConfigFile("./res/config_dev.yaml") // 指定配置文件路径
	viper.SetConfigName("config")            // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")              // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./res/")            // 查找配置文件所在的路径
	err = viper.ReadInConfig()               // 查找并读取配置文件
	if err != nil {                          // 处理读取配置文件的错误
		err = errors.New(fmt.Sprintf("Fatal error config file: %s \n", err))
		return
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		err = errors.New(fmt.Sprintf("Fatal error config file: %s \n", err))
		return
	}
	return
}

func init() {
	var err error
	if err = readConfig(); err != nil {
		panic(err)
	}
	var ok bool
	config.TokenSecret, ok = os.LookupEnv("TOKEN_SECRET")
	if !ok {
		config.TokenSecret = "baby"
	}
}
