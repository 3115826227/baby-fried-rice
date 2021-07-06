package db

import (
	"baby-fried-rice/internal/pkg/kit/db"
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/module/backend/config"
	"baby-fried-rice/internal/pkg/module/backend/log"
	"time"
)

var (
	accountClient interfaces.DB
)

func initRoot() {
	var root = tables.AccountRoot{
		LoginName:  "root",
		Username:   "后台管理账号",
		Password:   handle.EncodePassword("root1234"),
		EncodeType: "md5",
	}
	root.ID = "10000000"
	now := time.Now()
	root.CreatedAt, root.UpdatedAt = now, now
	if exist, err := accountClient.ExistObject(map[string]interface{}{"id": root.ID}, &tables.AccountRoot{}); err != nil {
		panic(err)
	} else if !exist {
		if err = accountClient.CreateObject(&root); err != nil {
			panic(err)
		}
	}
	return
}

func GetAccountDB() interfaces.DB {
	return accountClient
}

func InitDB() (err error) {
	var conf = config.GetConfig()
	accountClient, err = db.NewClientDB(conf.Mysqls.Account, log.Logger)
	if err != nil {
		return
	}
	initRoot()
	return
}
