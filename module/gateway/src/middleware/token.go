package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/gateway/src/redis"
	"github.com/3115826227/baby-fried-rice/module/gateway/src/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckToken(c *gin.Context) {
	token := c.GetHeader(HeaderToken)

	if token == "" {
		token = c.Query(HeaderToken)
		if token == "" {
			util.ErrorResp(c, http.StatusUnauthorized, util.CodeRequiredLogin, util.CodeRequiredLoginMsg)
			c.Abort()
			return
		}
	}

	tokenKey := fmt.Sprintf("%v:%v", TokenPrefix, token)
	str, err := redis.Get(tokenKey)
	if err != nil {
		util.ErrorResp(c, http.StatusUnauthorized, util.CodeRequiredLogin, util.CodeRequiredLoginMsg)
		c.Abort()
		return
	}
	var userMeta UserMeta
	err = json.Unmarshal([]byte(str), &userMeta)
	if err != nil {
		util.ErrorResp(c, http.StatusUnauthorized, util.CodeRequiredLogin, util.CodeRequiredLoginMsg)
		c.Abort()
		return
	}
	header := c.Request.Header
	header.Set(HeaderUserId, userMeta.UserId)
	header.Set(HeaderSchoolId, userMeta.SchoolId)
	header.Set(HeaderReqId, userMeta.ReqId)
	header.Set(HeaderIsSuper, userMeta.IsSuper)
	header.Set(HeaderPlatform, userMeta.Platform)
	c.Next()
}
