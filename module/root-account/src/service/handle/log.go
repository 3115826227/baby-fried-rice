package handle

import (
	"encoding/json"
	"github.com/3115826227/baby-fried-rice/module/root-account/src/config"
	"github.com/3115826227/baby-fried-rice/module/root-account/src/service/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RootLoginLog(c *gin.Context) {
	data, err := Get(config.Config.AccountDaoUrl + "/dao/account/root/login/log")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	var resp model.RspDaoRootLoginLog
	if err := json.Unmarshal(data, &resp); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	if resp.Code != 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	SuccessResp(c, "", resp.Data)
}

func AdminLoginLog(c *gin.Context) {
	data, err := Get(config.Config.AccountDaoUrl + "/dao/account/admin/login/log")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	var resp model.RspDaoAdminLoginLog
	if err := json.Unmarshal(data, &resp); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	if resp.Code != 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	SuccessResp(c, "", resp.Data)
}

func UserLoginLog(c *gin.Context) {
	data, err := Get(config.Config.AccountDaoUrl + "/dao/account/user/login/log")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	var resp model.RspDaoUserLoginLog
	if err := json.Unmarshal(data, &resp); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	if resp.Code != 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	SuccessResp(c, "", resp.Data)
}
