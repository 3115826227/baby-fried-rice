package handle

import (
	"github.com/3115826227/baby-fried-rice/module/account/src/config"
	"github.com/3115826227/baby-fried-rice/module/account/src/log"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model/db"
	"time"
)

func AdminRelationInit() {
	permissions := model.GetPermission()
	var now = time.Now().Format(config.TimeLayout)
	var beans = make([]interface{}, 0)
	for _, p := range permissions {
		var count = 0
		if err := db.DB.Debug().Model(&model.AdminRolePermissionRelation{}).Where(`role_id = 1 and permission_id = ?`, p.ID).Count(&count).Error; err != nil {
			log.Logger.Warn(err.Error())
			continue
		}
		if count == 0 {
			var relations = &model.AdminRolePermissionRelation{
				RoleId:       1,
				PermissionId: p.ID,
			}
			relations.CreatedAt = now
			relations.UpdatedAt = now
			beans = append(beans, relations)
		}
	}
	if err := db.CreateMulti(beans...); err != nil {
		panic(err)
	}
}

func GetRelationForId(role int)  {
}
