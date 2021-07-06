package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/module/backend/cache"
	"baby-fried-rice/internal/pkg/module/backend/config"
	"baby-fried-rice/internal/pkg/module/backend/log"
	"baby-fried-rice/internal/pkg/module/backend/query"
	"baby-fried-rice/internal/pkg/module/backend/service/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// 超级用户账号登录
func RootLogin(c *gin.Context) {
	var req model.ReqLogin
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	req.LoginName = strings.TrimSpace(req.LoginName)
	req.Password = strings.TrimSpace(req.Password)
	req.Password = handle.EncodePassword(req.Password)
	req.Ip = c.GetHeader("IP")
	root, err := query.GetRootByLogin(req.LoginName, req.Password)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}

	var token string
	token, err = handle.GenerateToken(root.ID, time.Now(), config.GetConfig().TokenSecret)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}

	var result = model.RspLogin{
		RspSuccess: handle.RspSuccess{Code: handle.SuccessCode},
		Data: model.LoginResult{
			UserInfo: model.RspUserData{
				UserId:    root.ID,
				LoginName: root.LoginName,
				Username:  root.Username,
			},
			Token: token,
		},
	}

	go func() {
		var userMeta = &handle.UserMeta{
			AccountId: root.ID,
			Platform:  "pc",
		}
		cache.GetCache().Add(fmt.Sprintf("%v:%v", constant.TokenPrefix, token), userMeta.ToString())
		cache.GetCache().Add(userMeta.AccountId, fmt.Sprintf("%v:%v", constant.TokenPrefix, token))
	}()

	c.JSON(http.StatusOK, result)
}

func RootLogout(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	token, err := cache.GetCache().Get(userMeta.AccountId)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	go func() {
		cache.GetCache().Del(token)
		cache.GetCache().Del(userMeta.AccountId)
	}()

	handle.SuccessResp(c, "", handle.RspOkResponse{})
}
