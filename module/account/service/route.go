package service

import (
	"github.com/3115826227/baby-fried-rice/module/account/service/login/handle"
	"github.com/gin-gonic/gin"
	"github.com/3115826227/baby-fried-rice/module/account/middlware"
)

func RegisterRoute(engine *gin.Engine) {

	engine.POST("/api/root/login", handle.RootLogin)

	app := engine.Group("/api")

	app.Use(middlware.MiddlewareSetUserMeta())

	app.POST("/client", handle.ClientAdd)
}
