package model

import (
	"github.com/3115826227/baby-fried-rice/module/account/src/log"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model/db"
)

func GetPermission() (permissions []AdminPermission) {
	permissions = make([]AdminPermission, 0)
	if err := db.DB.Model(&AdminPermission{}).Scan(&permissions).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}
