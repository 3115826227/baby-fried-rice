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

	api.POST("/admin/login", HandleBackendProxy)
	api.POST("/user/register", HandleAccountUserProxy)
	api.POST("/user/login", HandleAccountUserProxy)

	user := api.Group("")
	user.Use(kitMiddleware.GenerateUUID)
	user.Use(middleware.CheckToken)

	user.Any("/backend/*any", HandleBackendProxy)
	user.Any("/account/user/*any", HandleAccountUserProxy)
	user.Any("/im/*any", HandleImProxy)
	user.Any("/space/*any", HandleSpaceProxy)
	user.Any("/connect/*any", HandleConnectProxy)
	user.Any("/file/*any", HandleFileProxy)
}

func HandleAccountUserProxy(c *gin.Context) {
	userUrl, err := server.GetRegisterClient().GetServer(config.GetConfig().Servers.UserAccountServer)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var parserUserUrl *url.URL
	if parserUserUrl, err = url.Parse(userUrl); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(parserUserUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func HandleBackendProxy(c *gin.Context) {
	backendUrl, err := server.GetRegisterClient().GetServer(config.GetConfig().Servers.Backend)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var parseBackendUrl *url.URL
	if parseBackendUrl, err = url.Parse(backendUrl); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(parseBackendUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func HandleImProxy(c *gin.Context) {
	imUrl, err := server.GetRegisterClient().GetServer(config.GetConfig().Servers.ImServer)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var parserImUrl *url.URL
	if parserImUrl, err = url.Parse(imUrl); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(parserImUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func HandleSpaceProxy(c *gin.Context) {
	spaceUrl, err := server.GetRegisterClient().GetServer(config.GetConfig().Servers.SpaceServer)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var parserSpaceUrl *url.URL
	if parserSpaceUrl, err = url.Parse(spaceUrl); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(parserSpaceUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func HandleConnectProxy(c *gin.Context) {
	connectUrl, err := server.GetRegisterClient().GetServer(config.GetConfig().Servers.ConnectServer)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var parserConnectUrl *url.URL
	if parserConnectUrl, err = url.Parse(connectUrl); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(parserConnectUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func HandleFileProxy(c *gin.Context) {
	fileUrl, err := server.GetRegisterClient().GetServer(config.GetConfig().Servers.FileServer)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var parserFileUrl *url.URL
	if parserFileUrl, err = url.Parse(fileUrl); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(parserFileUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}
