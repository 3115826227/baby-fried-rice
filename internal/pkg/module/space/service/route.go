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

	common := engine.Group("/api/comment", middleware.SetUserMeta())
	common.POST("/operator", handle.OptAddHandle)
	common.POST("/comment", handle.CommentAddHandle)
	common.GET("/comment", handle.CommentQueryHandle)
	common.GET("/comment/reply", handle.CommentReplyQueryHandle)
	common.DELETE("/comment", handle.CommentDeleteHandle)
}
