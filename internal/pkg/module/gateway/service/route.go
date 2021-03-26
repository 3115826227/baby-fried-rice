package service

import (
	kitMiddleware "baby-fried-rice/internal/pkg/kit/middleware"
	"baby-fried-rice/internal/pkg/module/gateway/config"
	"baby-fried-rice/internal/pkg/module/gateway/middleware"
	"github.com/gin-gonic/gin"
	"net/http/httputil"
)

func Register(engine *gin.Engine) {
	api := engine.Group("/api")

	api.POST("/root/login", HandleAccountRootProxy)
	api.POST("/admin/login", HandleAccountAdminProxy)
	api.POST("/user/register", HandleAccountUserProxy)
	api.POST("/user/login", HandleAccountUserProxy)
	api.GET("/user/school", HandleAccountUserProxy)
	api.GET("/user/admin/permission", HandleAccountUserProxy)
	api.GET("/user/friend/chat", HandleImProxy)

	user := api.Group("")
	user.Use(kitMiddleware.GenerateUUID)
	user.Use(middleware.CheckToken)

	user.Any("/account/user/*any", HandleAccountUserProxy)
	user.Any("/account/admin/*any", HandleAccountAdminProxy)
	user.Any("/account/root/*any", HandleAccountRootProxy)
	user.Any("/public/*any", HandlePublicProxy)
	user.Any("/im/*any", HandleImProxy)
	user.Any("/square/*any", HandleSquareProxy)
}

func HandleAccountUserProxy(c *gin.Context) {
	proxy := httputil.NewSingleHostReverseProxy(config.GetConfig().ParserUserUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func HandleAccountAdminProxy(c *gin.Context) {
	proxy := httputil.NewSingleHostReverseProxy(config.GetConfig().ParserAdminUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func HandleAccountRootProxy(c *gin.Context) {
	proxy := httputil.NewSingleHostReverseProxy(config.GetConfig().ParserRootUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func HandlePublicProxy(c *gin.Context) {
	proxy := httputil.NewSingleHostReverseProxy(config.GetConfig().ParserPublicUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func HandleImProxy(c *gin.Context) {
	proxy := httputil.NewSingleHostReverseProxy(config.GetConfig().ParserImUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func HandleSquareProxy(c *gin.Context) {
	proxy := httputil.NewSingleHostReverseProxy(config.GetConfig().ParserSquareUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}
