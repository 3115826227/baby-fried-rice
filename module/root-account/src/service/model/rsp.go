package model

import "time"

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

type RspDaoRootLoginLog struct {
	Code int64 `json:"code"`
	Data []struct {
		RootID    string `json:"root_id"`
		LoginName string `json:"login_name"`
		Username  string `json:"username"`
		Phone     string `json:"phone"`
		Count     int    `json:"count"`
		Ip        string `json:"ip"`
		Area      string `json:"area"`
		Time      string `json:"time"`
	} `json:"data"`
	Message string `json:"message"`
}

type RspDaoRootDetail struct {
	Code int `json:"code"`
	Data struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		LoginName string    `json:"login_name"`
		Username  string    `json:"username"`
		ReqID     string    `json:"req_id"`
	} `json:"data"`
	Message string `json:"message"`
}

type RspDaoAdminLoginLog struct {
	Code int64 `json:"code"`
	Data []struct {
		AdminID   string `json:"admin_id"`
		LoginName string `json:"login_name"`
		Username  string `json:"username"`
		Phone     string `json:"phone"`
		Count     int    `json:"count"`
		IP        string `json:"ip"`
		Area      string `json:"area"`
		Time      string `json:"time"`
		School    string `json:"school"`
	} `json:"data"`
	Message string `json:"message"`
}

type RspDaoUserLoginLog struct {
	Code int64 `json:"code"`
	Data []struct {
		UserID    string `json:"user_id"`
		LoginName string `json:"login_name"`
		Username  string `json:"username"`
		Phone     string `json:"phone"`
		Count     int    `json:"count"`
		Ip        string `json:"ip"`
		Area      string `json:"area"`
		Time      string `json:"time"`
	} `json:"data"`
	Message string `json:"message"`
}

type RespAdmin struct {
	Code int64 `json:"code"`
	Data []struct {
		CreatedAt string `json:"created_at"`
		ID        string `json:"id"`
		LoginName string `json:"login_name"`
		Name      string `json:"name"`
		ReqID     string `json:"req_id"`
		SchoolID  string `json:"school_id"`
		Super     bool   `json:"super"`
		UpdatedAt string `json:"updated_at"`
		Username  string `json:"username"`
	} `json:"data"`
	Message string `json:"message"`
}

type RespUser struct {
	Code int64 `json:"code"`
	Data []struct {
		Addr      string `json:"addr"`
		Age       int64  `json:"age"`
		Birthday  string `json:"birthday"`
		CreatedAt string `json:"created_at"`
		Ethnic    string `json:"ethnic"`
		Gender    bool   `json:"gender"`
		Hometown  string `json:"hometown"`
		ID        string `json:"id"`
		Phone     string `json:"phone"`
		Qq        string `json:"qq"`
		SchoolID  string `json:"school_id"`
		UpdatedAt string `json:"updated_at"`
		Username  string `json:"username"`
		Verify    bool   `json:"verify"`
		Wx        string `json:"wx"`
	} `json:"data"`
	Message string `json:"message"`
}

type RespSchool struct {
	Code int64 `json:"code"`
	Data []struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Province string `json:"province"`
		City     string `json:"city"`
	} `json:"data"`
	Message string `json:"message"`
}

type RspSuccess struct {
	Code int `json:"code"`
}

type RspOkResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RespSuccessData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type RspUserData struct {
	UserId    string `json:"user_id"`
	Username  string `json:"username"`
	LoginName string `json:"login_name"`
}

type RspLogin struct {
	RspSuccess
	Data LoginResult `json:"data"`
}

type LoginResult struct {
	UserInfo RspUserData `json:"user_info"`
	Token    string      `json:"token"`
}
