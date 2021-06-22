package service

import (
	"baby-fried-rice/internal/pkg/kit/middleware"
	"baby-fried-rice/internal/pkg/module/connect/service/handle"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	handle.Init()

	app := engine.Group("/api/connect", middleware.SetUserMeta())
	app.GET("/websocket", handle.WebSocketHandle)
}
