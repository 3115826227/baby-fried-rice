package model

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"time"
)

type RspDaoUserLogin struct {
	Code int `json:"code"`
	Data struct {
		User struct {
			ID         string    `json:"id"`
			CreatedAt  time.Time `json:"created_at"`
			UpdatedAt  time.Time `json:"updated_at"`
			LoginName  string    `json:"login_name"`
			Password   string    `json:"password"`
			EncodeType string    `json:"encode_type"`
		} `json:"user"`
		Detail struct {
			ID        string    `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			AccountID string    `json:"account_id"`
			Username  string    `json:"username"`
			SchoolID  string    `json:"school_id"`
			Verify    bool      `json:"verify"`
			Biry      string    `json:"biry"`
			Gender    bool      `json:"gender"`
			Age       int       `json:"age"`
			Phone     string    `json:"phone"`
			Wx        string    `json:"wx"`
			Qq        string    `json:"qq"`
			Addr      string    `json:"addr"`
			Hometown  string    `json:"hometown"`
			Ethnic    string    `json:"ethnic"`
		} `json:"detail"`
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
