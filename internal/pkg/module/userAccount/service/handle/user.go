package handle

import (
	"baby-fried-rice/internal/pkg/kit/grpc/pbservices/user"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/module/userAccount/cache"
	"baby-fried-rice/internal/pkg/module/userAccount/config"
	"baby-fried-rice/internal/pkg/module/userAccount/grpc"
	"baby-fried-rice/internal/pkg/module/userAccount/log"
	"baby-fried-rice/internal/pkg/module/userAccount/server"
	"baby-fried-rice/internal/pkg/module/userAccount/service/model"
	"context"
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

func UserLoginHandle(c *gin.Context) {
	var err error
	var req requests.PasswordLoginReq
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	req.LoginName = strings.TrimSpace(req.LoginName)
	req.Password = strings.TrimSpace(req.Password)
	req.Ip = c.GetHeader("IP")

	client, err := grpc.GetClientGRPC(config.GetConfig().Servers.AccountDaoServer)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var reqLogin = &user.ReqPasswordLogin{
		LoginName: req.LoginName,
		Password:  req.Password,
	}
	resp, err := user.NewDaoUserClient(client.GetRpcClient()).
		UserDaoLogin(context.Background(), reqLogin)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}

	token, err := handle.GenerateToken(resp.User.AccountId, time.Now(), config.GetConfig().TokenSecret)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}

	var result = model.RspLogin{
		RspSuccess: handle.RspSuccess{Code: handle.SuccessCode},
		Data: model.LoginResult{
			UserInfo: model.RspUserData{
				UserId:    resp.User.AccountId,
				LoginName: resp.User.LoginName,
				Username:  resp.User.Username,
			},
			Token: token,
		},
	}

	go func() {
		var userMeta = &handle.UserMeta{
			UserId:   resp.User.AccountId,
			Platform: "pc",
		}
		cache.GetCache().Add(fmt.Sprintf("%v:%v", TokenPrefix, token), userMeta.ToString())
		cache.GetCache().Add(userMeta.UserId, fmt.Sprintf("%v:%v", TokenPrefix, token))
	}()

	c.JSON(http.StatusOK, result)
}

func UserRegisterHandle(c *gin.Context) {
	var err error
	var req requests.UserRegisterReq
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	req.LoginName = strings.TrimSpace(req.LoginName)
	req.Password = strings.TrimSpace(req.Password)
	req.Phone = strings.TrimSpace(req.Phone)

	payload, err := json.Marshal(req)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.SysErrResponse)
		return
	}

	accountDaoUrl, err := server.GetRegisterClient().GetServer(config.GetConfig().Servers.AccountDaoServer)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.SysErrResponse)
		return
	}
	data, err := handle.Post(accountDaoUrl+"/dao/account/user/register", payload, c.Request.Header.Clone())
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.SysErrResponse)
		return
	}
	ok, err := handle.ResponseHandle(data)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.SysErrResponse)
		return
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

func UserLogout(c *gin.Context) {
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
