package handle

import (
	"encoding/json"
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/user-account/src/config"
	"github.com/3115826227/baby-fried-rice/module/user-account/src/log"
	"github.com/3115826227/baby-fried-rice/module/user-account/src/redis"
	"github.com/3115826227/baby-fried-rice/module/user-account/src/service/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func UserRegister(c *gin.Context) {
	var err error
	var req model.ReqUserRegister
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	payload, err := json.Marshal(req)
	if err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}
	_, err = Post(config.Config.AccountDaoUrl+"/dao/account/user/register", payload)
	if err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}

	c.JSON(http.StatusOK, model.RspOkResponse{})
}

/*
	用户实名制认证
*/
func UserVerify(c *gin.Context) {
	var err error
	var req model.ReqUserVerify
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	userMeta := GetUserMeta(c)
	req.UserId = userMeta.UserId
	payload, err := json.Marshal(req)
	if err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}
	_, err = Post(config.Config.AccountDaoUrl+"/dao/account/user/verify", payload)
	if err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}

	SuccessResp(c, "", nil)
}

func UserLogin(c *gin.Context) {
	var err error
	var req model.ReqLogin
	if err = c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	req.LoginName = strings.TrimSpace(req.LoginName)
	req.Password = strings.TrimSpace(req.Password)

	payload, err := json.Marshal(req)
	if err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}
	data, err := LoginPost(config.Config.AccountDaoUrl+"/dao/account/user/login", payload, c.Request.Header.Clone())
	if err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}

	var resp model.RespUserLogin
	if err := json.Unmarshal(data, &resp); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}
	if resp.Code != 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}

	token, err := GenerateToken(resp.Data.User.ID, time.Now())
	if err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	var userInfo = model.RspUserData{
		UserId:    resp.Data.User.ID,
		LoginName: resp.Data.User.LoginName,
		Username:  resp.Data.Detail.Username,
		SchoolId:  resp.Data.Detail.SchoolID,
	}
	var loginResult = model.LoginResult{
		UserInfo:   userInfo,
		Token:      token,
		Permission: make([]int, 0),
	}
	var result = model.RspLogin{
		RspSuccess: model.RspSuccess{Code: 0},
		Data:       loginResult,
	}

	var userMeta = &model.UserMeta{
		UserId:   resp.Data.User.ID,
		IsSuper:  "0",
		Username: resp.Data.Detail.Username,
		SchoolId: resp.Data.Detail.SchoolID,
		ReqId:    resp.Data.User.ID,
		Platform: "pc",
	}

	log.Logger.Info(fmt.Sprint("user login header info:", userMeta.ToString()))
	redis.AddAccountToken(fmt.Sprintf("%v:%v", TokenPrefix, token), userMeta.ToString())
	redis.AddAccountToken(userMeta.UserId, fmt.Sprintf("%v:%v", TokenPrefix, token))

	c.JSON(http.StatusOK, result)
}

/*
	用户退出登录
*/
func UserLogout(c *gin.Context) {
	userMeta := GetUserMeta(c)
	token, err := redis.GetAccountToken(userMeta.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	redis.DeleteAccountToken(token)
	redis.DeleteAccountToken(userMeta.UserId)
	SuccessResp(c, "", model.RspOkResponse{})
}

/*
	用户登录token刷新
*/
func UserRefresh(c *gin.Context) {
	userMeta := GetUserMeta(c)
	token, err := GenerateToken(userMeta.UserId, time.Now())
	if err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	oldToken, err := redis.GetAccountToken(userMeta.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	//删除旧token，更新新的token
	redis.DeleteAccountToken(oldToken)
	redis.AddAccountToken(fmt.Sprintf("%v:%v", TokenPrefix, token), userMeta.ToString())
	redis.AddAccountToken(userMeta.UserId, fmt.Sprintf("%v:%v", TokenPrefix, token))
}

func UserDetail(c *gin.Context) {
	userMeta := GetUserMeta(c)

	data, err := Get(config.Config.AccountDaoUrl + "/dao/account/user/verify_detail?id=" + userMeta.UserId)
	if err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}

	var resp model.RespUserVerifyDetail
	if err := json.Unmarshal(data, &resp); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}
	if resp.Code != 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}

	var rsp model.RspUserDetail
	var login = resp.Data.User
	var detail = resp.Data.Detail
	rsp.LoginName = login.LoginName
	rsp.UserId = detail.ID
	rsp.Username = detail.Username
	rsp.SchoolId = detail.SchoolID
	rsp.Gender = detail.Gender
	rsp.Age = detail.Age
	rsp.Verify = detail.Verify
	rsp.AccountId = detail.AccountId
	if detail.Verify {
		schoolDetail := resp.Data.School
		label := resp.Data.SchoolLabel
		rsp.School = label.School
		rsp.Faculty = label.Faculty
		rsp.Grade = label.Grade
		rsp.Major = label.Major
		rsp.Name = schoolDetail.Name
		rsp.Identify = schoolDetail.Identify
		rsp.Number = schoolDetail.Number
	}

	SuccessResp(c, "", rsp)
}

func UserPasswordUpdate(c *gin.Context) {

}

func UserUpdate(c *gin.Context) {

}

func UserInfoGet(c *gin.Context) {
	id := c.Query("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	data, err := Get(config.Config.AccountDaoUrl + "/dao/account/user/info?id=" + id)
	if err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	var resp model.RespUserInfo
	if err := json.Unmarshal(data, &resp); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	if resp.Code != 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	SuccessResp(c, "", resp.Data)
}
