package db

import (
	"github.com/3115826227/baby-fried-rice/module/account/config"
	"github.com/3115826227/baby-fried-rice/module/account/log"
	"github.com/3115826227/baby-fried-rice/module/account/service/model"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
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
		new(model.AccountRoot),
		new(model.AccountUser),
		new(model.AccountClient),
		new(model.School),
		new(model.ClientSchoolRelation),
		new(model.Area),
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
