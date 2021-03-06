package gateway

import (
	"baby-fried-rice/internal/pkg/kit/middleware"
	"baby-fried-rice/internal/pkg/module/gateway/cache"
	"baby-fried-rice/internal/pkg/module/gateway/config"
	"baby-fried-rice/internal/pkg/module/gateway/log"
	"baby-fried-rice/internal/pkg/module/gateway/server"
	"baby-fried-rice/internal/pkg/module/gateway/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

var (
	conf config.Conf
)

func init() {
	// 初始化配置文件并获取
	conf = config.GetConfig()
	// 初始化日志
	if err := log.InitLog(conf.Server.Name, conf.Log.LogLevel, conf.Log.LogPath); err != nil {
		panic(err)
	}
	// 初始化缓存
	if err := cache.InitCache(conf.Redis.RedisUrl, conf.Redis.RedisPassword, conf.Redis.RedisDB, log.Logger); err != nil {
		panic(err)
	}
	if err := server.InitRegisterClient(conf.Etcd); err != nil {
		panic(err)
	}
}

func Main() {
	engine := gin.Default()

	engine.Use(middleware.Cors())
	service.Register(engine)

	engine.Run(fmt.Sprintf("%v:%v", conf.Server.Addr, conf.Server.Port))
}
