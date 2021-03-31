package model

type ReqPasswordLogin struct {
	LoginName string `json:"login_name" binding:"required"` // 账号名称
	Password  string `json:"password" binding:"required"`   // 密码
	Ip        string `json:"ip"`
}

type ReqUserRegister struct {
	ReqPasswordLogin
	Username string `json:"username" binding:"required"` //昵称
	Gender   bool   `json:"gender" binding:"required"`   //性别
	Phone    string `json:"phone" binding:"required"`    //手机号
}

type ReqUserSendPrivateMessage struct {
}
