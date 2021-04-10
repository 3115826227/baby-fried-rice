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
