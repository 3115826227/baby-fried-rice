package service

import (
	"baby-fried-rice/internal/pkg/module/file/service/handle"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	group := engine.Group("/api")
	group.POST("/upload", handle.FileUploadHandle)
}
