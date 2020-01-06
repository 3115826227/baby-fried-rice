package model

import (
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model/db"
	"github.com/3115826227/baby-fried-rice/module/im/src/log"
)

func FindCertification(identify, name string) (certify SchoolUserCertification, err error) {
	if err = db.DB.Where("identify = ? and name = ?", identify, name).Find(&certify).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}

func GetCertifications(id ...string) (certifications []SchoolUserCertification, err error) {
	if err = db.DB.Where("id in (?)", id).Find(&certifications).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}
