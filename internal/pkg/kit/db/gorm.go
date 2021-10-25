package db

import (
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/log"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ClientDB struct {
	db *gorm.DB
	lc log.Logging
}

func NewClientDB(mysqlUrl string, lc log.Logging) (client interfaces.DB, err error) {
	db, err := gorm.Open(mysql.Open(mysqlUrl), &gorm.Config{})
	if err != nil {
		return
	}
	client = &ClientDB{db: db, lc: lc}
	return
}

func (client *ClientDB) GetDB() *gorm.DB {
	return client.db.Debug()
}

func (client *ClientDB) InitTables(dos ...interfaces.DataObject) (err error) {
	var tables = make([]interface{}, 0)
	for _, do := range dos {
		tables = append(tables, do)
	}
	return client.db.Debug().AutoMigrate(tables...)
}

// 添加
func (client *ClientDB) CreateObject(object interfaces.DataObject) (err error) {
	return client.db.Debug().Create(object).Error
}

// 删除
func (client *ClientDB) DeleteObject(object interfaces.DataObject) (err error) {
	return client.db.Debug().Delete(object).Error
}

// 获取结果
func (client *ClientDB) GetObject(query map[string]interface{}, object interfaces.DataObject) (err error) {
	template := client.db.Debug().Table(object.TableName())
	for key, value := range query {
		template = template.Where(fmt.Sprintf("%v = ?", key), value)
	}
	return template.First(object).Error
}

// 更新数据
func (client *ClientDB) UpdateObject(object interfaces.DataObject) (err error) {
	return client.db.Debug().Table(object.TableName()).Save(object).Error
}

// 判断是否存在
func (client *ClientDB) ExistObject(query map[string]interface{}, do interfaces.DataObject) (exist bool, err error) {
	var count int64
	template := client.db.Debug().Table(do.TableName())
	for key, value := range query {
		template = template.Where(fmt.Sprintf("%v = ?", key), value)
	}
	err = template.First(do).Count(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return
	}
	exist = true
	return
}

func (client *ClientDB) CreateMulti(bean ...interface{}) error {
	var err error
	tx := client.db.Debug().Begin()
	defer func() {
		if err != nil {
			client.lc.Error("insert beans failed")
			tx.Rollback()
		}
	}()
	if err = tx.Error; err != nil {
		return err
	}

	for k := range bean {
		if err = tx.Create(bean[k]).Error; err != nil {
			client.lc.Error("insert beans failed")
			return err
		}
	}

	return tx.Commit().Error
}
