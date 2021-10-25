package db

import (
	"baby-fried-rice/internal/pkg/kit/db"
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/module/spaceDao/config"
	"baby-fried-rice/internal/pkg/module/spaceDao/log"
)

var (
	client interfaces.DB
)

func GetDB() interfaces.DB {
	if client == nil {
		if err := InitDB(config.GetConfig().Database.MainDatabase.GetMysqlUrl()); err != nil {
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
		&tables.Space{},
		&tables.SpaceDetail{},
		&tables.VisitedRelation{},
		&tables.OperatorRelation{},
		&tables.CommentRelation{},
		&tables.CommentDetail{},
	)
}
