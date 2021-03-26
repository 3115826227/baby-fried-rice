package model

type ReqLogin struct {
	LoginName string `json:"loginName" binding:"required"` // 账号名称
	Password  string `json:"password" binding:"required"`  // 密码
	Ip        string `json:"ip"`
}
