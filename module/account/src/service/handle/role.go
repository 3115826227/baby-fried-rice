package handle

import (
	"github.com/3115826227/baby-fried-rice/module/account/src/log"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RoleGet(c *gin.Context) {
	userMeta := GetUserMeta(c)

	var roles = make([]model.AdminRole, 0)
	if err := db.DB.Model(&model.AdminRole{}).Where("school_id = ?", userMeta.SchoolId).Find(&roles).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	SuccessResp(c, "", roles)
}

func RoleAdd(c *gin.Context) {
	var err error
	var req model.ReqRoleAdd
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
}
