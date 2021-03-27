package etcd

import (
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"context"
	"github.com/coreos/etcd/clientv3"
	"math/rand"
	"time"
)

type ClientETCD struct {
	address  []string
	username string
	password string
	kv       clientv3.KV
	client   *clientv3.Client
	ctx      context.Context
	lease    clientv3.Lease
	leaseID  clientv3.LeaseID
}

func NewClientETCD(addr []string) interfaces.RegisterClient {
	var client = &ClientETCD{
		ctx:     context.Background(),
		address: addr,
	}
	return client
}

func (client *ClientETCD) Connect() (err error) {
	client.client, err = clientv3.New(clientv3.Config{
		Endpoints:   client.address,
		DialTimeout: 5 * time.Second,
		TLS:         nil,
		Username:    client.username,
		Password:    client.password,
	})
	if err != nil {
		return
	}
	client.kv = clientv3.NewKV(client.client)
	client.ctx = context.Background()
	return
}

func (client *ClientETCD) list(prefix string) ([]string, error) {
	resp, err := client.kv.Get(client.ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	servers := make([]string, 0)
	for _, value := range resp.Kvs {
		if value != nil {
			servers = append(servers, string(value.Value))
		}
	}
	return servers, nil
}

func genRand(num int) int {
	return int(rand.Int31n(int32(num)))
}

func (client *ClientETCD) GetServer(serverName string) (string, error) {
	servers, err := client.list(serverName)
	if err != nil {
		return "", err
	}
	return servers[genRand(len(servers))], nil
}

func (client *ClientETCD) Close() (err error) {
	return client.client.Close()
}
