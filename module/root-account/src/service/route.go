package service

import (
	"github.com/3115826227/baby-fried-rice/module/root-account/src/middleware"
	"github.com/3115826227/baby-fried-rice/module/root-account/src/service/handle"
	"github.com/gin-gonic/gin"
)

func RegisterRoute(engine *gin.Engine) {

	engine.POST("/api/root/login", handle.RootLogin)

	root := engine.Group("/api/account/root", middleware.MiddlewareSetUserMeta())

	//学校操作
	root.POST("/school", handle.SchoolAdd)
	root.GET("/school", handle.SchoolGet)

	//学校管理员账号管理
	root.POST("/school/admin/init", handle.AdminInit)
	root.GET("/school/admin", handle.AdminGet)

	//用户账号
	root.GET("/user", handle.UserGet)

	//超级管理账号
	root.GET("/detail", handle.RootDetail)
	root.GET("/logout", handle.RootLogout)

	//日志管理
	root.GET("/login/log", handle.RootLoginLog)
	root.GET("/admin/login/log", handle.AdminLoginLog)
	root.GET("/user/login/log", handle.UserLoginLog)
}
