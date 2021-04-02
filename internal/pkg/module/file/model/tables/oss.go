package tables

import "baby-fried-rice/internal/pkg/kit/models/tables"

type OssMeta struct {
	tables.CommonIntField
	Domain    string `gorm:"column:domain" json:"domain"`
	Bucket    string `gorm:"column:bucket" json:"bucket"`
	Size      int64  `gorm:"column:size" json:"size"`
	SecretKey string `gorm:"column:secret_key"`
	AccessKey string `gorm:"column:access_key"`
}

func (table *OssMeta) TableName() string {
	return "baby_file_oss_meta"
}

func (table *OssMeta) Get() interface{} {
	return *table
}


