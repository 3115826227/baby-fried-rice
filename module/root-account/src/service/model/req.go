package model

type ReqLogin struct {
	LoginName string `json:"login_name" binding:"required"` // 账号名称
	Password  string `json:"password" binding:"required"`   // 密码
	Ip        string `json:"ip"`
}

type ReqSchoolAdd struct {
	Name     string `json:"name" binding:"required"`
	Province string `json:"province"`
	City     string `json:"city"`
}

type ReqAdminInit struct {
	LoginName string `json:"login_name" binding:"required"` // 账号登陆名称
	Name      string `json:"name" binding:"required"`       //账号名称
	SchoolId  string `json:"school_id" binding:"required"`  //学校id
}
