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

	api.POST("/admin/login", HandleManageProxy)
	api.POST("/user/register", HandleAccountUserProxy)
	api.POST("/user/login", HandleAccountUserProxy)

	user := api.Group("")
	user.Use(kitMiddleware.GenerateUUID)
	user.Use(middleware.CheckToken)

	user.Any("/manage/*any", HandleManageProxy)
	user.Any("/account/user/*any", HandleAccountUserProxy)
	user.Any("/im/*any", HandleImProxy)
	user.Any("/space/*any", HandleSpaceProxy)
	user.Any("/connect/*any", HandleConnectProxy)
	user.Any("/file/*any", HandleFileProxy)
	user.Any("/shop/*any", HandleShopProxy)
}

func HandleAccountUserProxy(c *gin.Context) {
	handleProxy(c, config.GetConfig().Servers.UserAccountServer)
}

func HandleManageProxy(c *gin.Context) {
	handleProxy(c, config.GetConfig().Servers.Manage)
}

func HandleImProxy(c *gin.Context) {
	handleProxy(c, config.GetConfig().Servers.ImServer)
}

func HandleSpaceProxy(c *gin.Context) {
	handleProxy(c, config.GetConfig().Servers.SpaceServer)
}

func HandleConnectProxy(c *gin.Context) {
	handleProxy(c, config.GetConfig().Servers.ConnectServer)
}

func HandleFileProxy(c *gin.Context) {
	handleProxy(c, config.GetConfig().Servers.FileServer)
}

func HandleShopProxy(c *gin.Context) {
	handleProxy(c, config.GetConfig().Servers.ShopServer)
}

func handleProxy(c *gin.Context, serverName string) {
	serverUrl, err := server.GetRegisterClient().GetServer(serverName)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var parseServerUrl *url.URL
	if parseServerUrl, err = url.Parse(serverUrl); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(parseServerUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}
