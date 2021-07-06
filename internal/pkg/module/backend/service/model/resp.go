package model

import "baby-fried-rice/internal/pkg/kit/handle"

type RspDaoRootLogin struct {
	Code int64 `json:"code"`
	Data struct {
		CreatedAt  string `json:"created_at"`
		EncodeType string `json:"encode_type"`
		ID         string `json:"id"`
		LoginName  string `json:"login_name"`
		Name       string `json:"name"`
		Password   string `json:"password"`
		ReqID      string `json:"req_id"`
		SchoolID   string `json:"school_id"`
		Super      bool   `json:"super"`
		UpdatedAt  string `json:"updated_at"`
		Username   string `json:"username"`
	} `json:"data"`
	Message string `json:"message"`
}

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
