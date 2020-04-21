package model

type UserMeta struct {
	//用户ID
	UserId string `json:"userId"`
	//用户名
	Username string `json:"username"`
	//学校ID
	SchoolId string `json:"schoolId"`
	//请求ID
	ReqId string `json:"reqId"`
	//平台
	Platform string `json:"platform"`
	//是否为超级管理员
	IsSuper string `json:"isSuper"`
}

