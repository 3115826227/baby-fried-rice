package service

import (
	"baby-fried-rice/internal/pkg/kit/middleware"
	"baby-fried-rice/internal/pkg/module/file/service/handle"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	handle.InitBackend()

	group := engine.Group("/api/file/", middleware.SetUserMeta())
	group.POST("/upload", handle.FileUploadHandle)

	group.GET("/file", handle.FileQueryHandle)
	group.DELETE("/file", handle.FileDeleteHandle)
}
