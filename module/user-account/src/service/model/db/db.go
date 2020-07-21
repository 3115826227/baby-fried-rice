package db

import (
	"github.com/3115826227/baby-fried-rice/module/user-account/src/config"
	"github.com/3115826227/baby-fried-rice/module/user-account/src/log"
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
