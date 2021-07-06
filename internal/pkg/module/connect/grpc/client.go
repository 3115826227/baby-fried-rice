package grpc

import (
	"baby-fried-rice/internal/pkg/kit/rpc"
	"baby-fried-rice/internal/pkg/module/connect/config"
	"baby-fried-rice/internal/pkg/module/connect/log"
	"baby-fried-rice/internal/pkg/module/connect/server"
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
	var b []byte
	if b, err = ioutil.ReadFile(config.GetConfig().Rpc.Client.CertFile); err != nil {
		return
	}
	var cp *x509.CertPool
	if cp = x509.NewCertPool(); !cp.AppendCertsFromPEM(b) {
		return
	}
	var c *rpc.ClientGRPC
	if c, err = rpc.NewClientGRPC(addr, log.Logger, cp); err != nil {
		return
	}
	client := &Client{c: c}
	clientMp[serverName][addr] = client
	return nil
}
