package handle

import (
	"encoding/json"
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/admin-account/src/config"
	"github.com/3115826227/baby-fried-rice/module/admin-account/src/log"
	"github.com/3115826227/baby-fried-rice/module/admin-account/src/redis"
	"github.com/3115826227/baby-fried-rice/module/admin-account/src/service/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func AdminLoginHandle(c *gin.Context) {
	var err error
	var req model.ReqAdminLogin
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
	var resp model.RespAdminLogin
	var data []byte
	data, err = LoginPost(config.Config.AccountDaoUrl+"/dao/account/admin/login", payload, c.Request.Header.Clone())
	if err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	admin := resp.Data.Admin
	token, err := GenerateToken(admin.ID, time.Now())
	if err != nil {
		log.Logger.Warn(err.Error())
		c.JSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	var userInfo = model.RspUserData{
		UserId:    admin.ID,
		LoginName: admin.LoginName,
		Username:  admin.Username,
		SchoolId:  admin.SchoolID,
		IsSuper:   admin.Super,
	}
	roles := resp.Data.Roles
	roleIDs := make([]int64, 0)
	for _, r := range roles {
		roleIDs = append(roleIDs, r.ID)
	}
	permissions := resp.Data.Permissions

	var loginResult = model.LoginResult{
		UserInfo:   userInfo,
		Token:      token,
		Role:       roles,
		Permission: permissions,
	}
	var result = model.RspLogin{
		RspSuccess: model.RspSuccess{Code: 0},
		Data:       loginResult,
	}

	var isSuper string
	if admin.Super {
		isSuper = "1"
	}
	var userMeta = &model.UserMeta{
		UserId:   admin.ID,
		SchoolId: admin.SchoolID,
		IsSuper:  isSuper,
		Username: admin.Name,
	}

	redis.AddAccountToken(fmt.Sprintf("%v:%v", TokenPrefix, token), userMeta.ToString())

	c.JSON(http.StatusOK, result)
}

func AdminLogoutHandle(c *gin.Context) {
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
