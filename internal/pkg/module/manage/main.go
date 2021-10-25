package manage

import (
	"baby-fried-rice/internal/pkg/kit/etcd"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/module/manage/cache"
	"baby-fried-rice/internal/pkg/module/manage/config"
	"baby-fried-rice/internal/pkg/module/manage/db"
	"baby-fried-rice/internal/pkg/module/manage/log"
	"baby-fried-rice/internal/pkg/module/manage/service"
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
	// 初始化数据库
	if err := db.InitDB(); err != nil {
		panic(err)
	}
	// 初始化缓存
	if err := cache.InitCache(conf.Cache.Redis.MainCache, log.Logger); err != nil {
		panic(err)
	}
	// 初始化缓存
	if err := cache.InitRedisClient(conf.Cache.Redis.MainCache); err != nil {
		panic(err)
	}
	log.Logger.Info("cache init successful")
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
}

func ServerRun() {
	engine := gin.Default()

	gin.SetMode(gin.ReleaseMode)
	//engine.Use(middleware.Cors())
	service.Register(engine)

	if err := engine.Run(fmt.Sprintf("%v:%v", conf.Server.HTTPServer.Addr, conf.Server.HTTPServer.Port)); err != nil {
		panic(err)
	}
}

func Main() {
	go ServerRun()
	log.Logger.Info("server run successful")
	select {
	case err := <-errChan:
		log.Logger.Error(err.Error())
	}
}
