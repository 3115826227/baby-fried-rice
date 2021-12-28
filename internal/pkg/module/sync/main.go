package sync

import (
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/module/sync/config"
	"baby-fried-rice/internal/pkg/module/sync/es"
	"baby-fried-rice/internal/pkg/module/sync/log"
)

var (
	conf       models.Conf
	serverName = "sync"
)

func init() {
	// 初始化配置文件并获取
	conf = config.GetConfig()
	// 初始化日志
	if err := log.InitLog(serverName, conf.Log.LogLevel, conf.Log.LogPath); err != nil {
		panic(err)
	}
	log.Logger.Info("log init successful")
}

func initTableSyncHandle(mysqlConf models.Mysql) {
	tbSync := es.NewTableSyncHandle(mysqlConf, log.Logger)
	//tbSync.AddSource(conf.Database.SubDatabase.ImDatabase)
	tbSync.AddSource(conf.Database.SubDatabase.SmsDatabase)
	tbSync.AddSource(conf.Database.SubDatabase.SpaceDatabase)
	tbSync.AddSource(conf.Database.SubDatabase.ShopDatabase)
	if err := tbSync.Run(); err != nil {
		panic(err)
	}
	defer tbSync.Close()
}

func Main() {
	initTableSyncHandle(conf.Database.SubDatabase.ImDatabase)
}
