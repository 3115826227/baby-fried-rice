package handle

import (
	"encoding/json"
	"github.com/3115826227/baby-fried-rice/module/admin-account/src/config"
	"github.com/3115826227/baby-fried-rice/module/admin-account/src/log"
	"github.com/3115826227/baby-fried-rice/module/admin-account/src/service/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func OfficialGroupAddHandle(c *gin.Context) {
	userMeta := GetUserMeta(c)
	var req model.ReqOfficialGroupAdd
	if err := c.BindJSON(&req); err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}

	req.Admin = userMeta.UserId
	payload, err := json.Marshal(req)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	data, err := Post(config.Config.ImUrl+"/api/im/official/group", payload)
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

func OfficialGroupAddStudentHandle(c *gin.Context) {

}

func OfficialGroupDeleteStudentHandle(c *gin.Context) {

}

func OfficialGroupGetHandle(c *gin.Context) {

}
