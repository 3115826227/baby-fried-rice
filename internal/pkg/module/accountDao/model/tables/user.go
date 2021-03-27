package tables

import "time"

type AccountUser struct {
	CommonField

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

	AccountID string `gorm:"column:account_id" json:"account_id"`
	Username  string `gorm:"column:username" json:"username"`
	SchoolId  string `gorm:"column:school_id" json:"school_id"`
	Verify    bool   `gorm:"column:verify" json:"verify"`
	Birthday  string `gorm:"column:birthday" json:"birthday"`
	Gender    bool   `gorm:"column:gender" json:"gender"`
	Age       int    `gorm:"column:age" json:"age"`
	//HeadImgUrl string `gorm:"column:head_img_url"`
	Phone    string `gorm:"column:phone" json:"phone"`
	Wx       string `gorm:"column:wx" json:"wx"`
	QQ       string `gorm:"column:qq" json:"qq"`
	Addr     string `gorm:"column:addr" json:"addr"`
	Hometown string `gorm:"column:hometown" json:"hometown"`
	Ethnic   string `gorm:"column:ethnic" json:"ethnic"`
}

func (table *AccountUserDetail) TableName() string {
	return "baby_account_user_detail"
}

func (table *AccountUserDetail) Get() interface{} {
	return *table
}

type UserDetail struct {
	UserId    string `gorm:"column:user_id;unique"`
	AccountId string `gorm:"column:account_id;not null"`
	Username  string `gorm:"column:username;not null"`
}

func (table *UserDetail) TableName() string {
	return "baby_im_user_detail"
}

func (table *UserDetail) Get() interface{} {
	return *table
}

// 好友关系表
type UserFriendRelation struct {
	CommonField

	Friend       string `gorm:"column:friend;unique_index:idx_friend_ref_category_friend_origin"`
	FriendRemark string `gorm:"column:friend_remark"`
	Origin       string `gorm:"column:origin;unique_index:idx_friend_ref_category_friend_origin"`
}

func (table *UserFriendRelation) TableName() string {
	return "baby_user_friend_relation"
}

func (table *UserFriendRelation) Get() interface{} {
	return *table
}

// 通知
type UserNotify struct {
	CommonField
	ReceiveUser string `gorm:"column:receive_user"`
	Content     string `gorm:"column:content"`
	ReadStatus  int    `gorm:"column:read_status"` // 通知读取状态： 0-未读 1-已读
}

func (table *UserNotify) TableName() string {
	return "baby_user_notify"
}

func (table *UserNotify) Get() interface{} {
	return *table
}
