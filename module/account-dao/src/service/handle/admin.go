package handle

import (
	"errors"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/config"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/log"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model"
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func AdminInit(c *gin.Context) {
	var err error
	var req model.ReqAdminInit
	if err = c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	var count int
	if err := db.DB.Debug().Model(&model.AccountAdmin{}).Where("school_id = ? and super = true", req.SchoolId).Count(&count).Error; err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}
	if count > 0 {
		err = errors.New("账号已分配")
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	schools := model.GetSchool(req.SchoolId)
	if len(schools) != 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	school := schools[0]

	var admin = new(model.AccountAdmin)
	admin.ID = GenerateID()
	admin.LoginName = req.LoginName
	admin.Super = true
	admin.Username = school.Name
	admin.Password = EncodePassword(config.DefaultAdminPassword)
	admin.EncodeType = config.DefaultUserEncryMd5
	admin.SchoolId = school.ID

	var now = time.Now()
	//添加默认角色
	var defaultRole = model.AdminRole{
		Name:     config.DefaultRoleName,
		SchoolId: admin.SchoolId,
	}
	defaultRole.CreatedAt = now
	defaultRole.UpdatedAt = now
	//添加管理员角色
	var adminRole = &model.AdminRole{
		Name:     config.AdminRoleName,
		SchoolId: admin.SchoolId,
	}
	adminRole.CreatedAt = now
	adminRole.UpdatedAt = now
	if err := db.DB.Create(&adminRole).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}
	if err := db.DB.Debug().Model(&model.AdminRole{}).Where(`name = ? and school_id = ?`, config.AdminRoleName, admin.SchoolId).Find(&adminRole).Error; err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}
	var relation = &model.AccountAdminRoleRelation{
		AdminId: admin.ID,
		RoleId:  adminRole.ID,
	}
	relation.CreatedAt = now
	relation.UpdatedAt = now

	var beans = make([]interface{}, 0)
	permissions := model.GetPermission()
	for _, p := range permissions {
		var count = 0
		if err := db.DB.Debug().Model(&model.AdminRolePermissionRelation{}).Where(`role_id = ? and permission_id = ?`, adminRole.ID, p.ID).Count(&count).Error; err != nil {
			log.Logger.Warn(err.Error())
			continue
		}
		if count == 0 {
			var relations = model.AdminRolePermissionRelation{
				RoleId:       adminRole.ID,
				PermissionId: p.ID,
			}
			relations.CreatedAt = now
			relations.UpdatedAt = now
			beans = append(beans, &relations)
		}
	}
	var organize = model.AccountSchoolOrganize{
		Label:    school.Name,
		ParentId: config.RootSchoolOrganizeId,
		SchoolId: school.ID,
		Status:   true,
		Count:    0,
	}
	organize.ID = GenerateID()
	organize.CreatedAt = time.Now()
	organize.UpdatedAt = time.Now()
	beans = append(beans, &organize)
	beans = append(beans, &admin)
	beans = append(beans, &relation)
	beans = append(beans, &defaultRole)
	if err := db.CreateMulti(beans...); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}

	SuccessResp(c, "", nil)
}

func AdminLogin(c *gin.Context) {
	var err error
	var req model.ReqPasswordAdminLogin
	if err = c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	req.LoginName = strings.TrimSpace(req.LoginName)
	req.Password = strings.TrimSpace(req.Password)
	req.Ip = c.GetHeader("IP")

	var admin = model.AccountAdmin{}
	err = db.DB.Debug().Where("login_name = ? and password = ? and school_id = ?", req.LoginName, EncodePassword(req.Password), req.SchoolId).Find(&admin).Error
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	roles := model.GetRoleByAdmin(admin.ID)
	roleIDs := make([]int, 0)
	for _, r := range roles {
		roleIDs = append(roleIDs, r.ID)
	}
	permissions := model.GetPermission(roleIDs...)
	var resp = &model.RespAdminLogin{
		Admin:       admin,
		Roles:       roles,
		Permissions: permissions,
	}

	go AdminLoginLogAdd(admin.ID, req.Ip, time.Now())

	SuccessResp(c, "", resp)
}

func AdminGet(c *gin.Context) {
	schoolId := c.Query("school_id")
	if schoolId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	SuccessResp(c, "", model.GetAdminBySchool(schoolId))
}
