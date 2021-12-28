package server

import (
	"baby-fried-rice/internal/pkg/kit/etcd"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/log"
	"baby-fried-rice/internal/pkg/module/game/config"
)

var (
	client interfaces.RegisterClient
)

func GetRegisterClient() interfaces.RegisterClient {
	if client == nil {
		if err := InitRegisterClient(config.GetConfig().Register.ETCD.Cluster, log.Logger); err != nil {
			panic(err)
		}
	}
	return client
}

func InitRegisterClient(addr []string, lc log.Logging) (err error) {
	client = etcd.NewClientETCD(addr, lc)
	return client.Connect()
}