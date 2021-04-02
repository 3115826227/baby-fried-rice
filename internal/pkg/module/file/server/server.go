package server

import (
	"baby-fried-rice/internal/pkg/kit/etcd"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/module/file/config"
	"baby-fried-rice/internal/pkg/module/file/log"
)

var (
	server interfaces.RegisterServer
)

func GetRegisterServer() interfaces.RegisterServer {
	if server == nil {
		if err := InitRegisterServer(config.GetConfig().Etcd); err != nil {
			panic(err)
		}
	}
	return server
}

func InitRegisterServer(addrs []string) (err error) {
	server = etcd.NewServerETCD(addrs, log.Logger)
	return server.Connect()
}
