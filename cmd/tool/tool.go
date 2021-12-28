package main

import (
	"baby-fried-rice/cmd/tool/proto/hello_http"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var (
	address = "127.0.0.1:50052"
)

// 定义helloService并实现约定的接口
type helloService struct{}

// HelloService Hello服务
var HelloService = helloService{}

// SayHello 实现Hello服务接口
func (h helloService) SayHello(ctx context.Context, req *hello_http.HelloHTTPRequest) (*hello_http.HelloHTTPResponse, ERROR) {
	resp := new(hello_http.HelloHTTPResponse)
	resp.Message = fmt.Sprintf("Hello %s.", req.Name)

	return resp, nil
}

func run() {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	// 实例化grpc Server
	s := grpc.NewServer()

	// 注册HelloService
	hello_http.RegisterHelloHTTPServer(s, HelloService)

	grpclog.Println("Listen on " + address)
	s.Serve(listen)
}

func main() {
	go run()
	// 1. 定义一个context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// grpc服务地址
	opts := []grpc.DialOption{grpc.WithInsecure()}

	mux := runtime.NewServeMux()
	// HTTP转grpc
	err := hello_http.RegisterHelloHTTPHandlerFromEndpoint(ctx, mux, address, opts)
	if err != nil {
		grpclog.Fatalf("Register handler err:%v\n", err)
	}

	grpclog.Println("HTTP Listen on 8080")
	http.ListenAndServe(":8080", mux)
}
