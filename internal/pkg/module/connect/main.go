package connect

import (
	"baby-fried-rice/internal/pkg/kit/etcd"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/module/connect/cache"
	"baby-fried-rice/internal/pkg/module/connect/config"
	"baby-fried-rice/internal/pkg/module/connect/log"
	"baby-fried-rice/internal/pkg/module/connect/server"
	"baby-fried-rice/internal/pkg/module/connect/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

var (
	conf    models.Conf
	errChan chan error
)

func init() {
	// 初始化配置文件并获取
	conf = config.GetConfig()
	// 初始化日志
	if err := log.InitLog(conf.Server.HTTPServer.Name, conf.Log.LogLevel, conf.Log.LogPath); err != nil {
		panic(err)
	}
	// 初始化缓存
	if err := cache.InitCache(conf.Cache.Redis.MainCache, log.Logger); err != nil {
		panic(err)
	}
	if err := server.InitRegisterClient(conf.Register.ETCD.Cluster); err != nil {
		panic(err)
	}
	srv := etcd.NewServerETCD(conf.Register.ETCD.Cluster, log.Logger)
	if err := srv.Connect(); err != nil {
		panic(err)
	}
	var serverInfo = interfaces.RegisterServerInfo{
		Addr:         conf.Server.HTTPServer.Register,
		ServerName:   conf.Server.HTTPServer.Name,
		ServerSerial: conf.Server.HTTPServer.Serial,
	}
	if err := srv.Register(serverInfo); err != nil {
		panic(err)
	}
	errChan = make(chan error, 1)
	go srv.HealthCheck(serverInfo, time.Duration(conf.Register.HealthyRollTime), errChan)
}

func Main() {
	engine := gin.Default()
	service.Register(engine)
	if err := engine.Run(fmt.Sprintf("%v:%v", conf.Server.HTTPServer.Addr, conf.Server.HTTPServer.Port)); err != nil {
		log.Logger.Error(err.Error())
		return
	}
}
