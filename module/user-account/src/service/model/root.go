package model

import (
	"github.com/3115826227/baby-fried-rice/module/user-account/src/log"
	"github.com/3115826227/baby-fried-rice/module/user-account/src/service/model/db"
)

func GetRootDetail(id string) (root AccountRoot, err error) {
	if err = db.DB.Where("id = ?", id).Find(&root).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return root, err
}
