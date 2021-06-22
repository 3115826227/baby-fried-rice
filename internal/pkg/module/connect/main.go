package connect

import (
	"baby-fried-rice/internal/pkg/kit/etcd"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/module/connect/config"
	"baby-fried-rice/internal/pkg/module/connect/log"
	"baby-fried-rice/internal/pkg/module/connect/server"
	"baby-fried-rice/internal/pkg/module/connect/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

var (
	conf    config.Conf
	errChan chan error
)

func init() {
	// 初始化配置文件并获取
	conf = config.GetConfig()
	// 初始化日志
	if err := log.InitLog(conf.Server.Name, conf.Log.LogLevel, conf.Log.LogPath); err != nil {
		panic(err)
	}
	if err := server.InitRegisterClient(conf.Etcd); err != nil {
		panic(err)
	}
	srv := etcd.NewServerETCD(conf.Etcd, log.Logger)
	if err := srv.Connect(); err != nil {
		panic(err)
	}
	var serverInfo = interfaces.RegisterServerInfo{
		Addr:         conf.Server.Register,
		ServerName:   conf.Server.Name,
		ServerSerial: conf.Server.Serial,
	}
	if err := srv.Register(serverInfo); err != nil {
		panic(err)
	}
	errChan = make(chan error, 1)
	go srv.HealthCheck(serverInfo, time.Duration(conf.HealthyRollTime), errChan)
}

func Main() {
	engine := gin.Default()
	service.Register(engine)
	if err := engine.Run(fmt.Sprintf("%v:%v", conf.Server.Addr, conf.Server.Port)); err != nil {
		log.Logger.Error(err.Error())
		return
	}
}
