package service

import (
	"baby-fried-rice/internal/pkg/kit/middleware"
	"baby-fried-rice/internal/pkg/module/userAccount/service/handle"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	InitBackend()

	engine.POST("/api/user/register", handle.UserRegisterHandle)
	engine.POST("/api/user/login", handle.UserLoginHandle)
	user := engine.Group("/api/account/user", middleware.SetUserMeta())

	user.GET("/logout", handle.UserLogoutHandle)
	user.GET("/detail", handle.UserDetailHandle)
	user.GET("/query", handle.UserQueryHandle)
	user.PATCH("/detail", handle.UserDetailUpdateHandle)
	user.PATCH("/pwd", handle.UserPwdUpdateHandle)

	// 私信模块
	user.POST("/private_message", handle.SendPrivateMessageHandle)
	user.GET("/private_message", handle.PrivateMessagesHandle)
	user.GET("/private_message/detail", handle.PrivateMessageDetailHandle)
	user.PATCH("/private_message/status", handle.UpdatePrivateMessageStatusHandle)
	user.DELETE("/private_message", handle.DeletePrivateMessageHandle)

	// 积分模块
	user.GET("/coin/log", handle.CoinLogHandle)
	user.DELETE("/coin/log", handle.DeleteCoinLogHandle)
	user.GET("/coin/rank", handle.CoinRankHandle)
	user.GET("/coin/rank/board", handle.CoinRankBoardHandle)

	// 签到模块
	user.GET("/sign_in", handle.SignInHandle)
	user.GET("/sign_in/log", handle.SignInLogHandle)
}
