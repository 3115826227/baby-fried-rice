package service

import (
	"github.com/3115826227/baby-fried-rice/module/account-dao/src/service/handle"
	"github.com/gin-gonic/gin"
)

func init() {
	//var id = handle.GenerateID()
	//var root = model.AccountRoot{
	//	LoginName:   "root",
	//	Username:    "root",
	//	Password:    handle.EncodePassword("root"),
	//	EncodeType:  "md5",
	//}
	//root.ID = id
	//if err := db.DB.Debug().Create(&root).Error; err != nil {
	//	log.Logger.Warn(err.Error())
	//	return
	//}
}

func RegisterRoute(engine *gin.Engine) {
	app := engine.Group("/dao/account")

	app.GET("/school", handle.SchoolGet)
	app.POST("/school", handle.SchoolAdd)

	app.POST("/user/register", handle.UserRegister)
	app.POST("/user/login", handle.UserLogin)
	app.POST("/admin/login", handle.AdminLogin)
	app.POST("/root/login", handle.RootLogin)

	app.GET("/user/login/log", handle.UserLoginLog)
	app.GET("/user/verify_detail", handle.UserVerifyDetail)
	app.GET("/user/detail", handle.UserDetail)
	app.GET("/user/school/detail", handle.UserSchoolDetail)
	app.POST("/user/verify", handle.UserVerify)
	app.POST("/user/update", handle.UserUpdate)
	app.GET("/user/info", handle.UserInfoHandle)

	app.GET("/admin/login/log", handle.AdminLoginLog)
	app.POST("/admin/role", handle.RoleAdd)
	app.GET("/admin/role", handle.RoleGet)
	app.POST("/admin/role/permission", handle.RoleUpdatePermission)
	app.DELETE("/admin/role", handle.RoleDelete)

	app.GET("/admin/permission/init", handle.PermissionInit)
	app.GET("/admin/permission", handle.PermissionGet)

	app.POST("/admin/student", handle.StudentAdd)
	app.GET("/admin/student", handle.StudentGet)

	app.POST("/admin/organize", handle.OrganizeAdd)
	app.PATCH("/admin/organize", handle.OrganizeUpdate)
	app.PATCH("/admin/organize/status", handle.OrganizeStatus)
	app.GET("/admin/organize", handle.OrganizeGetHandle)
	app.GET("/admin/organize/exist", handle.OrganizeExistHandle)
	app.DELETE("/admin/organize", handle.OrganizeDelete)

	app.GET("/admin", handle.AdminGet)

	app.POST("/admin/sub_admin", handle.SubAdminAdd)
	app.GET("/admin/sub_admin", handle.SubAdminGet)

	app.POST("/admin/sub_admin/role", handle.SubAdminUpdateRole)

	app.GET("/root/login/log", handle.RootLoginLog)
	app.POST("/root/admin/init", handle.AdminInit)
	app.GET("/root/detail", handle.RootDetail)
}
