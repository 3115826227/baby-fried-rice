package service

import (
	"baby-fried-rice/internal/pkg/kit/middleware"
	"baby-fried-rice/internal/pkg/module/userAccount/service/handle"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	engine.POST("/api/user/register", handle.UserRegisterHandle)
	engine.POST("/api/user/login", handle.UserLoginHandle)
	root := engine.Group("/api/account/user", middleware.SetUserMeta())

	root.GET("/logout", handle.UserLogout)

}
