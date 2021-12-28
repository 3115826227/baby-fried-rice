package handle

import (
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

const (
	HeaderUUID  = "requestID"
	HeaderToken = "token"
	TokenPrefix = "token"
	HeaderIP    = "IP"

	HeaderAccountId  = "accountId"
	HeaderUsername   = "username"
	HeaderSchoolId   = "schoolId"
	HeaderPlatform   = "platform"
	HeaderReqId      = "reqId"
	HeaderIsOfficial = "isOfficial"
	HeaderPhone      = "phone"

	GinContextKeyUserMeta = "userMeta"

	QueryId           = "id"
	QueryAccountId    = "account_id"
	QueryLikeUsername = "username"
	QueryLikeName     = "name"

	QueryPage     = "page"
	QueryPageSize = "page_size"
)

type UserMeta struct {
	//用户ID
	AccountId string `json:"accountId"`
	//用户名
	Username string `json:"username"`
	//学校ID
	SchoolId string `json:"schoolId"`
	//请求ID
	ReqId string `json:"reqId"`
	//平台
	Platform string `json:"platform"`
	//是否为超级管理员
	IsOfficial bool `json:"isOfficial"`
	//手机号 如果为空表示未认证用户
	Phone string `json:"phone"`
}

func (meta *UserMeta) ToString() string {
	data, _ := json.Marshal(meta)
	return string(data)
}

func GetUserMeta(c *gin.Context) *UserMeta {
	return c.MustGet(GinContextKeyUserMeta).(*UserMeta)
}

func (meta *UserMeta) GetUser() rsp.User {
	var phoneVerify = false
	if meta.Phone != "" {
		phoneVerify = true
	}
	return rsp.User{
		AccountID:   meta.AccountId,
		Username:    meta.Username,
		IsOfficial:  meta.IsOfficial,
		PhoneVerify: phoneVerify,
	}
}
