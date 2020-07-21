package model

import (
	"github.com/3115826227/baby-fried-rice/module/user-account/src/log"
	"github.com/3115826227/baby-fried-rice/module/user-account/src/service/model/db"
)

func GetSchoolById(id string) (school School, err error) {
	if err = db.DB.Where("id = ?", id).Find(&school).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}
