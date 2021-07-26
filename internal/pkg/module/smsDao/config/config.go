package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Conf struct {
	Log struct {
		LogLevel string `json:"log_level"`
		LogPath  string `json:"log_path"`
	} `json:"log"`

	Redis struct {
		RedisUrl      string `json:"redis_url"`
		RedisPassword string `json:"redis_password"`
		RedisDB       int    `json:"redis_db"`
	} `json:"redis"`

	Server struct {
		Name     string `json:"name"`
		Serial   int    `json:"serial"`
		Addr     string `json:"addr"`
		Port     int    `json:"port"`
		Register string `json:"register"`
	} `json:"server"`

	Rpc struct {
		Server struct {
			CertFile string `json:"cert_file"`
			KeyFile  string `json:"key_file"`
		} `json:"server"`
	} `json:"rpc"`

	Etcd            []string `json:"etcd"`
	HealthyRollTime int64    `json:"healthy_roll_time"`

	MysqlUrl string `json:"mysql_url"`

	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
}

var (
	config Conf
)

func GetConfig() Conf {
	return config
}

func readConfig() (err error) {
	var data []byte
	if data, err = ioutil.ReadFile("res/config_dev.yaml"); err != nil {
		err = errors.New(fmt.Sprintf("failed read config file: %s \n", err))
		return
	}
	if err = yaml.Unmarshal(data, &config); err != nil {
		err = errors.New(fmt.Sprintf("failed unmarshal config file: %s \n", err))
	}
	return
}

func init() {
	if err := readConfig(); err != nil {
		panic(err)
	}
}
