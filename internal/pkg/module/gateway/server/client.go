package server

import (
	"baby-fried-rice/internal/pkg/kit/etcd"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/module/gateway/config"
)

var (
	client interfaces.RegisterClient
)

func GetRegisterClient() interfaces.RegisterClient {
	if client == nil {
		if err := InitRegisterClient(config.GetConfig().Etcd); err != nil {
			panic(err)
		}
	}
	return client
}

func InitRegisterClient(addr []string) (err error) {
	client = etcd.NewClientETCD(addr)
	return client.Connect()
}
