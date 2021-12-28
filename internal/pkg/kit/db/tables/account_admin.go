package tables

import (
	"time"
)

type AccountAdmin struct {
	CommonField

	LoginName  string `gorm:"column:login_name;type:varchar(255);" json:"login_name"`
	Username   string `json:"username"`
	Password   string `json:"-"`
	EncodeType string `json:"-"`
	HeadImgUrl string `json:"head_img_url"`
	Phone      string `json:"phone"`
	ReqId      string `gorm:"column:req_id;type:varchar(255);" json:"req_id"`
}

func (table *AccountAdmin) TableName() string {
	return "account_admin"
}

func (table *AccountAdmin) Parse(rows []interface{}) {
	table.CommonField.ID = rows[0].(string)
	table.LoginName = rows[3].(string)
	table.Username = rows[4].(string)
}

type AccountAdminLoginLog struct {
	ID         int       `gorm:"column:id;AUTO_INCREMENT" json:"id"`
	AccountId  string    `json:"account_id"`
	IP         string    `json:"ip"`
	LoginCount int       `json:"login_count"`
	LoginTime  time.Time `gorm:"column:login_time;type:timestamp" json:"login_time"`
}

func (table *AccountAdminLoginLog) TableName() string {
	return "account_admin_login_log"
}

func (table *AccountAdminLoginLog) Get() interface{} {
	return *table
}
