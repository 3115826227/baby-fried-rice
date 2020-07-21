package service

import (
	"github.com/3115826227/baby-fried-rice/module/user-account/src/middleware"
	"github.com/3115826227/baby-fried-rice/module/user-account/src/service/handle"
	"github.com/gin-gonic/gin"
)

func RegisterRoute(engine *gin.Engine) {

	//用户注册
	engine.POST("/api/user/register", handle.UserRegister)
	//用户登陆
	engine.POST("/api/user/login", handle.UserLogin)

	public := engine.Group("/api/user/")
	//获取学校列表
	public.GET("/school", handle.SchoolGet)

	app := engine.Group("/api/account")

	app.Use(middleware.MiddlewareSetUserMeta())

	app.GET("/user/info", handle.UserInfoGet)
	app.GET("/user/detail", handle.UserDetail)
	//用户修改密码
	app.PATCH("/user/password", handle.UserPasswordUpdate)
	//用户更新详细信息
	app.PATCH("/user/detail", handle.UserUpdate)
	//用户退出登陆
	app.GET("/user/logout", handle.UserLogout)
	//刷新用户token
	app.GET("/user/refresh", handle.UserRefresh)
	//用户实名制认证
	app.POST("/user/verify", handle.UserVerify)
}
