package grpc

import (
	"baby-fried-rice/internal/pkg/kit/rpc"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/privatemessage"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/shop"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/userAccount/config"
	"baby-fried-rice/internal/pkg/module/userAccount/log"
	"baby-fried-rice/internal/pkg/module/userAccount/server"
	"crypto/x509"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
)

var (
	clientMp = make(map[string]map[string]*Client)
)

type Client struct {
	c *rpc.ClientGRPC
}

func GetClientGRPC(serverName string) (*rpc.ClientGRPC, error) {
	if _, ok := clientMp[serverName]; !ok {
		clientMp[serverName] = make(map[string]*Client)
	}
	addr, err := server.GetRegisterClient().GetServer(serverName)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to get server %v", serverName))
		return nil, err
	}
	addr = strings.Split(addr, "//")[1]
	if err = initClient(serverName, addr); err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to init rpc client %v %v", serverName, addr))
		return nil, err
	}
	return clientMp[serverName][addr].c, nil
}

func initClient(serverName, addr string) (err error) {
	if _, exist := clientMp[serverName][addr]; exist {
		return nil
	}
	b, err := ioutil.ReadFile(config.GetConfig().Rpc.Cert.Client.ClientCertFile)
	if err != nil {
		return
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return
	}
	c, err := rpc.NewClientGRPC(addr, log.Logger, cp)
	if err != nil {
		return
	}
	client := &Client{c: c}
	clientMp[serverName][addr] = client
	return nil
}

func GetUserClient() (user.DaoUserClient, error) {
	cli, err := GetClientGRPC(config.GetConfig().Rpc.SubServers.AccountDaoServer)
	if err != nil {
		log.Logger.Error(err.Error())
		return nil, err
	}
	return user.NewDaoUserClient(cli.GetRpcClient()), nil
}

func GetShopClient() (shop.DaoShopClient, error) {
	cli, err := GetClientGRPC(config.GetConfig().Rpc.SubServers.ShopDaoServer)
	if err != nil {
		log.Logger.Error(err.Error())
		return nil, err
	}
	return shop.NewDaoShopClient(cli.GetRpcClient()), nil
}

func GetPrivateMessageClient() (privatemessage.DaoPrivateMessageClient, error) {
	cli, err := GetClientGRPC(config.GetConfig().Rpc.SubServers.AccountDaoServer)
	if err != nil {
		log.Logger.Error(err.Error())
		return nil, err
	}
	return privatemessage.NewDaoPrivateMessageClient(cli.GetRpcClient()), nil
}

func GetImClient() (im.DaoImClient, error) {
	cli, err := GetClientGRPC(config.GetConfig().Rpc.SubServers.ImDaoServer)
	if err != nil {
		log.Logger.Error(err.Error())
		return nil, err
	}
	return im.NewDaoImClient(cli.GetRpcClient()), nil
}
