package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
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

	TokenSecret string `json:"token_secret"`
	Mysqls      struct {
		Account string `json:"account"`
		Shop    string `json:"shop"`
		Sms     string `json:"sms"`
	} `json:"mysqls"`
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
	var ok bool
	config.TokenSecret, ok = os.LookupEnv("TOKEN_SECRET")
	if !ok {
		config.TokenSecret = "baby"
	}
}
