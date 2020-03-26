package handle

import (
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/account/src/config"
	"github.com/3115826227/baby-fried-rice/module/account/src/log"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
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
	role := c.Query("role_id")
	var permissions = make([]int, 0)
	var relations = make([]model.AdminRolePermissionRelation, 0)
	if err := db.DB.Debug().Model(&model.AdminRolePermissionRelation{}).Where("role_id = ?", role).Find(&relations).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	for _, r := range relations {
		permissions = append(permissions, r.PermissionId)
	}
	var rsp = make([]model.AdminPermission, 0)
	if err := db.DB.Debug().Model(&model.AdminPermission{}).Where("id in (?)", permissions).Find(&rsp).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	SuccessResp(c, "", rsp)
}

func PermissionAllGet(c *gin.Context) {

	var permissions = make([]model.AdminPermission, 0)
	if err := db.DB.Debug().Model(&model.AdminPermission{}).Find(&permissions).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	var rsp = make([]model.RspAdminPermissions, 0)
	var mp = make(map[int]model.RspAdminPermissions)
	for _, p := range permissions {
		mp[p.ID] = model.RspAdminPermissions{
			Id:       p.ID,
			Name:     p.Name,
			Method:   p.Method,
			Path:     p.Path,
			Types:    p.Types,
			ParentId: p.ParentId,
			Children: make([]model.RspAdminPermissions, 0),
		}
	}
	for _, p := range mp {
		if p.ParentId == 0 {
			DFSGetAdminPermission(&p, mp)
			rsp = append(rsp, p)
		}
	}

	SuccessResp(c, "", rsp)
}

func DFSGetAdminPermission(permission *model.RspAdminPermissions, mp map[int]model.RspAdminPermissions) {
	for _, p := range mp {
		if p.ParentId == permission.Id {
			DFSGetAdminPermission(&p, mp)
			permission.Children = append(permission.Children, p)
		}
	}
}
