package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/userAccount/cache"
	"baby-fried-rice/internal/pkg/module/userAccount/config"
	"baby-fried-rice/internal/pkg/module/userAccount/grpc"
	"baby-fried-rice/internal/pkg/module/userAccount/log"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// 用户登录
func UserLoginHandle(c *gin.Context) {
	var req requests.PasswordLoginReq
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	req.LoginName = strings.TrimSpace(req.LoginName)
	req.Password = strings.TrimSpace(req.Password)
	req.Ip = c.GetHeader(handle.HeaderIP)
	var reqLogin = &user.ReqPasswordLogin{
		LoginName: req.LoginName,
		Password:  req.Password,
		Ip:        req.Ip,
	}
	userClient, err := grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.LoginErrResponse)
		return
	}
	var resp *user.RspDaoUserLogin
	resp, err = userClient.UserDaoLogin(context.Background(), reqLogin)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.LoginErrResponse)
		return
	}

	var token string
	token, err = handle.GenerateToken(resp.User.AccountId, time.Now(), config.GetConfig().TokenSecret)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}

	var loginResult = rsp.LoginResult{
		UserInfo: rsp.UserDataResp{
			UserId:    resp.User.AccountId,
			LoginName: resp.User.LoginName,
			Username:  resp.User.Username,
		},
		Token: token,
	}

	go func() {
		var userMeta = &handle.UserMeta{
			AccountId: resp.User.AccountId,
			Username:  resp.User.Username,
			Platform:  "pc",
		}
		if err = cache.GetCache().Add(fmt.Sprintf("%v:%v", constant.TokenPrefix, token), userMeta.ToString()); err != nil {
			log.Logger.Error(err.Error())
		}
		if err = cache.GetCache().Add(userMeta.AccountId, fmt.Sprintf("%v:%v", constant.TokenPrefix, token)); err != nil {
			log.Logger.Error(err.Error())
		}
	}()

	handle.SuccessResp(c, "", loginResult)
}

// 用户注册
func UserRegisterHandle(c *gin.Context) {
	var req requests.UserRegisterReq
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	req.LoginName = strings.TrimSpace(req.LoginName)
	req.Password = strings.TrimSpace(req.Password)
	req.Phone = strings.TrimSpace(req.Phone)

	userClient, err := grpc.GetUserClient()
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
		Phone:    req.Phone,
	}
	_, err = userClient.UserDaoRegister(context.Background(), reqRegister)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 用户退出登录
func UserLogoutHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	token, err := cache.GetCache().Get(userMeta.AccountId)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	go func() {
		if err = cache.GetCache().Del(token); err != nil {
			log.Logger.Error(err.Error())
		}
		if err = cache.GetCache().Del(userMeta.AccountId); err != nil {
			log.Logger.Error(err.Error())
		}
	}()

	handle.SuccessResp(c, "", nil)
}

// 查看用户自己信息
func UserDetailHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	userClient, err := grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var resp *user.RspDaoUserDetail
	resp, err = userClient.UserDaoDetail(context.Background(),
		&user.ReqDaoUserDetail{AccountId: userMeta.AccountId})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var detailRsp = rsp.UserDetailResp{
		AccountId:  resp.Detail.AccountId,
		Describe:   resp.Detail.Describe,
		HeadImgUrl: resp.Detail.HeadImgUrl,
		Username:   resp.Detail.Username,
		SchoolId:   resp.Detail.SchoolId,
		Gender:     resp.Detail.Gender,
		Age:        resp.Detail.Age,
		Phone:      resp.Detail.Phone,
		Coin:       resp.Detail.Coin,
		IsOfficial: resp.Detail.IsOfficial,
	}
	handle.SuccessResp(c, "", detailRsp)
}

// 用户更新自己信息
func UserDetailUpdateHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var req requests.UserDetailUpdateReq
	if err := c.ShouldBind(&req); err != nil {
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
	userClient, err := grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	_, err = userClient.UserDaoDetailUpdate(context.Background(), updateReq)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 用户更新密码
func UserPwdUpdateHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var req requests.UserPwdUpdateReq
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	userClient, err := grpc.GetUserClient()
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
	_, err = userClient.UserDaoPwdUpdate(context.Background(), updateReq)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 查看他人用户信息
func UserQueryHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	accountId := c.Query("account_id")
	if accountId == "" {
		err := fmt.Errorf("account_id can't null")
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	userClient, err := grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var resp *user.RspDaoUserDetail
	resp, err = userClient.UserDaoDetail(context.Background(),
		&user.ReqDaoUserDetail{AccountId: accountId})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var imClient im.DaoImClient
	imClient, err = grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var imReq = im.ReqIsFriendDao{
		Origin:    userMeta.AccountId,
		AccountId: accountId,
	}
	var imResp *im.RspIsFriendDao
	imResp, err = imClient.FriendIsDao(context.Background(), &imReq)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var detailRsp = rsp.UserDetailResp{
		AccountId:  resp.Detail.AccountId,
		Describe:   resp.Detail.Describe,
		HeadImgUrl: resp.Detail.HeadImgUrl,
		Username:   resp.Detail.Username,
		SchoolId:   resp.Detail.SchoolId,
		Gender:     resp.Detail.Gender,
		Age:        resp.Detail.Age,
		IsFriend:   imResp.IsFriend,
		Remark:     imResp.Remark,
		IsOfficial: resp.Detail.IsOfficial,
	}
	handle.SuccessResp(c, "", detailRsp)
}
