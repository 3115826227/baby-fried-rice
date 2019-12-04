package service

import (
	"github.com/3115826227/baby-fried-rice/module/gateway/middleware"
	"github.com/3115826227/baby-fried-rice/module/gateway/service/common"
	"github.com/3115826227/baby-fried-rice/module/gateway/service/proxy"
	"github.com/3115826227/baby-fried-rice/module/gateway/service/register"
	"github.com/3115826227/baby-fried-rice/module/gateway/service/user"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(engine *gin.Engine) {
	api := engine.Group("/api")

	api.Use(middleware.GenerateUUID)

	common.Route(api)
	register.Route(api)
	user.Route(api)
	proxy.Route(api)
}
