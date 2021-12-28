package etcd

import (
	"fmt"
	"testing"
)

func TestNewClientETCD(t *testing.T) {
	rc := NewClientETCD([]string{"http://127.0.0.1:23791"}, nil)
	if err := rc.Connect(); err != nil {
		panic(err)
	}
	defer rc.Close()
	addr, err := rc.GetServer("root-account")
	if err != nil {
		panic(err)
	}
	fmt.Println(addr)
}

func TestClientETCD_AddAuthPath(t *testing.T) {
	rc := NewClientETCD([]string{"http://127.0.0.1:23791"}, nil)
	if err := rc.Connect(); err != nil {
		panic(err)
	}
	defer rc.Close()
	if err := rc.AddAuthPath("put:/api/v1/space/comment", "添加评论"); err != nil {
		panic(err)
	}
	mp, err := rc.GetAuthPath()
	if err != nil {
		panic(err)
	}
	fmt.Println(mp)
	exist, err := rc.IsAuthPathConfig("put:/api/v1/space/comment")
	if err != nil {
		panic(err)
	}
	fmt.Println(exist)
}
