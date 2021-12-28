package middleware

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/module/gateway/log"
	"baby-fried-rice/internal/pkg/module/gateway/server"
	"fmt"
	"github.com/gin-gonic/gin"
)

// 用户访问鉴权
func Auth(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	// 手机未认证用户部分功能将不可用
	if userMeta.Phone != "" {
		var path = fmt.Sprintf("%v:%v", c.Request.Method, c.Request.URL.Path)
		exist, err := server.GetRegisterClient().IsAuthPathConfig(path)
		if err != nil {
			log.Logger.Error(err.Error())
			handle.SystemErrorResponse(c)
			c.Abort()
			return
		}
		if exist {
			handle.FailedResp(c, handle.CodeUnVerifyForbidden)
			c.Abort()
			return
		}
	}
	c.Next()
}
