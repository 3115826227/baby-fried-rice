package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

type Conf struct {
	Log struct {
		LogLevel string `json:"log_level"`
		LogPath  string `json:"log_path"`
	} `json:"log"`

	Server struct {
		Name     string `json:"name"`
		Serial   int    `json:"serial"`
		Addr     string `json:"addr"`
		Port     int    `json:"port"`
		Register string `json:"register"`
	} `json:"server"`

	Redis struct {
		RedisUrl      string `json:"redis_url"`
		RedisPassword string `json:"redis_password"`
		RedisDB       int    `json:"redis_db"`
	} `json:"redis"`

	Etcd            []string `json:"etcd"`
	HealthyRollTime int64    `json:"healthy_roll_time"`

	Rpc struct {
		Client struct {
			CertFile string `json:"cert_file"`
		} `json:"client"`
	} `json:"rpc"`

	Servers struct {
		AccountDaoServer string `json:"account_dao_server"`
		SpaceDaoServer   string `json:"space_dao_server"`
	}
}

var (
	config Conf
)

func GetConfig() Conf {
	return config
}

func readConfig() (err error) {
	viper.SetConfigFile("./res/config_dev.yaml") // 指定配置文件路径
	viper.SetConfigName("config_dev")            // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")                  // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./res/")                // 查找配置文件所在的路径
	err = viper.ReadInConfig()                   // 查找并读取配置文件
	if err != nil {                              // 处理读取配置文件的错误
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
}
