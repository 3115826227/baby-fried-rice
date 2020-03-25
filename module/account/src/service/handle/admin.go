package handle

import (
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/account/src/config"
	"github.com/3115826227/baby-fried-rice/module/account/src/log"
	"github.com/3115826227/baby-fried-rice/module/account/src/redis"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func AdminLogin(c *gin.Context) {
	var err error
	var req model.ReqAdminLogin
	if err = c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	req.LoginName = strings.TrimSpace(req.LoginName)
	req.Password = strings.TrimSpace(req.Password)

	var admin = model.AccountAdmin{}
	err = db.DB.Debug().Where("login_name = ? and password = ? and school_id = ?", req.LoginName, EncodePassword(req.Password), req.SchoolId).Find(&admin).Error
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	token, err := GenerateToken(admin.ID, time.Now())
	if err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	var userInfo = model.RspUserData{
		UserId:    admin.ID,
		LoginName: admin.LoginName,
		Username:  admin.Username,
		SchoolId:  admin.SchoolId,
		IsSuper:   admin.Super,
	}
	roles := model.GetRoleByAdminId(admin.ID)
	roleIDs := make([]int, 0)
	for _, r := range roles {
		roleIDs = append(roleIDs, r.ID)
	}
	permissions := model.GetPermissionByRole(roleIDs)

	var loginResult = model.LoginResult{
		UserInfo:   userInfo,
		Token:      token,
		Role:       roles,
		Permission: permissions,
	}
	var result = model.RspLogin{
		RspSuccess: model.RspSuccess{Code: 0},
		Data:       loginResult,
	}

	var isSuper string
	if admin.Super {
		isSuper = "1"
	}
	var userMeta = &model.UserMeta{
		UserId:   admin.ID,
		SchoolId: admin.SchoolId,
		IsSuper:  isSuper,
	}

	redis.AddAccountToken(fmt.Sprintf("%v:%v", TokenPrefix, token), userMeta.ToString())

	c.JSON(http.StatusOK, result)
}

func SubAdminAdd(c *gin.Context) {
	var err error
	var req model.ReqAdminAdd
	if err = c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	userMeta := GetUserMeta(c)
	_, err = model.GetRootDetail(userMeta.UserId)
	if err != nil {
		ErrorResp(c, http.StatusBadRequest, ErrCodeAccountNotFound, ErrCodeM[ErrCodeAccountNotFound])
		return
	}

	var admin = new(model.AccountAdmin)
	admin.ID = GenerateID()
	admin.LoginName = req.LoginName
	admin.Password = EncodePassword(AdminPassword)
	admin.EncodeType = UserEncryMd5

	var beans = make([]interface{}, 0)
	beans = append(beans, &admin)
	if err := db.CreateMulti(beans...); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}

	c.JSON(http.StatusOK, model.RspOkResponse{})
}

func SubAdminUpdate(c *gin.Context) {

}

func SubAdminGet(c *gin.Context) {
	userMeta := GetUserMeta(c)

	var admins = make([]model.AccountAdmin, 0)
	if err := db.DB.Debug().Where("school_id = ?", userMeta.SchoolId).Find(&admins).Error; err != nil {
		log.Logger.Warn(err.Error())
		return
	}

	var rsp = make([]model.RspSubAdmin, 0)
	for _, ad := range admins {
		rsp = append(rsp, model.RspSubAdmin{
			Id:       ad.ID,
			Username: ad.Username,
			Name:     ad.Name,
		})
	}
	SuccessResp(c, "", rsp)
}

/*
	admin子账号删除
		删除账号、角色关联表
*/
func SubAdminDelete(c *gin.Context) {

}

/*
	初始化管理账号
*/
func InitAdmin(c *gin.Context) {
	var err error
	var req model.ReqAdminInit
	if err = c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	school, err := model.GetSchoolById(req.SchoolId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	userMeta := GetUserMeta(c)
	_, err = model.GetRootDetail(userMeta.UserId)
	if err != nil {
		ErrorResp(c, http.StatusBadRequest, ErrCodeAccountNotFound, ErrCodeM[ErrCodeAccountNotFound])
		return
	}

	var admin = new(model.AccountAdmin)
	admin.ID = GenerateID()
	admin.LoginName = req.LoginName
	admin.Super = true
	admin.Username = school.Name
	admin.Password = EncodePassword(AdminPassword)
	admin.EncodeType = UserEncryMd5
	admin.SchoolId = school.ID

	var now = time.Now().Format(config.TimeLayout)
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
	beans = append(beans, &admin)
	beans = append(beans, &relation)
	beans = append(beans, &defaultRole)
	if err := db.CreateMulti(beans...); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}

	c.JSON(http.StatusOK, model.RspOkResponse{})
}
