package rsp

import (
	"baby-fried-rice/internal/pkg/kit/constant"
)

type User struct {
	AccountID  string `json:"account_id"`
	Username   string `json:"username"`
	HeadImgUrl string `json:"head_img_url"`
	Remark     string `json:"remark"`
}

type UserDataResp struct {
	UserId    string `json:"user_id"`
	Username  string `json:"username"`
	LoginName string `json:"login_name"`
}

type LoginResult struct {
	UserInfo UserDataResp `json:"user_info"`
	Token    string       `json:"token"`
}

// 用户信息
type UserDetailResp struct {
	AccountId  string `json:"account_id"`
	Describe   string `json:"describe"`
	HeadImgUrl string `json:"head_img_url"`
	Username   string `json:"username"`
	SchoolId   string `json:"school_id"`
	Gender     bool   `json:"gender"`
	Age        int64  `json:"age"`
	Phone      string `json:"phone,"`
	Coin       int64  `json:"coin"`
}

// 后台管理用户信息
type UserBackendResp struct {
	AccountId    string `json:"account_id"`
	HeadImgUrl   string `json:"head_img_url"`
	Username     string `json:"username"`
	SchoolId     string `json:"school_id"`
	Gender       bool   `json:"gender"`
	Age          int64  `json:"age"`
	Phone        string `json:"phone"`
	RegisterTime int64  `json:"register_time"`
}

// 后台管理用户信息列表
type UserBackendListResp struct {
	List     []UserBackendResp `json:"list"`
	Page     int64             `json:"page"`
	PageSize int64             `json:"page_size"`
	Total    int64             `json:"total"`
}

type UserSignInResp struct {
	Ok       bool   `json:"ok"`
	Describe string `json:"describe"`
	Coin     int64  `json:"coin"`
}

// 用户签到日志信息
type UserSignInLogResp struct {
	SignInType constant.SignInType `json:"sign_in_type"`
	Coin       int64               `json:"coin"`
	Timestamp  int64               `json:"timestamp"`
}

// 用户积分信息
type UserCoinResp struct {
	List     []UserCoin `json:"list"`
	Page     int64      `json:"page"`
	PageSize int64      `json:"page_size"`
	Total    int64      `json:"total"`
}

type UserCoin struct {
	User            User  `json:"user"`
	Coin            int64 `json:"coin"`
	CoinTotal       int64 `json:"coin_total"`
	UpdateTimestamp int64 `json:"update_timestamp"`
}

type UserCoinLogResp struct {
	List     []UserCoinLog `json:"list"`
	Page     int64         `json:"page"`
	PageSize int64         `json:"page_size"`
	Total    int64         `json:"total"`
}

// 用户积分日志信息
type UserCoinLog struct {
	// 积分记录id
	Id int64 `json:"id"`
	// 积分变动值
	Coin int64 `json:"coin"`
	// 积分使用类型
	CoinType constant.CoinType `json:"coin_type"`
	// 积分使用描述
	Describe string `json:"describe"`
	// 积分使用时间
	Timestamp int64 `json:"timestamp"`
}

// 用户积分排名信息
type UserCoinRank struct {
	User User `json:"user"`
	// 排名
	Rank int64 `json:"rank"`
	// 积分数
	Coin int64 `json:"coin"`
	// 相同积分用户数
	SameCoinUsers int64 `json:"same_coin_users"`
}

type UserCoinRankBoard struct {
	User User `json:"user"`
	// 排名
	Rank int64 `json:"rank"`
	// 积分数
	Coin int64 `json:"coin"`
	// 获取时间
	UpdateTimestamp int64 `json:"update_timestamp"`
}

type UserCoinRankBoardResp struct {
	// 用户积分排名列表
	List []UserCoinRankBoard `json:"list"`
	// 统计时间
	StatisticTimestamp int64 `json:"statistic_timestamp"`
}

// 用户登录日志列表
type UserLoginLogListResp struct {
	List     []UserLoginLogResp `json:"list"`
	Page     int64              `json:"page"`
	PageSize int64              `json:"page_size"`
	Total    int64              `json:"total"`
}

type UserLoginLogResp struct {
	ID             int    `json:"id"`
	User           User   `json:"user"`
	LoginCount     int    `json:"login_count"`
	IP             string `json:"ip"`
	LoginTimestamp int64  `json:"login_timestamp"`
}
