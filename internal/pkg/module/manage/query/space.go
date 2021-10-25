package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/manage/db"
)

type SpacesQueryParam struct {
	Id          string
	AccountId   string
	VisitorType string
	AuditStatus string
	Page        int64
	PageSize    int64
}

func GetSpaces(param SpacesQueryParam) (spaces []tables.Space, total int64, err error) {
	var (
		offset = int((param.Page - 1) * param.PageSize)
		limit  = int(param.PageSize)
	)
	template := db.GetSpaceDB().GetDB().Model(&tables.Space{})
	if param.Id != "" {
		template = template.Where("id = ?", param.Id)
	}
	if param.AccountId != "" {
		template = template.Where("origin = ?", param.AccountId)
	}
	if param.VisitorType != "" {
		template = template.Where("visitor_type = ?", param.VisitorType)
	}
	if param.AuditStatus != "" {
		template = template.Where("audit_status = ?", param.AuditStatus)
	}
	if err = template.Count(&total).Error; err != nil {
		return
	}
	err = template.Order("create_time desc").Offset(offset).Limit(limit).Find(&spaces).Error
	return
}
