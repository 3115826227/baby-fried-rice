package es

import (
	"baby-fried-rice/internal/pkg/kit/db"
	Log "baby-fried-rice/internal/pkg/kit/log"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/module/sync/config"
	"database/sql"
	"fmt"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"github.com/siddontang/go-mysql-elasticsearch/river"
	"strings"
)

type mysqlMessage struct {
	lc       Log.Logging
	conf     models.Mysql
	rr       *river.River
	rrConfig *river.Config
}

func (msg *mysqlMessage) GetTables() ([]string, error) {
	dbClient, err := db.NewClientDB(msg.conf.GetMysqlUrl(), msg.lc)
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	rows, err = dbClient.GetDB().Raw("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA=?;", msg.conf.DBName).Rows()
	if err != nil {
		return nil, err
	}
	var tableNames = make([]string, 0)
	for rows.Next() {
		var tableName string
		if err = rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tableNames = append(tableNames, tableName)
	}
	return tableNames, nil
}

type TableSync interface {
	Close()
	Run() error
	AddSource(mysqlConf models.Mysql)
}

type tableSyncHandle struct {
	//canal.DummyEventHandler
	mysqlMessage *mysqlMessage
}

func NewTableSyncHandle(mysqlConf models.Mysql, lc Log.Logging) TableSync {
	conf := config.GetConfig()
	var riverConfig = &river.Config{
		MyAddr:     fmt.Sprintf("%v:%v", mysqlConf.Host, mysqlConf.Port),
		MyUser:     mysqlConf.Username,
		MyPassword: mysqlConf.Password,
		ESAddr:     strings.Join(conf.ElasticSearch.Urls, ","),
		ESUser:     conf.ElasticSearch.Username,
		ESPassword: conf.ElasticSearch.Password,
		ServerID:   1001,
		StatAddr:   fmt.Sprintf("%v:%v", conf.Server.HTTPServer.Addr, conf.Server.HTTPServer.Port),
		StatPath:   "/metrics",
		Sources: []river.SourceConfig{
			{
				Schema: mysqlConf.DBName,
				Tables: []string{"*"},
			},
		},
		BulkSize: 128,
	}
	var msg = &mysqlMessage{
		conf:     mysqlConf,
		lc:       lc,
		rrConfig: riverConfig,
	}
	return &tableSyncHandle{
		mysqlMessage: msg,
	}
}

func (handle *tableSyncHandle) AddSource(mysqlConf models.Mysql) {
	sources := handle.mysqlMessage.rrConfig.Sources
	sources = append(sources, river.SourceConfig{
		Schema: mysqlConf.DBName,
		Tables: []string{"*"},
	})
	handle.mysqlMessage.rrConfig.Sources = sources
}

func (handle *tableSyncHandle) Run() error {
	rr, err := river.NewRiver(handle.mysqlMessage.rrConfig)
	if err != nil {
		return err
	}
	handle.mysqlMessage.rr = rr
	return handle.mysqlMessage.rr.Run()
}

func (handle *tableSyncHandle) Close() {
	handle.mysqlMessage.rr.Close()
}
