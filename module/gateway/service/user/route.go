package user

import (
	"github.com/3115826227/baby-fried-rice/module/gateway/config"
	"github.com/gin-gonic/gin"
	"net/http/httputil"
)

/*
	用户登录鉴权部分
*/
func Route(apiGroup *gin.RouterGroup) {
	userGroup := apiGroup.Group("/user")

	userGroup.POST("/root/login", HandleUserProxy)
	userGroup.GET("/root/logout", HandleUserProxy)
	userGroup.POST("/login", HandleUserProxy)
	userGroup.GET("/logout", HandleUserProxy)
}

func HandleUserProxy(c *gin.Context) {
	proxy := httputil.NewSingleHostReverseProxy(config.Config.ParseUserUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}
