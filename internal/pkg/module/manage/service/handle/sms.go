package handle

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/module/manage/log"
	"baby-fried-rice/internal/pkg/module/manage/query"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SmsLogHandle(c *gin.Context) {
	reqPage, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var (
		logs  []tables.SendMessageLog
		total int64
	)
	var param = query.SmsLogsQueryParam{
		AccountId: c.Query(handle.QueryAccountId),
		Phone:     c.Query("phone"),
		Page:      reqPage.Page,
		PageSize:  reqPage.PageSize,
	}
	logs, total, err = query.GetSmsLog(param)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var list = make([]interface{}, 0)
	for _, l := range logs {
		list = append(list, rsp.SmsLogModelToRsp(l))
	}
	handle.SuccessListResp(c, "", list, total, reqPage.Page, reqPage.PageSize)
}
