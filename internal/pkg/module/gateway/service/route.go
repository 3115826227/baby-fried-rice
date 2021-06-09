package service

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	kitMiddleware "baby-fried-rice/internal/pkg/kit/middleware"
	"baby-fried-rice/internal/pkg/module/gateway/config"
	"baby-fried-rice/internal/pkg/module/gateway/log"
	"baby-fried-rice/internal/pkg/module/gateway/middleware"
	"baby-fried-rice/internal/pkg/module/gateway/server"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
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
	userUrl, err := server.GetRegisterClient().GetServer(config.GetConfig().Servers.UserAccountServer)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	parserUserUrl, err := url.Parse(userUrl)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(parserUserUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func HandleAccountAdminProxy(c *gin.Context) {
	proxy := httputil.NewSingleHostReverseProxy(config.GetConfig().ParserAdminUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func HandleAccountRootProxy(c *gin.Context) {
	adminUrl, err := server.GetRegisterClient().GetServer(config.GetConfig().Servers.RootAccountServer)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	parserAdminUrl, err := url.Parse(adminUrl)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(parserAdminUrl)
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
