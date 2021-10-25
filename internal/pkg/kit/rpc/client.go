package rpc

import (
	"baby-fried-rice/internal/pkg/kit/log"
	"context"
	"crypto/x509"
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"time"
)

type ClientGRPC struct {
	addr      string
	lc        log.Logging
	rpcClient *grpc.ClientConn
}

type invokeLog struct {
	Req      interface{} `json:"req"`
	Reply    interface{} `json:"reply"`
	Success  bool        `json:"success"`
	Method   string      `json:"method"`
	Duration string      `json:"duration"`
	Error    string      `json:"error"`
}

func (il invokeLog) ToString() string {
	data, _ := json.Marshal(il)
	return string(data)
}

func withClientInterceptor(lc log.Logging) grpc.DialOption {
	return grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		var il = invokeLog{
			Req:      req,
			Reply:    reply,
			Method:   method,
			Duration: time.Since(start).String(),
		}
		if err != nil {
			il.Error = err.Error()
		} else {
			il.Success = true
		}
		lc.Debug(il.ToString())
		return err
	})
}

func NewClientGRPC(addr string, lc log.Logging, cp *x509.CertPool) (client *ClientGRPC, err error) {
	client = &ClientGRPC{
		addr: addr,
		lc:   lc,
	}
	cred := credentials.NewClientTLSFromCert(cp, "www.eline.com")
	client.rpcClient, err = grpc.Dial(addr,
		grpc.WithTransportCredentials(cred),
		withClientInterceptor(lc))
	if err != nil {
		return
	}
	return
}

func (client *ClientGRPC) GetRpcClient() *grpc.ClientConn {
	return client.rpcClient
}
