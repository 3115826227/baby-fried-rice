package rpc

import (
	"baby-fried-rice/internal/pkg/kit/log"
	"context"
	"crypto/tls"
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"time"
)

type ServerGRPC struct {
	addr      string
	lc        log.Logging
	rpcServer *grpc.Server
}

type handleLog struct {
	Req      interface{} `json:"req"`
	Reply    interface{} `json:"reply"`
	Success  bool        `json:"success"`
	Method   string      `json:"method"`
	Duration string      `json:"duration"`
	Error    string      `json:"error"`
}

func (hl handleLog) ToString() string {
	data, _ := json.Marshal(hl)
	return string(data)
}

func withServerInterceptor(lc log.Logging) grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		reply, err := handler(ctx, req)
		var hl = handleLog{
			Req:      req,
			Reply:    reply,
			Method:   info.FullMethod,
			Duration: time.Since(start).String(),
		}
		if err != nil {
			hl.Error = err.Error()
		} else {
			hl.Success = true
		}
		lc.Debug(hl.ToString())
		return reply, err
	})
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
		server.rpcServer = grpc.NewServer(grpc.Creds(cred), withServerInterceptor(lc))
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
