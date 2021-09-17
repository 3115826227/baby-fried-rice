package shop

import (
	"baby-fried-rice/internal/pkg/kit/etcd"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/module/shop/config"
	"baby-fried-rice/internal/pkg/module/shop/log"
	"baby-fried-rice/internal/pkg/module/shop/server"
	"baby-fried-rice/internal/pkg/module/shop/service"
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
	log.Logger.Info("log init successful")
	//// 初始化缓存
	//if err := cache.InitCache(conf.Redis.RedisUrl, conf.Redis.RedisPassword, conf.Redis.RedisDB, log.Logger); err != nil {
	//	panic(err)
	//}
	//log.Logger.Info("cache init successful")
	// 初始化注册中心
	srv := etcd.NewServerETCD(conf.Register.ETCD.Cluster, log.Logger)
	if err := srv.Connect(); err != nil {
		panic(err)
	}
	log.Logger.Info("register server init successful")
	// 注册本地服务到注册中心
	var serverInfo = interfaces.RegisterServerInfo{
		Addr:         conf.Server.HTTPServer.Register,
		ServerName:   conf.Server.HTTPServer.Name,
		ServerSerial: conf.Server.HTTPServer.Serial,
	}
	if err := srv.Register(serverInfo); err != nil {
		panic(err)
	}
	log.Logger.Info("server register successful")
	errChan = make(chan error, 1)
	// 开启后台协程向注册中心发送心跳机制
	go srv.HealthCheck(serverInfo, time.Duration(conf.Register.HealthyRollTime), errChan)
	if err := server.InitRegisterClient(conf.Register.ETCD.Cluster); err != nil {
		panic(err)
	}
}

func ServerRun() {
	engine := gin.Default()

	gin.SetMode(gin.ReleaseMode)
	//engine.Use(middleware.Cors())
	service.Register(engine)

	engine.Run(fmt.Sprintf("%v:%v", conf.Server.HTTPServer.Addr, conf.Server.HTTPServer.Port))
}

func Main() {
	go ServerRun()
	log.Logger.Info("server run successful")
	select {
	case err := <-errChan:
		log.Logger.Error(err.Error())
	}
}
