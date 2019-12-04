package middlware

import (
	"github.com/gin-gonic/gin"
	"github.com/3115826227/baby-fried-rice/module/account/service/login/handle"
	"net/http"
	"github.com/3115826227/baby-fried-rice/module/account/service/model"
)

func Middleware(c *gin.Context) {

}

func MiddlewareSetUserMeta() gin.HandlerFunc {
	return func(context *gin.Context) {
		header := context.Request.Header
		userId := header.Get(handle.HeaderUserId)
		clientId := header.Get(handle.HeaderClientId)
		platform := header.Get(handle.HeaderPlatform)
		reqId := header.Get(handle.HeaderReqId)
		isSuper := header.Get(handle.HeaderIsSuper)

		if isSuper == "" || userId == "" || clientId == "" || platform == "" || reqId == "" {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求头错误",
			})
		}

		userMeta := model.UserMeta{
			UserId:   userId,
			ClientId: clientId,
			ReqId:    reqId,
			Platform: platform,
			IsSuper:  isSuper,
		}

		context.Set(handle.GinContextKeyUserMeta, &userMeta)
	}
}
