package grpc

import (
	"baby-fried-rice/internal/pkg/kit/rpc"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/live"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/live/config"
	"baby-fried-rice/internal/pkg/module/live/log"
	"baby-fried-rice/internal/pkg/module/live/server"
	"crypto/x509"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
	"sync"
)

var (
	clientMap = make(map[string]map[string]*Client)
	locker    sync.RWMutex
)

type Client struct {
	c *rpc.ClientGRPC
}

func GetClientGRPC(serverName string) (*rpc.ClientGRPC, error) {
	locker.RLock()
	srvMap, ok := clientMap[serverName]
	if !ok {
		srvMap = make(map[string]*Client)
	}
	locker.RUnlock()
	addr, err := server.GetRegisterClient().GetServer(serverName)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to get server %v", serverName))
		return nil, err
	}
	addr = strings.Split(addr, "//")[1]
	client, exist := srvMap[addr]
	if exist {
		return client.c, nil
	}
	if client, err = initClient(addr); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to init rpc client %v %v", serverName, addr))
		return nil, err
	}
	srvMap[addr] = client
	locker.Lock()
	clientMap[serverName] = srvMap
	locker.Unlock()
	return client.c, nil
}

func initClient(addr string) (client *Client, err error) {
	b, err := ioutil.ReadFile(config.GetConfig().Rpc.Cert.Client.ClientCertFile)
	if err != nil {
		return
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return
	}
	var c *rpc.ClientGRPC
	c, err = rpc.NewClientGRPC(addr, log.Logger, cp)
	if err != nil {
		return
	}
	return &Client{c: c}, nil
}

func GetUserClient() (user.DaoUserClient, error) {
	cli, err := GetClientGRPC(config.GetConfig().Rpc.SubServers.AccountDaoServer)
	if err != nil {
		log.Logger.Error(err.Error())
		return nil, err
	}
	return user.NewDaoUserClient(cli.GetRpcClient()), nil
}

func GetLiveClient() (live.DaoLiveClient, error) {
	cli, err := GetClientGRPC(config.GetConfig().Rpc.SubServers.LiveDaoServer)
	if err != nil {
		log.Logger.Error(err.Error())
		return nil, err
	}
	return live.NewDaoLiveClient(cli.GetRpcClient()), nil
}
