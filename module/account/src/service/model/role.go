package model

import (
	"github.com/3115826227/baby-fried-rice/module/account/src/log"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model/db"
)

func GetAdminDetail(id string) (admin AccountAdmin, err error) {
	if err = db.DB.Where("id = ?", id).Find(&admin).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}

func GetRoleByAdminId(admin string) (roles []AdminRole) {
	var relations = make([]AccountAdminRoleRelation, 0)
	if err := db.DB.Debug().Where("admin_id = ?", admin).Find(&relations).Error; err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	var ids = make([]int, 0)
	for _, r := range relations {
		ids = append(ids, r.RoleId)
	}
	if err := db.DB.Debug().Where("id in (?)", ids).Find(&roles).Error; err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	return
}
