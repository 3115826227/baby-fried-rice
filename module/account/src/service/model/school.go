package model

import (
	"github.com/3115826227/baby-fried-rice/module/account/src/log"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model/db"
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

func GetSchoolById(id string) (school School, err error) {
	if err = db.DB.Where("id = ?", id).Find(&school).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}
