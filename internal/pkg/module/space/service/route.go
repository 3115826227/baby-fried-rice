package service

import (
	"baby-fried-rice/internal/pkg/kit/middleware"
	"baby-fried-rice/internal/pkg/module/space/service/handle"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	handle.Init()
	app := engine.Group("/api/space", middleware.SetUserMeta())
	app.GET("/space", handle.SpacesQueryHandle)
	app.POST("/space", handle.SpaceAddHandle)
	app.DELETE("/space", handle.SpaceDeleteHandle)

	app.POST("/operator", handle.SpaceOptAddHandle)
	app.DELETE("/operator", handle.SpaceOptCancelHandle)
	app.POST("/comment", handle.SpaceCommentAddHandle)
	app.DELETE("/comment", handle.SpaceCommentDeleteHandle)
}
