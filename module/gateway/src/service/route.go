package service

import (
	"github.com/gin-gonic/gin"
	"github.com/3115826227/baby-fried-rice/module/gateway/src/middleware"
	"github.com/3115826227/baby-fried-rice/module/gateway/src/config"
	"net/http/httputil"
)

func RegisterRouter(engine *gin.Engine) {
	api := engine.Group("/api")

	api.POST("/root/login", HandleAccountProxy)
	api.POST("/user/register", HandleAccountProxy)
	api.POST("/user/login", HandleAccountProxy)

	user := api.Group("")
	user.Use(middleware.GenerateUUID)
	user.Use(middleware.CheckToken)

	user.Any("/account/*any", HandleAccountProxy)
	user.Any("/public/", HandlePublicProxy)
	user.Any("/im/", HandleImProxy)
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