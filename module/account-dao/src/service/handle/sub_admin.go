package handle

import (
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/config"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/log"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func SubAdminAdd(c *gin.Context) {
	var err error
	var req model.ReqSubAdminAdd
	if err = c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	var admin = new(model.AccountAdmin)
	admin.ID = GenerateID()
	admin.LoginName = req.LoginName
	admin.Password = EncodePassword(config.DefaultSubAdminPassword)
	admin.EncodeType = config.DefaultUserEncryMd5
	admin.ReqId = req.ReqId

	var beans = make([]interface{}, 0)
	beans = append(beans, &admin)
	if err := db.CreateMulti(beans...); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}

	SuccessResp(c, "", nil)
}

func SubAdminUpdateRole(c *gin.Context) {
	var err error
	var req model.ReqSubAdminUpRole
	if err = c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	var rel = &model.AccountAdminRoleRelation{
		RoleId:  req.Role,
		AdminId: req.Admin,
	}

	if req.Status {
		if err = db.DB.Debug().Model(&rel).Update("update_at", time.Now()).Error; err != nil {
			log.Logger.Warn(err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
			return
		}
	} else {
		if err = db.DB.Debug().Where("role_id = ? and admin_id = ?", req.Role, req.Admin).Delete(&model.AccountAdminRoleRelation{}).Error; err != nil {
			log.Logger.Warn(err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
			return
		}
	}

	SuccessResp(c, "", nil)
}

func SubAdminGet(c *gin.Context) {
	schoolId := c.Query("school_id")
	var admins = make([]model.AccountAdmin, 0)
	if err := db.DB.Debug().Where("school_id = ?", schoolId).Find(&admins).Error; err != nil {
		log.Logger.Warn(err.Error())
		return
	}

	SuccessResp(c, "", admins)
}
