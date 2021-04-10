package grpc

import (
	"baby-fried-rice/internal/pkg/kit/grpc"
	"baby-fried-rice/internal/pkg/kit/log"
	"baby-fried-rice/internal/pkg/module/space/config"
	"baby-fried-rice/internal/pkg/module/space/server"
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
	c *grpc.ClientGRPC
}

func GetClientGRPC(serverName string) (*grpc.ClientGRPC, error) {
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
		err = errors.Wrap(err, fmt.Sprintf("failed to init grpc client %v %v", serverName, addr))
		return nil, err
	}
	return clientMp[serverName][addr].c, nil
}

func initClient(serverName, addr string) (err error) {
	if _, exist := clientMp[serverName][addr]; exist {
		return nil
	}
	b, err := ioutil.ReadFile(config.GetConfig().Rpc.Client.CertFile)
	if err != nil {
		return
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return
	}
	c, err := grpc.NewClientGRPC(addr, log.Logger, cp)
	if err != nil {
		return
	}
	client := &Client{c: c}
	clientMp[serverName][addr] = client
	return nil
}
