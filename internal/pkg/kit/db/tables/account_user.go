package tables

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"time"
)

type AccountUser struct {
	CommonField

	AccountId  string `gorm:"column:account_id;unique_index:unq_idx_account_id"`
	LoginName  string `gorm:"column:login_name;type:varchar(255);index:idx_user_login_name"`
	Password   string `gorm:"column:password;type:varchar(255);"`
	EncodeType string `gorm:"column:encode_type"`
	// 冻结状态
	Freeze bool `gorm:"column:freeze"`
	// 注销状态
	Cancel bool `gorm:"column:cancel"`
}

func (table *AccountUser) TableName() string {
	return "baby_account_user"
}

type AccountUserLoginLog struct {
	ID         int       `gorm:"column:id;AUTO_INCREMENT" json:"id"`
	AccountId  string    `json:"account_id"`
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

	AccountID  string `gorm:"column:account_id;pk"`
	Username   string `gorm:"column:username"`
	Describe   string `gorm:"column:describe"`
	SchoolId   string `gorm:"column:school_id"`
	Birthday   string `gorm:"column:birthday"`
	Gender     int32  `gorm:"column:gender"`
	Age        int64  `gorm:"column:age"`
	HeadImgUrl string `gorm:"column:head_img_url"`
	Phone      string `gorm:"column:phone;unique_index:user_detail_phone_idx"`
	// 是否为官方账号
	IsOfficial bool `gorm:"column:is_official"`
}

func (table *AccountUserDetail) TableName() string {
	return "baby_account_user_detail"
}

type AccountUserPhone struct {
	CommonIntField
	Phone string `gorm:"column:phone;unique"`
}

func (table *AccountUserPhone) TableName() string {
	return "baby_account_user_phone"
}

// 用户积分
type AccountUserCoin struct {
	AccountID       string `gorm:"column:account_id;primaryKey" json:"account_id"`
	Coin            int64  `gorm:"column:coin" json:"coin"`
	CoinTotal       int64  `gorm:"column:coin_total" json:"coin_total"`
	UpdateTimestamp int64  `gorm:"column:update_timestamp" json:"update_timestamp"`
}

func (table *AccountUserCoin) TableName() string {
	return "baby_account_user_coin"
}

// 用户积分日志
type AccountUserCoinLog struct {
	ID        int64  `gorm:"column:id;AUTO_INCREMENT" json:"id"`
	AccountID string `gorm:"column:account_id" json:"account_id"`
	Coin      int64  `gorm:"column:coin" json:"coin"`
	/*
		积分类型：
			1、每日登陆奖励  2、签到
			101、商城消费
	*/
	CoinType  constant.CoinType `gorm:"column:coin_type" json:"coin_type"`
	Timestamp int64             `gorm:"column:timestamp" json:"timestamp"`
}

func (table *AccountUserCoinLog) TableName() string {
	return "baby_account_user_coin_log"
}

// 用户封禁日志
type AccountUserBanLog struct {
	CommonIntField
	// 被封禁的用户id
	AccountId string `gorm:"column:account_id"`
	// 封禁原因
	Result string `gorm:"column:result"`
	// 封禁时长
	BanTime string `gorm:"column:ban_time"`
	// 解封时间
	DismissTimestamp int64 `gorm:"column:dismiss_timestamp"`
}

func (table *AccountUserBanLog) TableName() string {
	return "baby_account_user_ban_log"
}

// 用户举报日志
type AccountUserTipLog struct {
	CommonIntField
	// 举报用户
	ReportAccountId string `gorm:"column:report_account_id"`
	// 被举报用户
	ReportedAccountId string `gorm:"column:reported_account_id"`
	// 描述
	Describe string `gorm:"column:describe"`
	// 日志状态 0-待审核 1-有效举报 2-无效举报
	Status int `gorm:"column:status"`
}

func (table *AccountUserTipLog) TableName() string {
	return "baby_account_user_tip_log"
}
