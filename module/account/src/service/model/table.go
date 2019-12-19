package model

import (
	"time"
	"github.com/jinzhu/gorm"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model/db"
	"github.com/3115826227/baby-fried-rice/module/account/src/log"
)

func init() {
	Sync(db.DB)
}

func Sync(engine *gorm.DB) {
	err := engine.AutoMigrate(
		new(AccountRoot),
		new(AccountUser),
		new(AccountClient),
		new(School),
		new(ClientSchoolRelation),
		new(Area),
	).Error
	if err != nil {
		log.Logger.Warn(err.Error())
	}
}

type CommonField struct {
	ID        string    `gorm:"column:id;type:char(36);primary_key;not null"`
	CreatedAt time.Time `gorm:"column:create_time;type:timestamp with time zone" json:"-"`
	UpdatedAt time.Time `gorm:"column:update_time;type:timestamp with time zone" json:"-"`
}

type AccountRoot struct {
	CommonField

	LoginName  string `gorm:"column:login_name;type:varchar(255);"`
	Username   string
	Password   string
	EncodeType string
	ReqId      string `gorm:"column:req_id;type:varchar(255);"`
}

type AccountClient struct {
	CommonField

	Name   string
	Origin string
	Status int
}

func (client *AccountClient) TableName() string {
	return "account_client"
}

type ClientSchoolRelation struct {
	SchoolId string `gorm:"column:school_id;not null"`
	ClientId string `gorm:"column:client_id;not null"`
}

func (relation *ClientSchoolRelation) TableName() string {
	return "account_rel_client_school"
}

type School struct {
	ID       string `gorm:"column:id"`
	Name     string
	Province string
	City     string
}

type AccountUser struct {
	CommonField

	LoginName  string `gorm:"column:login_name;type:varchar(255);"`
	Username   string `gorm:"column:username;type:varchar(255);"`
	SchoolID   string `gorm:"column:school_id"`
	Password   string `gorm:"column:password;type:varchar(255);"`
	EncodeType string
	Verify     int
}

type Area struct {
	Code       string
	Name       string
	ParentCode string
	Level      int
}
