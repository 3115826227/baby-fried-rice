package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"net/url"
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

	Connect struct {
		UserUrl   string `json:"user_url"`
		AdminUrl  string `json:"admin_url"`
		RootUrl   string `json:"root_url"`
		PublicUrl string `json:"public_url"`
		SquareUrl string `json:"square_url"`
		ImUrl     string `json:"im_url"`
	} `json:"connect"`

	Servers struct {
		UserAccountServer string `json:"user_account_server"`
		RootAccountServer string `json:"root_account_server"`
	}

	ParserUserUrl   *url.URL
	ParserAdminUrl  *url.URL
	ParserPublicUrl *url.URL
	ParserSquareUrl *url.URL
	ParserImUrl     *url.URL
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
	config.ParserUserUrl, err = url.Parse(config.Connect.UserUrl)
	if err != nil {
		panic("error user url")
	}
	config.ParserAdminUrl, err = url.Parse(config.Connect.AdminUrl)
	if err != nil {
		panic("error admin url")
	}
	config.ParserPublicUrl, err = url.Parse(config.Connect.PublicUrl)
	if err != nil {
		panic("error user url")
	}
	config.ParserSquareUrl, err = url.Parse(config.Connect.SquareUrl)
	if err != nil {
		panic("error user url")
	}
	config.ParserImUrl, err = url.Parse(config.Connect.ImUrl)
	if err != nil {
		panic("error user url")
	}

}
