package requests

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/errors"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
)

// 用户账号密码登录
type PasswordLoginReq struct {
	LoginName string `json:"login_name" binding:"required"` // 账号名称
	Password  string `json:"password" binding:"required"`   // 密码
	Ip        string `json:"ip"`
}

func (req PasswordLoginReq) Validate() error {
	if req.LoginName == "" {
		return errors.NewCommonError(constant.CodeLoginNameEmpty)
	}
	if req.Password == "" {
		return errors.NewCommonError(constant.CodePasswordEmpty)
	}
	return nil
}

// 用户注册
type UserRegisterReq struct {
	PasswordLoginReq
	Username string `json:"username" binding:"required"` //昵称
}

func (req UserRegisterReq) Validate() error {
	if err := req.PasswordLoginReq.Validate(); err != nil {
		return err
	}
	if req.Username == "" {
		return errors.NewCommonError(constant.CodeUsernameEmpty)
	}
	return nil
}

// 管理平台添加用户
type AddUserReq struct {
	UserRegisterReq
	// 是否为官方用户
	IsOfficial bool `json:"is_official"`
}

// 用户详情更新
type UserDetailUpdateReq struct {
	HeadImgUrl string `json:"head_img_url"`
	Describe   string `json:"describe"`
	Username   string `json:"username"`
	Gender     int32  `json:"gender"`
	Phone      string `json:"phone"`
	Age        int64  `json:"age"`
}

func (req UserDetailUpdateReq) Validate() error {
	return nil
}

type UserPwdUpdateReq struct {
	Password    string `json:"password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (req UserPwdUpdateReq) Validate() error {
	return nil
}

type UserPhoneVerifyReq struct {
	Code string `json:"code" binding:"required,len=4"`
}

type UserCoinGiveawayReq struct {
	Coin int64    `json:"coin"`
	Ids  []string `json:"ids"`
}

func (req UserCoinGiveawayReq) Validate() error {
	return nil
}

type UserCommunicationAddReq struct {
	Title             string                 `json:"title" binding:"required"`
	CommunicationType user.CommunicationType `json:"communication_type"`
	Content           string                 `json:"content" binding:"required"`
	Images            []string               `json:"images"`
}

func (req UserCommunicationAddReq) Validate() error {
	return nil
}
