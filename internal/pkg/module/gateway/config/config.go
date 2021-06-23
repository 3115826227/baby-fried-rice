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

	Server struct {
		Name string `json:"name"`
		Addr string `json:"addr"`
		Port int    `json:"port"`
	} `json:"server"`

	Redis struct {
		RedisUrl      string `json:"redis_url"`
		RedisPassword string `json:"redis_password"`
		RedisDB       int    `json:"redis_db"`
	} `json:"redis"`

	Etcd []string `json:"etcd"`

	Servers struct {
		UserAccountServer string `json:"user_account_server"`
		RootAccountServer string `json:"root_account_server"`
		SpaceServer       string `json:"space_server"`
		ConnectServer     string `json:"connect_server"`
		ImServer          string `json:"im_server"`
		FileServer        string `json:"file_server"`
	}
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
