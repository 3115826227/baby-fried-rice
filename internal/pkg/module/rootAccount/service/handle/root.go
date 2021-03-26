package handle

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/module/rootAccount/cache"
	"baby-fried-rice/internal/pkg/module/rootAccount/config"
	"baby-fried-rice/internal/pkg/module/rootAccount/log"
	"baby-fried-rice/internal/pkg/module/rootAccount/service/model"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

const (
	TokenPrefix = "token"
)

// 超级用户账号登录
func RootLogin(c *gin.Context) {
	var err error
	var req model.ReqLogin
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	req.LoginName = strings.TrimSpace(req.LoginName)
	req.Password = strings.TrimSpace(req.Password)
	req.Password = handle.EncodePassword(req.Password)
	req.Ip = c.GetHeader("IP")

	payload, err := json.Marshal(req)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.SysErrResponse)
		return
	}

	data, err := handle.Post(config.GetConfig().Connect.AccountDaoUrl+"/dao/account/root/login", payload, c.Request.Header.Clone())
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.SysErrResponse)
		return
	}
	var resp model.RspDaoRootLogin
	err = json.Unmarshal(data, &resp)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.SysErrResponse)
		return
	}
	if resp.Code != handle.SuccessCode {
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.SysErrResponse)
		return
	}

	token, err := handle.GenerateToken(resp.Data.ID, time.Now(), config.GetConfig().TokenSecret)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}

	var result = model.RspLogin{
		RspSuccess: handle.RspSuccess{Code: handle.SuccessCode},
		Data: model.LoginResult{
			UserInfo: model.RspUserData{
				UserId:    resp.Data.ID,
				LoginName: resp.Data.LoginName,
				Username:  resp.Data.Username,
			},
			Token: token,
		},
	}

	go func() {
		var userMeta = &handle.UserMeta{
			UserId:   resp.Data.ID,
			Platform: "pc",
		}
		cache.GetCache().Add(fmt.Sprintf("%v:%v", TokenPrefix, token), userMeta.ToString())
		cache.GetCache().Add(userMeta.UserId, fmt.Sprintf("%v:%v", TokenPrefix, token))
	}()

	c.JSON(http.StatusOK, result)
}

func RootLogout(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	token, err := cache.GetCache().Get(userMeta.UserId)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	go func() {
		cache.GetCache().Del(token)
		cache.GetCache().Del(userMeta.UserId)
	}()

	handle.SuccessResp(c, "", handle.RspOkResponse{})
}
