package model

import (
	"time"
)

type CommonField struct {
	ID        string    `gorm:"column:id;type:char(36);primary_key;not null"`
	CreatedAt time.Time `gorm:"column:create_time;type:timestamp" json:"-"`
	UpdatedAt time.Time `gorm:"column:update_time;type:timestamp" json:"-"`
}

type AccountRoot struct {
	CommonField

	LoginName  string `gorm:"column:login_name;type:varchar(255);"`
	Username   string
	Password   string
	EncodeType string
	ReqId      string `gorm:"column:req_id;type:varchar(255);"`
}

func (table *AccountRoot) TableName() string {
	return "account_root"
}
