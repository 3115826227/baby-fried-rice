package model

type ReqLogin struct {
	LoginName string `json:"login_name" binding:"required"` // 账号名称
	Password  string `json:"password" binding:"required"`   // 密码
	Ip        string `json:"ip"`
}

type ReqUserRegister struct {
	LoginName string `json:"login_name" binding:"required"` // 账号名称
	Password  string `json:"password" binding:"required"`   // 密码
	Username  string `json:"username" binding:"required"`   //用户名
	Gender    bool   `json:"gender" binding:"required"`     //性别
	Phone     string `json:"phone" binding:"required"`      //手机号
}
