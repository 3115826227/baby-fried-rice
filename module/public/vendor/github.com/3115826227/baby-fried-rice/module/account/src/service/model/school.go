package model

import (
	"github.com/3115826227/baby-fried-rice/module/account/src/log"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model/db"
)

func InitSchool() {
	var schools = []School{
		{ID: "1", Name: "邵阳学院", Province: "湖南省", City: "邵阳市"},
		{ID: "2", Name: "中南大学", Province: "湖南省", City: "长沙市"},
		{ID: "3", Name: "湖南大学", Province: "湖南省", City: "长沙市"},
		{ID: "4", Name: "长沙学院", Province: "湖南省", City: "长沙市"},
	}
	for _, school := range schools {
		err := InsertSchool(school)
		if err != nil {
			log.Logger.Warn(err.Error())
		}
	}
}

func InsertSchool(school School) (err error) {
	var count = 0
	if err = db.DB.Debug().Model(&School{}).Where("id = ?", school.ID).Count(&count).Error; err != nil {
		return
	}
	if count != 0 {
		return
	}
	return db.DB.Debug().Create(&school).Error
}

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
