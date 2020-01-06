package model

import (
	"github.com/3115826227/baby-fried-rice/module/account/src/log"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model/db"
	"github.com/jinzhu/gorm"
	"time"
)

func init() {
	Sync(db.DB)
}

func Sync(engine *gorm.DB) {
	err := engine.AutoMigrate(
		new(AccountRoot),
		new(AccountUser),
		new(AccountClient),
		new(AccountUserDetail),
		new(AccountUserSchoolDetail),
		new(School),
		new(SchoolDepartment),
		new(SchoolCommunity),
		new(SchoolUserCertification),
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
	Name     string `gorm:"column:name"`
	Province string `gorm:"column:province"`
	City     string `gorm:"column:city"`
}

func (table *School) TableName() string {
	return "account_school"
}

type SchoolDepartment struct {
	CommonField

	SchoolId string `gorm:"school_id"`
	Name     string `gorm:"name"`
	FullName string `gorm:"full_name"`
	IsLeaf   int    `gorm:"is_leaf"`
	ParentId string `gorm:"parent_id"`
}

func (table *SchoolDepartment) TableName() string {
	return "account_school_department"
}

type SchoolCommunity struct {
	CommonField

	SchoolId      string    `gorm:"school_id"`
	Name          string    `gorm:"name"`
	Origin        string    `gorm:"origin"`
	EstablishTime time.Time `gorm:"column:establish_time;type:timestamp with time zone"`
}

func (table *SchoolCommunity) TableName() string {
	return "account_school_community"
}

type SchoolUserCertification struct {
	CommonField

	Identify           string `gorm:"column:identify;type:varchar(255);unique"`
	Name               string `gorm:"column:name;type:varchar(255)"`
	SchoolDepartmentId string `gorm:"column:school_department_id;type:varchar(255)"`
}

func (table *SchoolUserCertification) TableName() string {
	return "account_school_user_certification"
}

type AccountUser struct {
	CommonField

	LoginName  string `gorm:"column:login_name;type:varchar(255);"`
	Password   string `gorm:"column:password;type:varchar(255);"`
	EncodeType string `gorm:"column:encode_type"`
}

type AccountUserDetail struct {
	CommonField

	Username   string `gorm:"column:username"`
	SchoolId   string `gorm:"column:school_id"`
	Verify     int    `gorm:"column:verify"`
	Birthday   string `gorm:"column:birthday"`
	Gender     int    `gorm:"column:gender"`
	Age        int    `gorm:"column:age"`
	HeadImgUrl string `gorm:"column:head_img_url"`
	Phone      string `gorm:"column:phone"`
	Wx         string `gorm:"column:wx"`
	QQ         string `gorm:"column:qq"`
	Addr       string `gorm:"column:addr"`
	Hometown   string `gorm:"column:hometown"`
	Ethnic     string `gorm:"column:ethnic"`
}

func (table *AccountUserDetail) TableName() string {
	return "account_user_detail"
}

type AccountUserSchoolDetail struct {
	CommonField

	Name               string `gorm:"column:name"`
	Identify           string `gorm:"column:identify"`
	SchoolDepartmentId string `gorm:"school_department_id"`
}

func (table *AccountUserSchoolDetail) TableName() string {
	return "account_user_school_detail"
}

type Area struct {
	Code       string
	Name       string
	ParentCode string
	Level      int
}

func (table *Area) TableName() string {
	return "area"
}
