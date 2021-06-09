package model

import (
	"baby-fried-rice/internal/pkg/kit/handle"
)

type RspUserData struct {
	UserId    string `json:"user_id"`
	Username  string `json:"username"`
	LoginName string `json:"login_name"`
}

type RspLogin struct {
	handle.RspSuccess
	Data LoginResult `json:"data"`
}

type LoginResult struct {
	UserInfo RspUserData `json:"user_info"`
	Token    string      `json:"token"`
}

type RspUserDetail struct {
	AccountId  string `json:"account_id"`
	Describe   string `json:"describe"`
	HeadImgUrl string `json:"head_img_url"`
	Username   string `json:"username"`
	SchoolId   string `json:"school_id"`
	Gender     bool   `json:"gender"`
	Age        int64  `json:"age"`
	Phone      string `json:"phone"`
}
