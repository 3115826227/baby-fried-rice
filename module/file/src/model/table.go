package model

import (
	"github.com/3115826227/baby-fried-rice/module/file/src/config"
	"github.com/3115826227/baby-fried-rice/module/file/src/log"
	"github.com/3115826227/baby-fried-rice/module/file/src/model/db"
	"github.com/jinzhu/gorm"
	"time"
)

var flag bool

func init() {
	Sync(db.DB)
	flag = false
}

func Sync(engine *gorm.DB) {
	err := engine.AutoMigrate(
		new(OssMeta),
	).Error
	if err != nil {
		log.Logger.Warn(err.Error())
	}
}

type CommonIntField struct {
	ID        int       `gorm:"column:id;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time `gorm:"column:create_time;type:timestamp" json:"-"`
	UpdatedAt time.Time `gorm:"column:update_time;type:timestamp" json:"update_time"`
}

type OssMeta struct {
	CommonIntField
	Domain string `json:"domain"`
	Bucket string `json:"bucket"`
	Size   int64  `json:"size"`
}

func (table *OssMeta) TableName() string {
	return "file_oss_meta"
}

type OssInfo struct {
	OssMeta

	SecretKey string
	AccessKey string
}

func OssInfoQuery() (info OssInfo) {
	defer func() {
		flag = !flag
	}()

	var meta OssMeta
	var id int
	if flag {
		id = 2
	} else {
		id = 1
	}
	for _, data := range config.Key.Key {
		if id == data.Id {
			info.SecretKey = data.SecretKey
			info.AccessKey = data.AccessKey
		}
	}
	if err := db.DB.Debug().Where("id = ?", id).Find(&meta).Error; err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	info.OssMeta = meta
	return
}

func OssInfoDataUpdate(info OssInfo, fileSize int64) {
	var updateMap = map[string]interface{}{
		"update_time": time.Now(),
		"size":        info.Size + fileSize,
	}
	if err := db.DB.Debug().Model(OssMeta{}).Where("id = ?", info.ID).Updates(updateMap).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
}
