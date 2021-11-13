package requests

import "baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"

type PasswordLoginReq struct {
	LoginName string `json:"login_name" binding:"required"` // 账号名称
	Password  string `json:"password" binding:"required"`   // 密码
	Ip        string `json:"ip"`
}

type UserRegisterReq struct {
	PasswordLoginReq
	Username string `json:"username" binding:"required"` //昵称
	Phone    string `json:"phone" binding:"required"`    //手机号
}

// 管理平台添加用户
type AddUserReq struct {
	UserRegisterReq
	// 是否为官方用户
	IsOfficial bool `json:"is_official"`
}

type UserDetailUpdateReq struct {
	HeadImgUrl string `json:"head_img_url"`
	Describe   string `json:"describe"`
	Username   string `json:"username"`
	Gender     int32  `json:"gender"`
	Phone      string `json:"phone"`
	Age        int64  `json:"age"`
}

type UserPwdUpdateReq struct {
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}

type UserCoinGiveawayReq struct {
	Coin int64    `json:"coin"`
	Ids  []string `json:"ids"`
}

type UserCommunicationAddReq struct {
	Title             string                 `json:"title"`
	CommunicationType user.CommunicationType `json:"communication_type"`
	Content           string                 `json:"content"`
	Images            []string               `json:"images"`
}
