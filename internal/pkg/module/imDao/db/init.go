package db

import (
	"baby-fried-rice/internal/pkg/kit/db"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/module/imDao/config"
	"baby-fried-rice/internal/pkg/module/imDao/log"
	"baby-fried-rice/internal/pkg/module/imDao/model/tables"
)

var (
	client interfaces.DB
)

func GetDB() interfaces.DB {
	if client == nil {
		if err := InitDB(config.GetConfig().MysqlUrl); err != nil {
			panic(err)
		}
	}
	return client
}

func InitDB(mysqlUrl string) (err error) {
	client, err = db.NewClientDB(mysqlUrl, log.Logger)
	if err != nil {
		return
	}
	return client.InitTables(
		&tables.Session{},
		&tables.SessionUserRelation{},
		&tables.Message{},
		&tables.MessageUserRelation{},
	)
}
