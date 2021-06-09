package service

import (
	"baby-fried-rice/internal/pkg/kit/middleware"
	"baby-fried-rice/internal/pkg/module/userAccount/service/handle"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	engine.POST("/api/user/register", handle.UserRegisterHandle)
	engine.POST("/api/user/login", handle.UserLoginHandle)
	user := engine.Group("/api/account/user", middleware.SetUserMeta())

	user.GET("/logout", handle.UserLogoutHandle)
	user.GET("/detail", handle.UserDetailHandle)
	user.PATCH("/detail", handle.UserDetailUpdateHandle)
	user.PATCH("/pwd", handle.UserPwdUpdateHandle)
}
