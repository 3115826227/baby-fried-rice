package config

import "github.com/spf13/viper"

func ReadConfig(conf interface{}) (err error) {
	viper.SetConfigFile("./res/config.yaml") // 指定配置文件路径
	viper.SetConfigName("config")            // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")              // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./res/")            // 查找配置文件所在的路径
	err = viper.ReadInConfig()               // 查找并读取配置文件
	if err != nil {                          // 处理读取配置文件的错误
		err = errors.New(fmt.Sprintf("Fatal error config file: %s \n", err))
		return
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		err = errors.New(fmt.Sprintf("Fatal error config file: %s \n", err))
		return
	}
	return
}
