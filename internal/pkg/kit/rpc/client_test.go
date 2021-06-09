package rpc

import (
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/kit/log"
	"context"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestNewClientGRPC(t *testing.T) {
	lc, err := log.NewLoggerClient("rpc-server", log.DebugLog, "")
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadFile("./cert/server.pem")
	if err != nil {
		panic(err)
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		panic(err)
	}
	client, err := NewClientGRPC("localhost:18040", lc, cp)
	if err != nil {
		panic(err)
	}
	c := user.NewDaoUserClient(client.GetRpcClient())
	var req = &user.ReqUserLogin{
		Password: "234",
	}
	resp, err := c.UserDaoLogin(context.Background(), req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
