package server

import (
	"baby-fried-rice/internal/pkg/kit/etcd"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/module/connect/config"
)

var (
	client interfaces.RegisterClient
)

func GetRegisterClient() interfaces.RegisterClient {
	if client == nil {
		if err := InitRegisterClient(config.GetConfig().Register.ETCD.Cluster); err != nil {
			panic(err)
		}
	}
	return client
}

func InitRegisterClient(addr []string) (err error) {
	client = etcd.NewClientETCD(addr)
	return client.Connect()
}
