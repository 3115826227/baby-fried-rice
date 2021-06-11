package accountDao

import (
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/rpc"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/privatemessage"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/accountDao/cache"
	"baby-fried-rice/internal/pkg/module/accountDao/config"
	"baby-fried-rice/internal/pkg/module/accountDao/db"
	"baby-fried-rice/internal/pkg/module/accountDao/log"
	"baby-fried-rice/internal/pkg/module/accountDao/server"
	"baby-fried-rice/internal/pkg/module/accountDao/service/application"
	"crypto/tls"
	"fmt"
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
	// 初始化数据库
	if err := db.InitDB(conf.MysqlUrl); err != nil {
		panic(err)
	}
	// 初始化缓存
	if err := cache.InitCache(conf.Redis.RedisUrl, conf.Redis.RedisPassword, conf.Redis.RedisDB, log.Logger); err != nil {
		panic(err)
	}
	if err := server.InitRegisterServer(conf.Etcd); err != nil {
		panic(err)
	}
	var serverInfo = interfaces.RegisterServerInfo{
		Addr:         conf.Server.Register,
		ServerName:   conf.Server.Name,
		ServerSerial: conf.Server.Serial,
	}
	if err := server.GetRegisterServer().Register(serverInfo); err != nil {
		panic(err)
	}
	errChan = make(chan error, 1)
	go server.GetRegisterServer().HealthCheck(serverInfo, time.Duration(conf.HealthyRollTime), errChan)
}

func ServerRun() {
	cert, err := tls.LoadX509KeyPair(conf.Rpc.Server.CertFile, conf.Rpc.Server.KeyFile)
	if err != nil {
		panic(err)
	}
	svr := rpc.NewServerGRPC(fmt.Sprintf("%v:%v", conf.Server.Addr, conf.Server.Port),
		log.Logger, &cert)
	user.RegisterDaoUserServer(svr.GetRpcServer(), &application.UserService{})
	privatemessage.RegisterDaoPrivateMessageServer(svr.GetRpcServer(), &application.PrivateMessageService{})
	if err = svr.Run(); err != nil {
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
