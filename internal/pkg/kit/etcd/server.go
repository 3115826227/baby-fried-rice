package etcd

import (
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/log"
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

const (
	DefaultLeaseTTL        = 5
	DefaultHealthyRollTime = time.Duration(2000)
)

type ServerETCD struct {
	addrs    []string
	username string
	password string
	kv       clientv3.KV
	client   *clientv3.Client
	ctx      context.Context
	lease    clientv3.Lease
	leaseID  clientv3.LeaseID
	leaseTTL int64
	lc       log.Logging
}

func NewRegisterServerInfo(addr, serverName string, serverSerial int) interfaces.RegisterServerInfo {
	return interfaces.RegisterServerInfo{
		Addr:         addr,
		ServerName:   serverName,
		ServerSerial: serverSerial,
	}
}

func NewServerETCD(addrs []string, lc log.Logging) interfaces.RegisterServer {
	return &ServerETCD{
		addrs:    addrs,
		ctx:      context.Background(),
		leaseTTL: DefaultLeaseTTL,
		lc:       lc,
	}
}

func (server *ServerETCD) Connect() (err error) {
	server.client, err = clientv3.New(clientv3.Config{
		Endpoints:   server.addrs,
		DialTimeout: 5 * time.Second,
		TLS:         nil,
		Username:    server.username,
		Password:    server.password,
	})
	if err != nil {
		return
	}
	server.kv = clientv3.NewKV(server.client)
	server.ctx = context.Background()
	return
}

func (server *ServerETCD) Register(rs interfaces.RegisterServerInfo) (err error) {
	if server.lease == nil {
		server.lease = clientv3.NewLease(server.client)
	}
	leaseResp, err := server.lease.Grant(server.ctx, server.leaseTTL)
	if err != nil {
		return
	}
	server.leaseID = leaseResp.ID
	_, err = server.kv.Put(server.ctx, fmt.Sprintf("%v-%v", rs.ServerName, rs.ServerSerial), rs.Addr, clientv3.WithLease(leaseResp.ID))
	return
}

func (server *ServerETCD) leaseKeepAlive(rs interfaces.RegisterServerInfo) (err error) {
	if server.lease == nil {
		server.lease = clientv3.NewLease(server.client)
	}
	server.lease.Revoke(server.ctx, server.leaseID)
	return server.Register(rs)
}

func (server *ServerETCD) HealthCheck(rs interfaces.RegisterServerInfo, rollTime time.Duration, errChan chan error) {
	if rollTime == 0 {
		rollTime = DefaultHealthyRollTime
	}
	var tick = time.NewTicker(rollTime * time.Millisecond)
	defer close(errChan)
	for {
		select {
		case <-server.ctx.Done():
			return
		case <-tick.C:
			err := server.leaseKeepAlive(rs)
			if err != nil {
				server.lc.Error(fmt.Sprintf("server lease keep alive error: %v", err.Error()))
				errChan <- err
				return
			}
			server.lc.Info("server lease keep alive successful")
		}
	}
}

func (server *ServerETCD) Close() (err error) {
	server.ctx, _ = context.WithCancel(server.ctx)
	return server.client.Close()
}
