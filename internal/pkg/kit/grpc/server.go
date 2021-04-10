package grpc

import (
	"baby-fried-rice/internal/pkg/kit/log"
	"crypto/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

type ServerGRPC struct {
	addr      string
	lc        log.Logging
	rpcServer *grpc.Server
}

func NewServerGRPC(addr string, lc log.Logging, cert *tls.Certificate) *ServerGRPC {
	server := &ServerGRPC{
		addr: addr,
		lc:   lc,
	}
	if cert == nil {
		server.rpcServer = grpc.NewServer()
	} else {
		cred := credentials.NewServerTLSFromCert(cert)
		server.rpcServer = grpc.NewServer(grpc.Creds(cred))
	}
	return server
}

func (server *ServerGRPC) GetRpcServer() *grpc.Server {
	return server.rpcServer
}

func (server *ServerGRPC) Run() (err error) {
	lis, err := net.Listen("tcp", server.addr)
	if err != nil {
		server.lc.Error(err.Error())
		return
	}
	err = server.rpcServer.Serve(lis)
	if err != nil {
		server.lc.Error(err.Error())
		return
	}
	return
}
