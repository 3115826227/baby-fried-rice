package rpc

import (
	"baby-fried-rice/internal/pkg/kit/log"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/common"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"context"
	"crypto/tls"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"testing"
)

type UserService struct {
}

func GenerateID() string {
	UUID := uuid.NewV4()
	return UUID.String()
}

/*
	输出输入的用户名和密码，返回随机生成的UUID
*/
func (service *UserService) UserDaoRegister(context context.Context, request *user.ReqUserRegister) (*common.CommonResponse, error) {
	username := request.Username
	password := request.Login.Password
	fmt.Println(username, password)
	return &common.CommonResponse{}, nil
}

/*
	输出输入的用户名和密码，返回随机生成的UUID
*/
func (service *UserService) UserDaoLogin(context context.Context, request *user.ReqPasswordLogin) (*user.RspDaoUserLogin, error) {
	loginName := request.LoginName
	password := request.Password
	fmt.Println(loginName, password)
	return &user.RspDaoUserLogin{User: &user.RspDaoUser{AccountId: GenerateID()}}, nil
}

func TestNewServerGRPC(t *testing.T) {
	lc, err := log.NewLoggerClient("rpc-server", log.DebugLog, "")
	if err != nil {
		panic(err)
	}
	cert, err := tls.LoadX509KeyPair("./cert/server.pem", "./cert/server.key")
	if err != nil {
		panic(err)
	}
	server := NewServerGRPC("localhost:18040", lc, &cert)
	user.RegisterDaoUserServer(server.GetRpcServer(), &UserService{})
	server.Run()
}
