package service

import (
	"baby-fried-rice/internal/pkg/kit/middleware"
	"baby-fried-rice/internal/pkg/module/space/service/handle"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	root := engine.Group("/api/space", middleware.SetUserMeta())
	root.GET("", handle.QuerySpacesHandle)
	root.POST("", handle.AddSpaceHandle)
	root.DELETE("", handle.DeleteSpaceHandle)
}
