package middleware

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetUserMeta() gin.HandlerFunc {
	return func(context *gin.Context) {
		header := context.Request.Header
		accountId := header.Get(handle.HeaderAccountId)
		username := header.Get(handle.HeaderUsername)
		schoolId := header.Get(handle.HeaderSchoolId)
		platform := header.Get(handle.HeaderPlatform)
		reqId := header.Get(handle.HeaderReqId)
		isSuper := header.Get(handle.HeaderIsSuper)

		if accountId == "" {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求头错误",
			})
		}

		userMeta := handle.UserMeta{
			AccountId: accountId,
			Username:  username,
			SchoolId:  schoolId,
			ReqId:     reqId,
			Platform:  platform,
			IsSuper:   isSuper,
		}

		context.Set(handle.GinContextKeyUserMeta, &userMeta)
	}
}
