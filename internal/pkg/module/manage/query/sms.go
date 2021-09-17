package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/manage/db"
)

type SmsLogsQueryParam struct {
	AccountId string
	Phone     string
	Page      int64
	PageSize  int64
}

func GetSmsLog(param SmsLogsQueryParam) (logs []tables.SendMessageLog, total int64, err error) {
	var (
		offset = int((param.Page - 1) * param.PageSize)
		limit  = int(param.PageSize)
	)
	template := db.GetSmsDB().GetDB().Model(&tables.SendMessageLog{})
	if param.AccountId != "" {
		template = template.Where("account_id = ?", param.AccountId)
	}
	if param.Phone != "" {
		template = template.Where("phone = ?", param.Phone)
	}
	if err = template.Count(&total).Error; err != nil {
		return
	}
	err = template.Order("id desc").Offset(offset).Limit(limit).Find(&logs).Error
	return
}
