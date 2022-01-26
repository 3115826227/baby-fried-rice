package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/module/manage/db"
	"baby-fried-rice/internal/pkg/module/manage/log"
	"baby-fried-rice/internal/pkg/module/manage/query"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 管理用户登录日志查询
func AdminLoginLogHandle(c *gin.Context) {
	reqPage, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	var (
		logs  []tables.AccountAdminLoginLog
		total int64
	)
	var param = query.LoginLogsQueryParam{
		AccountId: c.Query(handle.QueryAccountId),
		Page:      reqPage.Page,
		PageSize:  reqPage.PageSize,
	}
	logs, total, err = query.GetAdminLoginLogs(param)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var ids = make([]string, 0)
	for _, l := range logs {
		ids = append(ids, l.AccountId)
	}
	var admins []tables.AccountAdmin
	if err = db.GetAccountDB().GetDB().Where("id in (?)", ids).Find(&admins).Error; err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var idsMap = make(map[string]tables.AccountAdmin)
	for _, admin := range admins {
		idsMap[admin.ID] = admin
	}
	var list = make([]interface{}, 0)
	for _, l := range logs {
		list = append(list, rsp.UserLoginLogResp{
			ID: l.ID,
			User: rsp.User{
				AccountID:  l.AccountId,
				Username:   idsMap[l.AccountId].Username,
				HeadImgUrl: idsMap[l.AccountId].HeadImgUrl,
			},
			LoginCount:     l.LoginCount,
			IP:             l.IP,
			LoginTimestamp: l.LoginTime.Unix(),
		})
	}
	handle.SuccessListResp(c, "", list, total, reqPage.Page, reqPage.PageSize)
}

// 用户登录日志查询
func UserLoginLogHandle(c *gin.Context) {
	reqPage, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, constant.ParamErrResponse)
		return
	}
	var (
		logs  []tables.AccountUserLoginLog
		total int64
	)
	var param = query.LoginLogsQueryParam{
		AccountId: c.Query(handle.QueryAccountId),
		Page:      reqPage.Page,
		PageSize:  reqPage.PageSize,
	}
	logs, total, err = query.GetUserLoginLogs(param)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var ids = make([]string, 0)
	for _, l := range logs {
		ids = append(ids, l.AccountId)
	}
	var idsMap = make(map[string]tables.AccountUserDetail)
	if idsMap, err = query.GetUsersByIds(ids); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, constant.SysErrResponse)
		return
	}
	var list = make([]interface{}, 0)
	for _, l := range logs {
		list = append(list, rsp.UserLoginLogResp{
			ID: l.ID,
			User: rsp.User{
				AccountID:  idsMap[l.AccountId].AccountID,
				Username:   idsMap[l.AccountId].Username,
				HeadImgUrl: idsMap[l.AccountId].HeadImgUrl,
			},
			LoginCount:     l.LoginCount,
			IP:             l.IP,
			LoginTimestamp: l.LoginTime.Unix(),
		})
	}
	handle.SuccessListResp(c, "", list, total, reqPage.Page, reqPage.PageSize)
}
