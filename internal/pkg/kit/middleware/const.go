package middleware

import "encoding/json"

const (
	HeaderUUID  = "requestID"
	HeaderToken = "token"

	HeaderUserId   = "userId"
	HeaderUsername = "username"
	HeaderSchoolId = "schoolId"
	HeaderPlatform = "platform"
	HeaderReqId    = "reqId"
	HeaderIsSuper  = "isSuper"

	TokenPrefix = "token"
)

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

func (meta *UserMeta) ToString() string {
	data, _ := json.Marshal(meta)
	return string(data)
}
