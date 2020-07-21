package middleware

import (
	"github.com/3115826227/baby-fried-rice/module/user-account/src/service/handle"
	"github.com/3115826227/baby-fried-rice/module/user-account/src/service/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Middleware(c *gin.Context) {

}

func MiddlewareSetUserMeta() gin.HandlerFunc {
	return func(context *gin.Context) {
		header := context.Request.Header
		userId := header.Get(handle.HeaderUserId)
		username := header.Get(handle.HeaderUsername)
		schoolId := header.Get(handle.HeaderSchoolId)
		platform := header.Get(handle.HeaderPlatform)
		reqId := header.Get(handle.HeaderReqId)
		isSuper := header.Get(handle.HeaderIsSuper)

		if isSuper == "" || userId == "" {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求头错误",
			})
		}

		userMeta := model.UserMeta{
			UserId:   userId,
			Username: username,
			SchoolId: schoolId,
			ReqId:    reqId,
			Platform: platform,
			IsSuper:  isSuper,
		}

		context.Set(handle.GinContextKeyUserMeta, &userMeta)
	}
}
