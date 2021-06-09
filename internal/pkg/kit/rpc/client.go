package rpc

import (
	"baby-fried-rice/internal/pkg/kit/log"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ClientGRPC struct {
	addr      string
	lc        log.Logging
	rpcClient *grpc.ClientConn
}

func NewClientGRPC(addr string, lc log.Logging, cp *x509.CertPool) (client *ClientGRPC, err error) {
	client = &ClientGRPC{
		addr: addr,
		lc:   lc,
	}
	cred := credentials.NewClientTLSFromCert(cp, "www.eline.com")
	client.rpcClient, err = grpc.Dial(addr, grpc.WithTransportCredentials(cred))
	if err != nil {
		return
	}
	return
}

func (client *ClientGRPC) GetRpcClient() *grpc.ClientConn {
	return client.rpcClient
}
