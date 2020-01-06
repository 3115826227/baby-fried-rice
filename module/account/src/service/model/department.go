package model

import (
	"github.com/3115826227/baby-fried-rice/module/account/src/log"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model/db"
)

func FindDepartmentById(id string) (department SchoolDepartment, err error) {
	if err = db.DB.Where("id = ?", id).First(&department).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}

func GetDepartmentsBySchool(schoolId string) (departments []SchoolDepartment, err error) {
	if err = db.DB.Where("school_id = ?", schoolId).Find(&departments).Error; err != nil {
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
