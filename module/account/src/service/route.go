package service

import (
	"github.com/3115826227/baby-fried-rice/module/account/src/middlware"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/handle"
	"github.com/gin-gonic/gin"
)

func init() {
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

	app := engine.Group("/api/account")

	app.Use(middlware.MiddlewareSetUserMeta())

	app.GET("/user", handle.UserDetail)
	app.PATCH("/user/password", handle.UserPasswordUpdate)
	app.PATCH("/user", handle.UserUpdate)
	app.GET("/user/logout", handle.UserLogout)
	app.GET("/user/refresh", handle.UserRefresh)
	app.POST("/user/verify", handle.UserVerify)

	app.POST("/school/department", handle.SchoolDepartmentAdd)
	app.PATCH("/school/department", handle.SchoolDepartmentUpdate)
	app.GET("/school/department", handle.SchoolDepartments)
	app.DELETE("/school/department", handle.SchoolDepartmentDelete)

	app.DELETE("/school/certification", handle.SchoolCertificationDelete)

	app.GET("/root", handle.RootDetail)
	app.POST("/root/admin", handle.AddAdmin)

	app.POST("/client", handle.ClientAdd)
}
