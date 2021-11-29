package service

import (
	"baby-fried-rice/internal/pkg/kit/middleware"
	"baby-fried-rice/internal/pkg/module/live/service/handle"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	app := engine.Group("/api/live", middleware.SetUserMeta())
	app.GET("/room", handle.LiveRoomHandle)
	app.GET("/room/detail", handle.LiveRoomDetailHandle)
	app.GET("/room/user", handle.LiveRoomUserHandle)
	app.POST("/room/origin", handle.LiveRoomOriginUpdateHandle)
	app.POST("/room/user/opt", handle.LiveRoomUserOptUpdateHandle)
	app.GET("/room/message", handle.LiveRoomUserMessageHandle)
}
