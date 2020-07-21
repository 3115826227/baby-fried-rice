package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"shop/src/config"
	"shop/src/log"
)

var (
	MainDB *gorm.DB
)

func GetDB() *gorm.DB {
	return MainDB
}

func init() {
	var err error
	MainDB, err = gorm.Open("mysql", config.Config.MysqlUrl)
	if err != nil {
		panic(err)
	} else {
		MainDB.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false)
		MainDB.SingularTable(true)
	}
}

// 批量删除
func DeleteMulti(conditions [][]interface{}) error {
	var err error
	tx := MainDB.Debug().Begin()
	defer func() {
		if err != nil {
			log.Logger.Warn("delete beans failed", zap.String("err", err.Error()))
			tx.Rollback()
		}
	}()
	if err = tx.Error; err != nil {
		return err
	}
	for _, values := range conditions {
		if len(values) > 0 {
			if err = tx.Delete(values[0], values[1:]...).Error; err != nil {
				return err
			}
		}
	}

	return tx.Commit().Error
}

func CreateMulti(bean ...interface{}) error {
	var err error
	tx := MainDB.Begin()
	defer func() {
		if err != nil {
			log.Logger.Warn("insert beans failed", zap.String("err", err.Error()))
			tx.Rollback()
		}
	}()
	if err = tx.Error; err != nil {
		return err
	}

	for k := range bean {
		if err = tx.Create(bean[k]).Error; err != nil {
			log.Logger.Warn("insert beans failed", zap.String("err", err.Error()))
			return err
		}
	}

	return tx.Commit().Error
}
