package tables

import (
	"time"
)

type AccountUser struct {
	CommonField

	AccountId  string `gorm:"column:account_id;" json:"account_id"`
	LoginName  string `gorm:"column:login_name;type:varchar(255);" json:"login_name"`
	Password   string `gorm:"column:password;type:varchar(255);" json:"password"`
	EncodeType string `gorm:"column:encode_type" json:"encode_type"`
}

func (table *AccountUser) TableName() string {
	return "baby_account_user"
}

func (table *AccountUser) Get() interface{} {
	return *table
}

type AccountUserLoginLog struct {
	ID         int       `gorm:"column:id;AUTO_INCREMENT" json:"id"`
	UserID     string    `json:"user_id"`
	LoginCount int       `json:"login_count"`
	IP         string    `json:"ip"`
	LoginTime  time.Time `gorm:"column:login_time;type:timestamp" json:"login_time"`
}

func (table *AccountUserLoginLog) TableName() string {
	return "baby_account_user_login_log"
}

func (table *AccountUserLoginLog) Get() interface{} {
	return *table
}

type AccountUserDetail struct {
	CommonField

	AccountID  string `gorm:"column:account_id;pk" json:"account_id"`
	Username   string `gorm:"column:username" json:"username"`
	Describe   string `gorm:"column:describe" json:"describe"`
	SchoolId   string `gorm:"column:school_id" json:"school_id"`
	Verify     bool   `gorm:"column:verify" json:"verify"`
	Birthday   string `gorm:"column:birthday" json:"birthday"`
	Gender     bool   `gorm:"column:gender" json:"gender"`
	Age        int64  `gorm:"column:age" json:"age"`
	HeadImgUrl string `gorm:"column:head_img_url"`
	Phone      string `gorm:"column:phone" json:"phone"`
	Wx         string `gorm:"column:wx" json:"wx"`
	QQ         string `gorm:"column:qq" json:"qq"`
	Addr       string `gorm:"column:addr" json:"addr"`
	Hometown   string `gorm:"column:hometown" json:"hometown"`
	Ethnic     string `gorm:"column:ethnic" json:"ethnic"`
}

func (table *AccountUserDetail) TableName() string {
	return "baby_account_user_detail"
}

func (table *AccountUserDetail) Get() interface{} {
	return *table
}
