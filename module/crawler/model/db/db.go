package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/3115826227/baby-fried-rice/module/crawler/config"
	"github.com/3115826227/baby-fried-rice/module/crawler/log"
	"github.com/3115826227/baby-fried-rice/module/crawler/model"
)

var (
	DB *gorm.DB
)

func init() {
	var err error
	DB, err = gorm.Open("postgres", config.Config.PostgresUrl)
	if err != nil {
		panic(err)
	} else {
		DB.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false)
		DB.SingularTable(true)
	}
	Sync(DB)
}

func Sync(engine *gorm.DB) {
	err := engine.AutoMigrate(
		new(model.Station),
		new(model.TrainMeta),
		new(model.TrainStationRelation),
		new(model.TrainStationSeatPrice),
	).Error
	if err != nil {
		log.Logger.Warn(err.Error())
	}
}
