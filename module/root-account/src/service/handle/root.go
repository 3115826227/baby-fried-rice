package handle

import (
	"encoding/json"
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/root-account/src/config"
	"github.com/3115826227/baby-fried-rice/module/root-account/src/log"
	"github.com/3115826227/baby-fried-rice/module/root-account/src/redis"
	"github.com/3115826227/baby-fried-rice/module/root-account/src/service/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func RootLogin(c *gin.Context) {
	var err error
	var req model.ReqLogin
	if err = c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	req.LoginName = strings.TrimSpace(req.LoginName)
	req.Password = strings.TrimSpace(req.Password)
	req.Ip = c.GetHeader("IP")
	payload, err := json.Marshal(req)
	if err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}
	data, err := LoginPost(config.Config.AccountDaoUrl+"/dao/account/root/login", payload, c.Request.Header.Clone())
	if err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}
	var resp model.RspDaoRootLogin
	err = json.Unmarshal(data, &resp)
	if err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}
	if resp.Code != config.SuccessCode {
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}

	token, err := GenerateToken(resp.Data.ID, time.Now())
	if err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	var userInfo = model.RspUserData{
		UserId:    resp.Data.ID,
		LoginName: resp.Data.LoginName,
		Username:  resp.Data.Username,
	}
	var loginResult = model.LoginResult{
		UserInfo: userInfo,
		Token:    token,
	}
	var result = model.RspLogin{
		RspSuccess: model.RspSuccess{Code: 0},
		Data:       loginResult,
	}

	var userMeta = &model.UserMeta{
		UserId:   resp.Data.ID,
		Platform: "pc",
	}

	redis.AddAccountToken(fmt.Sprintf("%v:%v", TokenPrefix, token), userMeta.ToString())
	redis.AddAccountToken(userMeta.UserId, fmt.Sprintf("%v:%v", TokenPrefix, token))

	c.JSON(http.StatusOK, result)
}

func RootLogout(c *gin.Context) {
	userMeta := GetUserMeta(c)
	token, err := redis.GetAccountToken(userMeta.UserId)
	if err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	redis.DeleteAccountToken(token)
	redis.DeleteAccountToken(userMeta.UserId)
	SuccessResp(c, "", model.RspOkResponse{})
}

func RootRefresh(c *gin.Context) {

}

func RootDetail(c *gin.Context) {
	userMeta := GetUserMeta(c)
	data, err := Get(config.Config.AccountDaoUrl + "/dao/account/root/detail?id=" + userMeta.UserId)
	if err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}
	var rsp model.RspDaoRootDetail
	err = json.Unmarshal(data, &rsp)
	if err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}

	SuccessResp(c, "", rsp.Data)
}

func RootUpdatePassword(c *gin.Context) {

}
