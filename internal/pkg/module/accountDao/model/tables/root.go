package tables

import "time"

type AccountRoot struct {
	CommonField

	LoginName  string `gorm:"column:login_name;type:varchar(255);" json:"login_name"`
	Username   string `json:"username"`
	Password   string `json:"-"`
	EncodeType string `json:"-"`
	Phone      string `json:"phone"`
	ReqId      string `gorm:"column:req_id;type:varchar(255);" json:"req_id"`
}

func (table *AccountRoot) TableName() string {
	return "account_root"
}

func (table *AccountRoot) Get() interface{} {
	return *table
}

type AccountRootLoginLog struct {
	ID         int       `gorm:"column:id;AUTO_INCREMENT" json:"id"`
	RootID     string    `json:"root_id"`
	IP         string    `json:"ip"`
	LoginCount int       `json:"login_count"`
	LoginTime  time.Time `gorm:"column:login_time;type:timestamp" json:"login_time"`
}

func (table *AccountRootLoginLog) TableName() string {
	return "account_root_login_log"
}

func (table *AccountRootLoginLog) Get() interface{} {
	return *table
}
