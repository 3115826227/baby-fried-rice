package server

import (
	"baby-fried-rice/internal/pkg/kit/etcd"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/log"
)

var (
	server interfaces.RegisterServer
)

func GetRegisterServer() interfaces.RegisterServer {
	return server
}

func InitRegisterServer(addrs []string, lc log.Logging) (err error) {
	server = etcd.NewServerETCD(addrs, lc)
	return server.Connect()
}
