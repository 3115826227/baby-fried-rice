package handle

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/module/accountDao/log"
	"baby-fried-rice/internal/pkg/module/accountDao/query"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserLoginLogsHandle(c *gin.Context) {
	logs, err := query.GetUserLoginLogs()
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", logs)
}

func AdminLoginLogsHandle(c *gin.Context) {

}

func RootLoginLogsHandle(c *gin.Context) {
	logs, err := query.GetRootLoginLogs()
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", logs)
}
