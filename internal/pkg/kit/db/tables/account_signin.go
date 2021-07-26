package tables

import "baby-fried-rice/internal/pkg/kit/constant"

type AccountUserSignInLog struct {
	// 签到用户
	AccountId string `gorm:"column:account_id;pk" json:"account_id"`
	// 签到奖励积分
	Coin constant.RewardCoinBySignedInType `gorm:"column:coin" json:"coin"`
	// 签到类型
	SignInType constant.SignInType `gorm:"column:sign_in_type" json:"sign_in_type"`
	// 签到时间
	Timestamp int64 `gorm:"column:timestamp" json:"timestamp"`
}

func (table *AccountUserSignInLog) TableName() string {
	return "baby_account_user_signin_log"
}
