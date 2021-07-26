package accountDao

import (
	"baby-fried-rice/internal/pkg/kit/etcd"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/rpc"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/privatemessage"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/accountDao/cache"
	"baby-fried-rice/internal/pkg/module/accountDao/config"
	"baby-fried-rice/internal/pkg/module/accountDao/db"
	"baby-fried-rice/internal/pkg/module/accountDao/log"
	"baby-fried-rice/internal/pkg/module/accountDao/service/application"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"net/http"
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
	if err := cache.InitRedisClient(conf.Redis.RedisUrl, conf.Redis.RedisPassword, conf.Redis.RedisDB); err != nil {
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

func HttpRun() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	mux := runtime.NewServeMux()
	// HTTP转grpc
	err := user.RegisterDaoUserHandlerFromEndpoint(ctx, mux, fmt.Sprintf("%v:%v", conf.Server.Addr, conf.Server.Port), opts)
	if err != nil {
		panic(err)
	}
	if err = http.ListenAndServe(fmt.Sprintf("%v:%v", conf.Server.Addr, conf.Server.HttpPort), mux); err != nil {
		panic(err)
	}
}

func Main() {
	go ServerRun()
	go HttpRun()
	log.Logger.Info("server run successful")
	select {
	case err := <-errChan:
		log.Logger.Error(err.Error())
	}
}
