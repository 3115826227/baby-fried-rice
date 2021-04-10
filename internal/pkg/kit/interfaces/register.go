package interfaces

import "time"

type RegisterClient interface {
	Connect() error
	GetServer(severName string) (addr string, err error)
	GetServers(serverName string) ([]string, error)
	Close() error
}

type RegisterServerInfo struct {
	Addr         string `json:"addr"`
	ServerName   string `json:"server_name"`
	ServerSerial int    `json:"server_serial"`
}

type RegisterServer interface {
	Connect() error
	Register(info RegisterServerInfo) error
	HealthCheck(rs RegisterServerInfo, rollTime time.Duration, errChan chan error)
	Close() error
}
