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
	app.POST("/comment", handle.SpaceCommentAddHandle)
	app.GET("/comment", handle.CommentQueryHandle)
	app.GET("/comment/reply", handle.CommentReplyQueryHandle)
	app.DELETE("/comment", handle.SpaceCommentDeleteHandle)
}
