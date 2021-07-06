package service

import (
	"baby-fried-rice/internal/pkg/kit/middleware"
	"baby-fried-rice/internal/pkg/module/backend/service/handle"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	engine.POST("/api/admin/login", handle.RootLogin)
	root := engine.Group("/api/backend", middleware.SetUserMeta())

	root.GET("/admin/logout", handle.RootLogout)

}
