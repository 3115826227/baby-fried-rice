package handle

import (
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/account/src/config"
	"github.com/3115826227/baby-fried-rice/module/account/src/log"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model/db"
	"github.com/gin-gonic/gin"
)

type Permission struct {
	Id       int          `json:"id"`
	Name     string       `json:"name"`
	Path     string       `json:"path"`
	Method   string       `json:"method"`
	Type     int          `json:"type"`
	Children []Permission `json:"children"`
}

func PermissionInit() {
	var permissions = make([]model.AdminPermission, 0)
	permissions = PermissionDFS(0, config.Permission.Permission, permissions)
	for _, p := range permissions {
		sql := fmt.Sprintf(`insert into admin_permission values(%v, '%v', '%v', '%v', %v, %v) ON DUPLICATE KEY UPDATE name = '%v',path = '%v',method = '%v',types = %v,parent_id = %v`,
			p.ID, p.Name, p.Path, p.Method, p.Types, p.ParentId, p.Name, p.Path, p.Method, p.Types, p.ParentId)
		if err := db.DB.Debug().Model(&model.AdminPermission{}).Exec(sql).Error; err != nil {
			log.Logger.Warn(err.Error())
			continue
		}
	}
}

func PermissionDFS(parentId int, data interface{}, permissions []model.AdminPermission) []model.AdminPermission {
	switch data.(type) {
	case map[interface{}]interface{}:
		var mp = make(map[string]interface{})
		for key, value := range data.(map[interface{}]interface{}) {
			mp[key.(string)] = value
		}
		var permission = &model.AdminPermission{
			ID:       mp["id"].(int),
			Name:     mp["name"].(string),
			Path:     mp["path"].(string),
			Method:   mp["method"].(string),
			Types:    mp["types"].(int),
			ParentId: parentId,
		}
		permissions = append(permissions, *permission)
		if mp["children"] != nil {
			permissions = PermissionDFS(permission.ID, mp["children"], permissions)
		}
	case []interface{}:
		for _, obj := range data.([]interface{}) {
			permissions = PermissionDFS(parentId, obj, permissions)
		}
	}
	return permissions
}

func PermissionGet(c *gin.Context) {

}
