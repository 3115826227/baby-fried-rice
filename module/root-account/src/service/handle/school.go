package handle

import (
	"encoding/json"
	"github.com/3115826227/baby-fried-rice/module/root-account/src/config"
	"github.com/3115826227/baby-fried-rice/module/root-account/src/log"
	"github.com/3115826227/baby-fried-rice/module/root-account/src/service/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SchoolAdd(c *gin.Context) {
	var req model.ReqSchoolAdd
	if err := c.BindJSON(&req); err != nil {
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
	_, err = Post(config.Config.AccountDaoUrl+"/dao/account/school", payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	SuccessResp(c, "", nil)
}

func SchoolGet(c *gin.Context) {
	data, err := Get(config.Config.AccountDaoUrl + "/dao/account/school")
	if err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}

	var resp model.RespSchool
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
