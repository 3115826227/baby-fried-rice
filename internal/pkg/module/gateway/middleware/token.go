package middleware

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/module/gateway/cache"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckToken(c *gin.Context) {
	var token = c.GetHeader(handle.HeaderToken)

	if token == "" {
		token = c.Query(handle.HeaderToken)
	}
	if token == "" {
		handle.ErrorResp(c, http.StatusUnauthorized, constant.CodeRequiredLogin, constant.CodeRequiredLoginMsg)
		c.Abort()
		return
	}

	var tokenKey = fmt.Sprintf("%v:%v", handle.TokenPrefix, token)
	var str, err = cache.GetCache().Get(tokenKey)
	if err != nil {
		handle.ErrorResp(c, http.StatusUnauthorized, constant.CodeRequiredLogin, constant.CodeRequiredLoginMsg)
		c.Abort()
		return
	}
	var userMeta handle.UserMeta
	err = json.Unmarshal([]byte(str), &userMeta)
	if err != nil {
		handle.ErrorResp(c, http.StatusUnauthorized, constant.CodeRequiredLogin, constant.CodeRequiredLoginMsg)
		c.Abort()
		return
	}
	var header = c.Request.Header
	header.Set(handle.HeaderAccountId, userMeta.AccountId)
	header.Set(handle.HeaderUsername, userMeta.Username)
	header.Set(handle.HeaderSchoolId, userMeta.SchoolId)
	header.Set(handle.HeaderReqId, userMeta.ReqId)
	header.Set(handle.HeaderIsOfficial, fmt.Sprintf("%v", userMeta.IsOfficial))
	header.Set(handle.HeaderPlatform, userMeta.Platform)
	header.Set(handle.HeaderPhone, userMeta.Phone)
	c.Next()
}
