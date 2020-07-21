package handle

import (
	"encoding/json"
	"github.com/3115826227/baby-fried-rice/module/root-account/src/config"
	"github.com/3115826227/baby-fried-rice/module/root-account/src/log"
	"github.com/3115826227/baby-fried-rice/module/root-account/src/service/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserGet(c *gin.Context)  {
	data, err := Get(config.Config.AccountDaoUrl + "/dao/account/user/detail")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	var resp model.RespUser
	if err := json.Unmarshal(data, &resp); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, sysErrResponse)
		return
	}
	SuccessResp(c, "", resp.Data)
}
