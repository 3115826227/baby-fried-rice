package handle

import (
	"encoding/json"
	"github.com/3115826227/baby-fried-rice/module/admin-account/src/config"
	"github.com/3115826227/baby-fried-rice/module/admin-account/src/log"
	"github.com/3115826227/baby-fried-rice/module/admin-account/src/service/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func OrganizeGetHandle(c *gin.Context) {
	userMeta := GetUserMeta(c)
	data, err := Get(config.Config.AccountDaoUrl + "/dao/account/admin/organize?school_id=" + userMeta.SchoolId)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	var resp model.RespOrganize
	if err := json.Unmarshal(data, &resp); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	SuccessResp(c, "", resp.Data)
}

func OrganizeAddHandle(c *gin.Context) {
	var req model.ReqSchoolOrganizeAdd
	if err := c.BindJSON(&req); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	payload, err := json.Marshal(req)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	data, err := Post(config.Config.AccountDaoUrl+"/dao/account/admin/organize", payload)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	var resp model.RspOkResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	if resp.Code != 0 {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	SuccessResp(c, "", nil)
}
