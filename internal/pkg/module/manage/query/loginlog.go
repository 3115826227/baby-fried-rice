package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/manage/db"
)

type LoginLogsQueryParam struct {
	AccountId string `json:"account_id"`
	Page      int64  `json:"page"`
	PageSize  int64  `json:"page_size"`
}

func GetUserLoginLogs(param LoginLogsQueryParam) (logs []tables.AccountUserLoginLog, total int64, err error) {
	var (
		offset = int((param.Page - 1) * param.PageSize)
		limit  = int(param.PageSize)
	)
	template := db.GetAccountDB().GetDB().Model(&tables.AccountUserLoginLog{})
	if param.AccountId != "" {
		template = template.Where("account_id = ?", param.AccountId)
	}
	if err = template.Count(&total).Error; err != nil {
		return
	}
	err = template.Order("login_time desc").Offset(offset).Limit(limit).Find(&logs).Error
	return
}

func GetAdminLoginLogs(param LoginLogsQueryParam) (logs []tables.AccountAdminLoginLog, total int64, err error) {
	var (
		offset = int((param.Page - 1) * param.PageSize)
		limit  = int(param.PageSize)
	)
	template := db.GetAccountDB().GetDB().Model(&tables.AccountAdminLoginLog{})
	if param.AccountId != "" {
		template = template.Where("account_id = ?", param.AccountId)
	}
	if err = template.Count(&total).Error; err != nil {
		return
	}
	err = template.Order("login_time desc").Offset(offset).Limit(limit).Find(&logs).Error
	return
}
