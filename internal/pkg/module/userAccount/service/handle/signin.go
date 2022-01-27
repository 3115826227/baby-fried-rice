package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/userAccount/grpc"
	"baby-fried-rice/internal/pkg/module/userAccount/log"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 用户签到
func SignInHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	userClient, err := grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var resp *user.RspUserSignInDao
	resp, err = userClient.UserSignInDao(context.Background(), &user.ReqUserSignInDao{AccountId: userMeta.AccountId})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var response = rsp.UserSignInResp{
		Ok:       resp.Ok,
		Describe: resp.Describe,
		Coin:     resp.Coin,
	}
	handle.SuccessResp(c, "", response)
}

// 用户签到日志查询
func SignInLogHandle(c *gin.Context) {
	var year, month, day int
	var err error
	if year, err = strconv.Atoi(c.Query("year")); err != nil {
		year = 0
	}
	if month, err = strconv.Atoi(c.Query("month")); err != nil {
		month = 0
	}
	if day, err = strconv.Atoi(c.Query("day")); err != nil {
		day = 0
	}
	userMeta := handle.GetUserMeta(c)
	var userClient user.DaoUserClient
	userClient, err = grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var req = &user.ReqUserSignInLogQueryDao{
		AccountId: userMeta.AccountId,
		Year:      int64(year),
		Month:     int64(month),
		Day:       int64(day),
	}
	var resp *user.RspUserSignInLogQueryDao
	resp, err = userClient.UserSignInLogQueryDao(context.Background(), req)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var list = make([]rsp.UserSignInLogResp, 0)
	for _, sl := range resp.List {
		var signInLog = rsp.UserSignInLogResp{
			SignInType: constant.SignInType(sl.SignInType),
			Coin:       sl.Coin,
			Timestamp:  sl.Timestamp,
		}
		list = append(list, signInLog)
	}
	handle.SuccessResp(c, "", list)
}
