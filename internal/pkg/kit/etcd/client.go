package etcd

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/log"
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type ClientETCD struct {
	address  []string
	username string
	password string
	kv       clientv3.KV
	client   *clientv3.Client
	ctx      context.Context
	cancel   context.CancelFunc
	lease    clientv3.Lease
	leaseID  clientv3.LeaseID

	lc log.Logging

	serversMapLock *sync.RWMutex
	serversMap     map[string][]string
}

func NewClientETCD(addr []string, lc log.Logging) interfaces.RegisterClient {
	ctx, cancel := context.WithCancel(context.Background())
	var client = &ClientETCD{
		ctx:            ctx,
		cancel:         cancel,
		address:        addr,
		lc:             lc,
		serversMapLock: &sync.RWMutex{},
		serversMap:     make(map[string][]string),
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
	go client.flush()
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

type serverInfo struct {
	serverName string
	servers    []string
}

func (client *ClientETCD) flush() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-client.ctx.Done():
			return
		case <-ticker.C:
			var serverList = make([]serverInfo, 0)
			client.serversMapLock.RLock()
			for serverName := range client.serversMap {
				serverList = append(serverList, serverInfo{serverName: serverName})
			}
			client.serversMapLock.RUnlock()
			for index, info := range serverList {
				servers, err := client.list(info.serverName)
				if err != nil {
					client.lc.Error(err.Error())
					continue
				}
				if len(servers) == 0 {
					err = fmt.Errorf("no %v server register", info.serverName)
					client.lc.Error(err.Error())
					continue
				}
				serverList[index].servers = servers
			}
			client.serversMapLock.Lock()
			for _, info := range serverList {
				client.serversMap[info.serverName] = info.servers
			}
			client.serversMapLock.Unlock()
			client.lc.Debug("server flush success")
		}
	}
}

func (client *ClientETCD) GetServer(serverName string) (string, error) {
	client.serversMapLock.RLock()
	servers, exist := client.serversMap[serverName]
	client.serversMapLock.RUnlock()
	if !exist {
		var err error
		servers, err = client.list(serverName)
		if err != nil {
			client.lc.Error(err.Error())
			return "", err
		}
		if len(servers) == 0 {
			err = fmt.Errorf("no %v server register", serverName)
			client.lc.Error(err.Error())
			return "", err
		}
		client.serversMapLock.Lock()
		client.serversMap[serverName] = servers
		client.serversMapLock.Unlock()
	}
	return servers[genRand(len(servers))], nil
}

func (client *ClientETCD) GetServers(serverName string) ([]string, error) {
	return client.list(serverName)
}

func (client *ClientETCD) Close() (err error) {
	client.cancel()
	return client.client.Close()
}

func newAuthPathKey(path string) string {
	return fmt.Sprintf("%v/%v", constant.AuthPathRegisterPrefix, path)
}

func (client *ClientETCD) AddAuthPath(path, name string) error {
	_, err := client.kv.Put(client.ctx, newAuthPathKey(path), name)
	return err
}

func (client *ClientETCD) IsAuthPathConfig(path string) (bool, error) {
	resp, err := client.kv.Get(client.ctx, newAuthPathKey(path))
	if err != nil {
		return false, err
	}
	if resp.Count != 1 {
		return false, nil
	}
	return true, nil
}

func (client *ClientETCD) GetAuthPath() (map[string]string, error) {
	var mp = make(map[string]string)
	resp, err := client.kv.Get(client.ctx, constant.AuthPathRegisterPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	for _, item := range resp.Kvs {
		var key = strings.Replace(string(item.Key), constant.AuthPathRegisterPrefix, "", 1)[1:]
		mp[key] = string(item.Value)
	}
	return mp, nil
}

func (client *ClientETCD) DelAuthPath(path string) error {
	_, err := client.kv.Delete(client.ctx, newAuthPathKey(path))
	return err
}
