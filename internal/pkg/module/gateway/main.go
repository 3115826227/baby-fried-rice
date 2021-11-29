package gateway

import (
	"baby-fried-rice/internal/pkg/kit/middleware"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/module/gateway/cache"
	"baby-fried-rice/internal/pkg/module/gateway/config"
	"baby-fried-rice/internal/pkg/module/gateway/log"
	"baby-fried-rice/internal/pkg/module/gateway/server"
	"baby-fried-rice/internal/pkg/module/gateway/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

var (
	conf models.Conf
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
}

func LoadTls() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     fmt.Sprintf("%v:%v", conf.Server.HTTPServer.Addr, conf.Server.HTTPServer.Port),
		}).Process(c.Writer, c.Request)
		if err != nil {
			//如果出现错误，请不要继续。
			panic(err)
			return
		}
		// 继续往下处理
		c.Next()
	}
}

func Main() {
	engine := gin.Default()

	engine.Use(middleware.Cors())
	service.Register(engine)

	engine.RunTLS(
		fmt.Sprintf("%v:%v", conf.Server.HTTPServer.Addr, conf.Server.HTTPServer.Port),
		conf.Server.HTTPServer.CertFile, conf.Server.HTTPServer.KeyFile)
}
