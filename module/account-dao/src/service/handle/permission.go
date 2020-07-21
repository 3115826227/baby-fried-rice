package handle

import (
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/config"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/log"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PermissionInit(c *gin.Context) {
	permissions := model.GetPermission()
	if len(permissions) != 0 {
		SuccessResp(c, "", nil)
		return
	}
	permissions = PermissionDFS(0, config.Permission.Permission, permissions)
	var beans = make([]interface{}, 0)
	for _, p := range permissions {
		beans = append(beans, &p)
	}
	if err := db.CreateMulti(beans...); err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	SuccessResp(c, "", nil)
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
	SuccessResp(c, "", model.GetPermission())
}
