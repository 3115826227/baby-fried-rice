package service

import (
	"baby-fried-rice/internal/pkg/module/accountDao/service/handle"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	app := engine.Group("/dao/account")

	app.POST("/user/register", handle.UserRegisterHandle)
	app.POST("/user/login", handle.UserLoginHandle)
	app.POST("/user/private_message", )
	app.POST("/root/login", handle.RootLoginHandle)

	app.GET("/root/users", handle.UsersHandle)

	app.GET("/root/log/root_login", handle.RootLoginLogsHandle)
	app.GET("/root/log/user_login", handle.UserLoginLogsHandle)
}
