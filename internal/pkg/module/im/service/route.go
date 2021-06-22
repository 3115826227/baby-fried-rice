package service

import (
	"baby-fried-rice/internal/pkg/kit/middleware"
	"baby-fried-rice/internal/pkg/module/im/service/handle"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	app := engine.Group("/api/im", middleware.SetUserMeta())

	app.POST("/session", handle.SessionAddHandle)
	app.PATCH("/session", handle.SessionUpdateHandle)
	app.GET("/session", handle.SessionQueryHandle)
	app.GET("/session/detail", handle.SessionDetailHandle)
	app.POST("/session/join", handle.SessionJoinHandle)
	app.GET("/session/leave", handle.SessionLeaveHandle)
	app.DELETE("/session", handle.SessionDeleteHandle)

	app.GET("/session/message", handle.SessionMessageQueryHandle)
	app.DELETE("/session/message/flush", handle.SessionMessageFlushHandle)
}
