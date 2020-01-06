package model

import (
	"github.com/3115826227/baby-fried-rice/module/account/src/log"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model/db"
	"gopkg.in/gin-gonic/gin.v1/json"
)

type UserMeta struct {
	//用户ID
	UserId string `json:"userId"`
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

func IsDuplicateLoginNameByUser(loginName string) bool {
	var users = make([]AccountUser, 0)
	var count = 0
	if err := db.DB.Where("login_name = ?", loginName).Find(&users).Count(&count).Error; err != nil {
		log.Logger.Warn(err.Error())
		return true
	}
	return count != 0
}

func GetUserDetail(id string) (user AccountUserDetail, err error) {
	if err = db.DB.Where("id = ?", id).Find(&user).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return user, err
}
