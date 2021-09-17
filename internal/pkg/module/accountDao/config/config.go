package config

import (
	"baby-fried-rice/internal/pkg/kit/models"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	DefaultUserEncryMd5 = "md5"
)

var (
	config models.Conf
)

func GetConfig() models.Conf {
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
