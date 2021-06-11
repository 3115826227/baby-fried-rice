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
	user.GET("/query", handle.UserQueryHandle)
	user.PATCH("/detail", handle.UserDetailUpdateHandle)
	user.PATCH("/pwd", handle.UserPwdUpdateHandle)

	user.POST("/private_message", handle.SendPrivateMessageHandle)
	user.GET("/private_message", handle.PrivateMessagesHandle)
	user.GET("/private_message/detail", handle.PrivateMessageDetailHandle)
	user.PATCH("/private_message/status", handle.UpdatePrivateMessageStatusHandle)
	user.DELETE("/private_message", handle.DeletePrivateMessageHandle)
}
