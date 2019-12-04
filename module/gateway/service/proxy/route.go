package proxy

import (
	"github.com/3115826227/baby-fried-rice/module/gateway/middleware"
	"github.com/gin-gonic/gin"
)

func Route(apiGroup *gin.RouterGroup) {
	proxyGroup := apiGroup.Group("/api/%s")

	proxyGroup.Use(middleware.Cors())
	proxyGroup.Use(middleware.CheckToken)
}
