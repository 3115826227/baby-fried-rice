package handle

import (
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/account/src/config"
	"github.com/3115826227/baby-fried-rice/module/account/src/log"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func RoleInit() {
	roleName := "管理员"
	now := time.Now().Format(config.TimeLayout)
	sql := fmt.Sprintf(`insert into admin_role values(1,'%v', '%v', '%v') ON DUPLICATE KEY UPDATE name = '%v',update_time = '%v'`,
		now, now, roleName, roleName, now)
	if err := db.DB.Debug().Model(&model.AdminRole{}).Exec(sql).Error; err != nil {
		log.Logger.Warn(err.Error())
		return
	}
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
