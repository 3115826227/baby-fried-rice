package service

import (
	"github.com/3115826227/baby-fried-rice/module/account/src/middlware"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/login/handle"
	"github.com/gin-gonic/gin"
)

func RegisterRoute(engine *gin.Engine) {

	engine.POST("/api/root/login", handle.RootLogin)

	engine.POST("/api/user/register", handle.UserRegister)
	engine.POST("/api/user/login", handle.UserLogin)

	app := engine.Group("/api/account")

	app.Use(middlware.MiddlewareSetUserMeta())

	app.GET("/user", handle.UserDetail)
	app.PATCH("/user", handle.UserUpdate)

	app.POST("/client", handle.ClientAdd)
}
