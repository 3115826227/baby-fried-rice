package model

import (
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/log"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model/db"
)

func GetRoleBySchool(school ...int) (roles []AdminRole) {
	roles = make([]AdminRole, 0)
	var template = db.DB.Debug()
	if len(school) != 0 {
		template = template.Where("school_id in (?)", school)
	}
	if err := template.Find(&roles).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}

func GetRole(id ...int) (roles []AdminRole) {
	roles = make([]AdminRole, 0)
	var template = db.DB.Debug()
	if len(id) != 0 {
		template = template.Where("id in (?)", id)
	}
	if err := template.Find(&roles).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}

func GetRoleByAdmin(admin string) (roles []AdminRole) {
	sql := fmt.Sprintf(`select role.* from admin_role as role
		left join account_admin_role_relation as rel
		on role.id = rel.role_id
		where rel.admin_id = ?`)
	if err := db.DB.Debug().Raw(sql, admin).Scan(&roles).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}

func GetPermission(role ...int) (permissions []AdminPermission) {
	permissions = make([]AdminPermission, 0)
	if len(role) == 0 {
		if err := db.DB.Debug().Model(&AdminPermission{}).Scan(&permissions).Error; err != nil {
			log.Logger.Warn(err.Error())
		}
	} else {
		sql := fmt.Sprintf(`select p.* from admin_permission as p
left join admin_role_permission_relation as r
on p.id = r.permission_id
where r.role_id in (?)`)
		if err := db.DB.Debug().Raw(sql, role).Scan(&permissions).Error; err != nil {
			log.Logger.Warn(err.Error())
		}
	}
	return
}

func GetSchool(id ...string) (schools []School) {
	schools = make([]School, 0)
	var template = db.DB.Debug()
	if len(id) != 0 {
		template = template.Where("id in (?)", id)
	}
	if err := template.Find(&schools).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}

func GetRootLoginLog(id ...string) (logs []AccountRootLoginLog) {
	logs = make([]AccountRootLoginLog, 0)
	var template = db.DB.Debug()
	if len(id) != 0 {
		template = template.Where("root_id in (?)", id)
	}
	if err := template.Order("login_time desc").Find(&logs).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}

func GetRoot(id string) (root AccountRoot) {
	if err := db.DB.Debug().Where("id = ?", id).Find(&root).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}

func GetRoots(ids []string) (roots []AccountRoot) {
	if err := db.DB.Debug().Where("id in (?)", ids).Find(&roots).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}

func GetAdminLoginLog(id ...string) (logs []AccountAdminLoginLog) {
	logs = make([]AccountAdminLoginLog, 0)
	var template = db.DB.Debug()
	if len(id) != 0 {
		template = template.Where("admin_id in (?)", id)
	}
	if err := template.Order("login_time desc").Find(&logs).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}

func GetAdminBySchool(schoolId string) (admins []AccountAdmin) {
	admins = make([]AccountAdmin, 0)
	if err := db.DB.Debug().Where("school_id = ?", schoolId).Find(&admins).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}

func GetAdmin(id ...string) (admins []AccountAdmin) {
	admins = make([]AccountAdmin, 0)
	var template = db.DB.Debug()
	if len(id) != 0 {
		template = template.Where("id in (?)", id)
	}
	if err := template.Find(&admins).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}

func GetUserLoginLog(id ...string) (logs []AccountUserLoginLog) {
	logs = make([]AccountUserLoginLog, 0)
	var template = db.DB.Debug()
	if len(id) != 0 {
		template = template.Where("user_id in (?)", id)
	}
	if err := template.Order("login_time desc").Find(&logs).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}

func GetUser(id string) (user AccountUser) {
	if err := db.DB.Debug().Where("id = ?", id).Find(&user).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}

func GetUsers(id []string) (users []AccountUser) {
	if err := db.DB.Debug().Where("id in (?)", id).Find(&users).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}

func GetUserDetail(id ...string) (users []AccountUserDetail) {
	users = make([]AccountUserDetail, 0)
	var template = db.DB.Debug()
	if len(id) != 0 {
		template = template.Where("id in (?)", id)
	}
	if err := template.Find(&users).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}

func GetUserSchoolDetail(id ...string) (detail []AccountUserSchoolDetail) {
	detail = make([]AccountUserSchoolDetail, 0)
	var template = db.DB.Debug()
	if len(id) != 0 {
		template = template.Where("id in (?)", id)
	}
	if err := template.Find(&detail).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}
