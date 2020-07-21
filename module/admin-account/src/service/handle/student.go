package handle

import (
	"encoding/json"
	"github.com/3115826227/baby-fried-rice/module/admin-account/src/config"
	"github.com/3115826227/baby-fried-rice/module/admin-account/src/log"
	"github.com/3115826227/baby-fried-rice/module/admin-account/src/service/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func StudentGetHandle(c *gin.Context) {
	organize := c.Query("organize")
	if organize == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	data, err := Get(config.Config.AccountDaoUrl + "/dao/account/admin/student?organize=" + organize)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	var resp model.RespSchoolStudent
	if err := json.Unmarshal(data, &resp); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	SuccessResp(c, "", resp.Data)
}

func StudentAddHandle(c *gin.Context) {
	var req model.ReqSchoolStudentAdd
	if err := c.ShouldBind(&req); err != nil {
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
	data, err := Post(config.Config.AccountDaoUrl+"/dao/account/admin/student", payload)
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
