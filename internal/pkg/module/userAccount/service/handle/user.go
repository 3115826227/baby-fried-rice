package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/grpc/pbservices/user"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/module/userAccount/cache"
	"baby-fried-rice/internal/pkg/module/userAccount/config"
	"baby-fried-rice/internal/pkg/module/userAccount/grpc"
	"baby-fried-rice/internal/pkg/module/userAccount/log"
	"baby-fried-rice/internal/pkg/module/userAccount/service/model"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
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
	req.Ip = c.GetHeader(handle.HeaderIP)

	client, err := grpc.GetClientGRPC(config.GetConfig().Servers.AccountDaoServer)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var reqLogin = &user.ReqPasswordLogin{
		LoginName: req.LoginName,
		Password:  req.Password,
		Ip:        req.Ip,
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
		cache.GetCache().Add(fmt.Sprintf("%v:%v", constant.TokenPrefix, token), userMeta.ToString())
		cache.GetCache().Add(userMeta.UserId, fmt.Sprintf("%v:%v", constant.TokenPrefix, token))
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

	client, err := grpc.GetClientGRPC(config.GetConfig().Servers.AccountDaoServer)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var reqRegister = &user.ReqUserRegister{
		Login: &user.ReqPasswordLogin{
			LoginName: req.LoginName,
			Password:  req.Password,
			Ip:        c.GetHeader(handle.HeaderIP),
		},
		Username: req.Username,
		Gender:   req.Gender,
		Phone:    req.Phone,
	}
	resp, err := user.NewDaoUserClient(client.GetRpcClient()).
		UserDaoRegister(context.Background(), reqRegister)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	if resp.Code != handle.SuccessCode {
		log.Logger.Error(resp.Message)
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
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
