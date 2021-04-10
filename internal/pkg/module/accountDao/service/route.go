package service

import (
	"baby-fried-rice/internal/pkg/module/accountDao/log"
	"baby-fried-rice/internal/pkg/module/accountDao/service/handle"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(engine *gin.Engine) {
	app := engine.Group("/dao/account")

	app.GET("/test/ping", func(c *gin.Context) {
		log.Logger.Debug("pong")
		c.JSON(http.StatusOK, "pong")
	})
	app.POST("/private_message", handle.SendPrivateMessageHandle)
	app.GET("/private_message", handle.PrivateMessageDetailHandle)
	app.GET("/private_messages", handle.PrivateMessagesHandle)
	app.DELETE("/private_message", handle.DeletePrivateMessageHandle)
	app.PATCH("/private_message/status", handle.UpdatePrivateMessageStatusHandle)
	app.POST("/root/login", handle.RootLoginHandle)

	app.GET("/root/users", handle.UsersHandle)

	app.GET("/root/log/root_login", handle.RootLoginLogsHandle)
	app.GET("/root/log/user_login", handle.UserLoginLogsHandle)
}
