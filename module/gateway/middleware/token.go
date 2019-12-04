package middleware

import (
	"github.com/3115826227/baby-fried-rice/module/gateway/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckToken(c *gin.Context) {
	token := c.GetHeader(HeaderToken)

	if token == "" {
		util.ErrorResp(c, http.StatusUnauthorized, util.CodeRequiredLogin, util.CodeRequiredLoginMsg)
		c.Abort()
		return
	}
}
