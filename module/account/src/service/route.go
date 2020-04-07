package service

import (
	"github.com/3115826227/baby-fried-rice/module/account/src/middlware"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/handle"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model"
	"github.com/gin-gonic/gin"
)

func init() {
	handle.PermissionInit()
	model.InitSchool()
	err := handle.RootAdd()
	if err != nil {
		panic(err)
	}
}

func RegisterRoute(engine *gin.Engine) {

	engine.POST("/api/root/login", handle.RootLogin)
	engine.POST("/api/admin/login", handle.AdminLogin)

	engine.POST("/api/user/register", handle.UserRegister)
	engine.POST("/api/user/login", handle.UserLogin)

	public := engine.Group("/api/user/")
	public.GET("/school", handle.SchoolGet)
	public.GET("/admin/permission", handle.PermissionAllGet)

	app := engine.Group("/api/account")

	app.Use(middlware.MiddlewareSetUserMeta())

	app.GET("/user", handle.UserDetail)
	app.GET("/user/info", handle.UserInfoGet)
	app.PATCH("/user/password", handle.UserPasswordUpdate)
	app.PATCH("/user", handle.UserUpdate)
	app.GET("/user/logout", handle.UserLogout)
	app.GET("/user/refresh", handle.UserRefresh)
	app.POST("/user/verify", handle.UserVerify)

	app.GET("/admin/role", handle.RoleGet)
	app.GET("/admin/permission", handle.PermissionGet)

	app.POST("/admin/sub", handle.SubAdminAdd)
	app.GET("/admin/sub", handle.SubAdminGet)
	app.PATCH("/admin/sub")
	app.DELETE("/admin/sub", )

	app.POST("/admin/school/organize", handle.OrganizeAdd)
	app.GET("/admin/school/organize", handle.OrganizeGet)
	app.PUT("/admin/school/organize", handle.OrganizeUpdate)
	app.PUT("/admin/school/organize/status", handle.OrganizeStatus)
	app.DELETE("/admin/school/organize", handle.OrganizeDelete)

	app.GET("/admin/school/student", handle.StudentGet)
	app.POST("/admin/school/student", handle.StudentAdd)

	app.POST("/school/department", handle.SchoolDepartmentAdd)
	app.PATCH("/school/department", handle.SchoolDepartmentUpdate)
	app.GET("/school/department", handle.SchoolDepartments)
	app.DELETE("/school/department", handle.SchoolDepartmentDelete)

	app.DELETE("/school/certification", handle.SchoolCertificationDelete)

	app.GET("/root", handle.RootDetail)
	app.POST("/root/admin/init", handle.InitAdmin)

	app.GET("/root/system/docker/image", handle.ImagesGet)
	app.GET("/root/system/docker/container", handle.ContainersGet)
	app.GET("/root/system/docker/stats", handle.StatsGet)

	app.POST("/client", handle.ClientAdd)
}
