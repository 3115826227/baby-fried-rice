package handle

import (
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/log"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func RoleAdd(c *gin.Context) {
	var err error
	var req model.ReqRoleAdd
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	tx := db.DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		} else {
			tx.Commit()
			SuccessResp(c, "", nil)
		}
	}()

	var now = time.Now()
	var role = &model.AdminRole{
		Name:     req.Name,
		SchoolId: req.SchoolId,
		Describe: req.Describe,
	}
	role.CreatedAt = now
	role.UpdatedAt = now
	if err = tx.Debug().Create(&role).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	for _, p := range req.Permissions {
		var rel = &model.AdminRolePermissionRelation{
			RoleId:       role.ID,
			PermissionId: p,
		}
		rel.CreatedAt = now
		rel.UpdatedAt = now
		if err = tx.Debug().Create(&rel).Error; err != nil {
			log.Logger.Warn(err.Error())
			return
		}
	}
}

func RoleGet(c *gin.Context) {
	schoolId := c.Query("school")
	admin := c.Query("admin")
	var roles = make([]model.AdminRole, 0)
	if schoolId != "" {
		roles = model.GetRoleByAdmin(admin)
	} else if admin != "" {
		school, err := strconv.Atoi(schoolId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
			return
		}
		roles = model.GetRoleBySchool(school)
	} else {
		roles = model.GetRoleBySchool()
	}
	SuccessResp(c, "", roles)
}

func RoleUpdatePermission(c *gin.Context) {
	var err error
	var req model.ReqRoleUpPermission
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	var rel = &model.AdminRolePermissionRelation{
		RoleId:       req.Role,
		PermissionId: req.Permission,
	}

	if req.Status {
		if err = db.DB.Debug().Model(&rel).Update("update_at", time.Now()).Error; err != nil {
			log.Logger.Warn(err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
			return
		}
	} else {
		if err = db.DB.Debug().Where("role_id = ? and permission_id = ?", req.Role, req.Permission).Delete(&model.AdminRolePermissionRelation{}).Error; err != nil {
			log.Logger.Warn(err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
			return
		}
	}

	SuccessResp(c, "", nil)
}

func RoleDelete(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	var roleId int
	var err error
	roleId, err = strconv.Atoi(id)
	if err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	tx := db.DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		} else {
			tx.Commit()
			SuccessResp(c, "", nil)
		}
	}()

	if err = tx.Debug().Where("role_id = ?", roleId).Delete(&model.AccountAdminRoleRelation{}).Error; err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	if err = tx.Debug().Where("role_id = ?", roleId).Delete(&model.AdminRolePermissionRelation{}).Error; err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	if err = tx.Debug().Where("id = ?", roleId).Delete(&model.AdminRole{}).Error; err != nil {
		log.Logger.Warn(err.Error())
		return
	}
}
