package db

import (
	"github.com/3115826227/baby-fried-rice/module/public/config"
	"github.com/3115826227/baby-fried-rice/module/public/log"
	"github.com/3115826227/baby-fried-rice/module/public/service/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB *gorm.DB
)

func init() {
	var err error
	DB, err = gorm.Open("mysql", config.Config.MysqlUrl)
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
		new(model.Subject),
		new(model.Grade),
		new(model.Course),
		new(model.MessageType),
		new(model.Tutor),
		new(model.Appointment),
		//new(model.Area),

		//new(model.Station),
		//new(model.TrainMeta),
	).Error
	if err != nil {
		log.Logger.Warn(err.Error())
	}
}

func CreateMulti(bean ...interface{}) error {
	var err error
	tx := DB.Begin()
	defer func() {
		if err != nil {
			log.Logger.Warn("insert beans failed")
			tx.Rollback()
		}
	}()
	if err = tx.Error; err != nil {
		return err
	}

	for k := range bean {
		if err = tx.Create(bean[k]).Error; err != nil {
			log.Logger.Warn("insert beans failed")
			return err
		}
	}

	return tx.Commit().Error
}
