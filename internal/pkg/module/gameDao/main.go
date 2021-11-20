package gameDao

import (
	"baby-fried-rice/internal/pkg/kit/etcd"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/kit/rpc"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/game"
	"baby-fried-rice/internal/pkg/module/gameDao/cache"
	"baby-fried-rice/internal/pkg/module/gameDao/config"
	"baby-fried-rice/internal/pkg/module/gameDao/db"
	"baby-fried-rice/internal/pkg/module/gameDao/log"
	"baby-fried-rice/internal/pkg/module/gameDao/service/application"
	"crypto/tls"
	"fmt"
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
	if err := log.InitLog(conf.Server.RPCServer.Name, conf.Log.LogLevel, conf.Log.LogPath); err != nil {
		panic(err)
	}
	// 初始化数据库
	if err := db.InitDB(conf.Database.MainDatabase.GetMysqlUrl()); err != nil {
		panic(err)
	}
	// 初始化缓存
	if err := cache.InitCache(conf.Cache.Redis.MainCache, log.Logger); err != nil {
		panic(err)
	}
	srv := etcd.NewServerETCD(conf.Register.ETCD.Cluster, log.Logger)
	if err := srv.Connect(); err != nil {
		panic(err)
	}
	var serverInfo = interfaces.RegisterServerInfo{
		Addr:         conf.Server.RPCServer.Register,
		ServerName:   conf.Server.RPCServer.Name,
		ServerSerial: conf.Server.RPCServer.Serial,
	}
	if err := srv.Register(serverInfo); err != nil {
		panic(err)
	}
	errChan = make(chan error, 1)
	go srv.HealthCheck(serverInfo, time.Duration(conf.Register.HealthyRollTime), errChan)
}

func ServerRun() {
	cert, err := tls.LoadX509KeyPair(conf.Rpc.Cert.Server.ServerCertFile, conf.Rpc.Cert.Server.ServerKeyFile)
	if err != nil {
		panic(err)
	}
	svr := rpc.NewServerGRPC(fmt.Sprintf("%v:%v", conf.Server.RPCServer.Addr, conf.Server.RPCServer.Port),
		log.Logger, &cert)
	game.RegisterDaoGameServer(svr.GetRpcServer(), &application.GameService{})
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
