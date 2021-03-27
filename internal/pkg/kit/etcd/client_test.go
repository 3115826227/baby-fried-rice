package etcd

import (
	"fmt"
	"testing"
)

func TestNewClientETCD(t *testing.T) {
	rc := NewClientETCD([]string{"http://127.0.0.1:23791"})
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
