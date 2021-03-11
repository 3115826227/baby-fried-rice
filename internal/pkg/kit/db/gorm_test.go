package db

import (
	"baby-fried-rice/internal/pkg/kit/log"
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
	"testing"
)

const (
	mysqlUrl = "root:123456@tcp(127.0.0.1:23306)/baby?charset=utf8mb4&parseTime=True&loc=Local"
)

func TestNewClientDB(t *testing.T) {
	lc, err := log.NewLoggerClient("kit-db", log.DebugLog, "")
	if err != nil {
		panic(err)
	}
	lc.Debug("log new success")
	client, err := NewClientDB(mysqlUrl, lc)
	if err != nil {
		panic(err)
	}
	lc.Debug("db client create success")
	if err = client.InitTables(&tables.AccountUser{}); err != nil {
		panic(err)
	}
	lc.Debug("db tables create success")
}
