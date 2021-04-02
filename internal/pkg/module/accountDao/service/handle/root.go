package handle

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/module/accountDao/db"
	"baby-fried-rice/internal/pkg/module/accountDao/log"
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
	"baby-fried-rice/internal/pkg/module/accountDao/query"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func RootLoginHandle(c *gin.Context) {
	var err error
	var req requests.PasswordLoginReq
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	req.LoginName = strings.TrimSpace(req.LoginName)
	req.Password = strings.TrimSpace(req.Password)
	req.Ip = c.GetHeader(handle.HeaderIP)
	root, err := query.GetRootByLogin(req.LoginName, req.Password)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}

	go func() {
		var rootLoginLog = tables.AccountRootLoginLog{
			RootID:    root.ID,
			IP:        req.Ip,
			LoginTime: time.Now(),
		}
		db.AddRootLoginLog(rootLoginLog)
	}()
	handle.SuccessResp(c, "", root)
}
