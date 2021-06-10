package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
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
			AccountId: resp.User.AccountId,
			Platform:  "pc",
		}
		cache.GetCache().Add(fmt.Sprintf("%v:%v", constant.TokenPrefix, token), userMeta.ToString())
		cache.GetCache().Add(userMeta.AccountId, fmt.Sprintf("%v:%v", constant.TokenPrefix, token))
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
	_, err = user.NewDaoUserClient(client.GetRpcClient()).
		UserDaoRegister(context.Background(), reqRegister)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

func UserLogoutHandle(c *gin.Context) {
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

	handle.SuccessResp(c, "", nil)
}

// 查看用户自己信息
func UserDetailHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	client, err := grpc.GetClientGRPC(config.GetConfig().Servers.AccountDaoServer)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	resp, err := user.NewDaoUserClient(client.GetRpcClient()).
		UserDaoDetail(context.Background(), &user.ReqDaoUserDetail{AccountId: userMeta.AccountId})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var detailRsp = model.RspUserDetail{
		AccountId:  resp.Detail.AccountId,
		Describe:   resp.Detail.Describe,
		HeadImgUrl: resp.Detail.HeadImgUrl,
		Username:   resp.Detail.Username,
		SchoolId:   resp.Detail.SchoolId,
		Gender:     resp.Detail.Gender,
		Age:        resp.Detail.Age,
		Phone:      resp.Detail.Phone,
	}
	handle.SuccessResp(c, "", detailRsp)
}

func UserDetailUpdateHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var err error
	var req requests.UserDetailUpdateReq
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var updateReq = &user.ReqDaoUserDetailUpdate{
		Detail: &user.DaoUserDetail{
			AccountId:  userMeta.AccountId,
			Describe:   req.Describe,
			HeadImgUrl: req.HeadImgUrl,
			Username:   req.Username,
			Gender:     req.Gender,
			Age:        req.Age,
			Phone:      req.Phone,
		},
	}
	client, err := grpc.GetClientGRPC(config.GetConfig().Servers.AccountDaoServer)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	_, err = user.NewDaoUserClient(client.GetRpcClient()).
		UserDaoDetailUpdate(context.Background(), updateReq)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

func UserPwdUpdateHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var err error
	var req requests.UserPwdUpdateReq
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	client, err := grpc.GetClientGRPC(config.GetConfig().Servers.AccountDaoServer)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var updateReq = &user.ReqDaoUserPwdUpdate{
		AccountId:   userMeta.AccountId,
		Password:    req.Password,
		NewPassword: req.NewPassword,
	}
	_, err = user.NewDaoUserClient(client.GetRpcClient()).
		UserDaoPwdUpdate(context.Background(), updateReq)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 查看他人用户信息
func UserQueryHandle(c *gin.Context) {
	accountId := c.Query("account_id")
	client, err := grpc.GetClientGRPC(config.GetConfig().Servers.AccountDaoServer)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	resp, err := user.NewDaoUserClient(client.GetRpcClient()).
		UserDaoDetail(context.Background(), &user.ReqDaoUserDetail{AccountId: accountId})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var detailRsp = model.RspUserDetail{
		AccountId:  resp.Detail.AccountId,
		Describe:   resp.Detail.Describe,
		HeadImgUrl: resp.Detail.HeadImgUrl,
		Username:   resp.Detail.Username,
		SchoolId:   resp.Detail.SchoolId,
		Gender:     resp.Detail.Gender,
		Age:        resp.Detail.Age,
	}
	handle.SuccessResp(c, "", detailRsp)
}
