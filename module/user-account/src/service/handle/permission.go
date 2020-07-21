package handle

import (
	"github.com/3115826227/baby-fried-rice/module/user-account/src/log"
	"github.com/3115826227/baby-fried-rice/module/user-account/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/user-account/src/service/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
)

type Permission struct {
	Id       int          `json:"id"`
	Name     string       `json:"name"`
	Path     string       `json:"path"`
	Method   string       `json:"method"`
	Type     int          `json:"type"`
	Children []Permission `json:"children"`
}

func PermissionAllGet(c *gin.Context) {

	var permissions = make([]model.AdminPermission, 0)
	if err := db.DB.Debug().Model(&model.AdminPermission{}).Find(&permissions).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	var rsp = make([]model.RspAdminPermission, 0)
	var mp = make(map[int]model.RspAdminPermission)
	for _, p := range permissions {
		mp[p.ID] = model.RspAdminPermission{
			Id:       p.ID,
			Name:     p.Name,
			Method:   p.Method,
			Path:     p.Path,
			Types:    p.Types,
			ParentId: p.ParentId,
			Children: make([]model.RspAdminPermission, 0),
		}
	}
	for _, p := range mp {
		if p.ParentId == 0 {
			DFSGetAdminPermission(&p, mp)
			rsp = append(rsp, p)
		}
	}

	sort.Sort(model.RspAdminPermissions(rsp))
	SuccessResp(c, "", rsp)
}

func DFSGetAdminPermission(permission *model.RspAdminPermission, mp map[int]model.RspAdminPermission) {
	for _, p := range mp {
		if p.ParentId == permission.Id {
			DFSGetAdminPermission(&p, mp)
			permission.Children = append(permission.Children, p)
		}
	}
}
