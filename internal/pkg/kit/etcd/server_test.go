package etcd

import (
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"testing"
	"time"
)

func TestNewServerETCD(t *testing.T) {
	rs := NewServerETCD([]string{"http://127.0.0.1:23791"}, nil)
	if err := rs.Connect(); err != nil {
		panic(err)
	}
	defer rs.Close()
	var info = interfaces.RegisterServerInfo{
		Addr:         "",
		ServerName:   "root-account",
		ServerSerial: 1,
	}
	if err := rs.Register(info); err != nil {
		panic(err)
	}
	errChan := make(chan error, 1)
	go rs.HealthCheck(info, time.Second, errChan)
	select {
	case <-errChan:
		return
	}
}
