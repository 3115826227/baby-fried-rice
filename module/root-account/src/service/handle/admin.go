package handle

import (
	"encoding/json"
	"github.com/3115826227/baby-fried-rice/module/root-account/src/config"
	"github.com/3115826227/baby-fried-rice/module/root-account/src/log"
	"github.com/3115826227/baby-fried-rice/module/root-account/src/service/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminInit(c *gin.Context) {
	var err error
	var req model.ReqAdminInit
	if err = c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	payload, err := json.Marshal(req)
	if err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, sysErrResponse)
		return
	}
	_, err = Post(config.Config.AccountDaoUrl+"/dao/account/root/admin/init", payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	SuccessResp(c, "", nil)
}

func AdminAdd(c *gin.Context) {

}

func AdminGet(c *gin.Context) {
	schoolId := c.Query("school_id")
	if schoolId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	data, err := Get(config.Config.AccountDaoUrl + "/dao/account/admin?school_id=" + schoolId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	var resp model.RespAdmin
	if err := json.Unmarshal(data, &resp); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	SuccessResp(c, "", resp.Data)
}
