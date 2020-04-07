package service

import (
	"github.com/3115826227/baby-fried-rice/module/gateway/src/config"
	"github.com/3115826227/baby-fried-rice/module/gateway/src/middleware"
	"github.com/gin-gonic/gin"
	"net/http/httputil"
)

func RegisterRouter(engine *gin.Engine) {
	api := engine.Group("/api")

	api.POST("/root/login", HandleAccountProxy)
	api.POST("/admin/login", HandleAccountProxy)
	api.POST("/user/register", HandleAccountProxy)
	api.POST("/user/login", HandleAccountProxy)
	api.GET("/user/school", HandleAccountProxy)
	api.GET("/user/admin/permission", HandleAccountProxy)
	api.GET("/user/friend/chat", HandleImProxy)

	user := api.Group("")
	user.Use(middleware.GenerateUUID)
	user.Use(middleware.CheckToken)

	user.Any("/account/*any", HandleAccountProxy)
	user.Any("/public/*any", HandlePublicProxy)
	user.Any("/im/*any", HandleImProxy)
	user.Any("/square/", HandleSquareProxy)
}

func HandleAccountProxy(c *gin.Context) {
	proxy := httputil.NewSingleHostReverseProxy(config.Config.ParserAccountUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func HandlePublicProxy(c *gin.Context) {
	proxy := httputil.NewSingleHostReverseProxy(config.Config.ParserPublicUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func HandleImProxy(c *gin.Context) {
	proxy := httputil.NewSingleHostReverseProxy(config.Config.ParserImUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func HandleSquareProxy(c *gin.Context) {
	proxy := httputil.NewSingleHostReverseProxy(config.Config.ParserSquareUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}
