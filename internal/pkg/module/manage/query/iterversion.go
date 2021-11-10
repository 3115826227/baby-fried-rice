package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/manage/db"
)

func GetIterativeVersionByVersion(version string) (iv tables.IterativeVersion, err error) {
	err = db.GetAccountDB().GetObject(map[string]interface{}{"version": version}, &iv)
	return
}

type IterativeVersionQueryParam struct {
	LikeVersion string
	Status      *bool
	Page        int64
	PageSize    int64
}

func GetIterativeVersion(param IterativeVersionQueryParam) (ivs []tables.IterativeVersion, total int64, err error) {
	var (
		offset = int((param.Page - 1) * param.PageSize)
		limit  = int(param.PageSize)
	)
	template := db.GetAccountDB().GetDB().Model(&tables.IterativeVersion{})
	if param.LikeVersion != "" {
		template = template.Where("version like ?%", param.LikeVersion)
	}
	if param.Status != nil {
		template = template.Where("status = ?", param.Status)
	}
	if err = template.Count(&total).Error; err != nil {
		return
	}
	err = template.Order("version desc").Offset(offset).Limit(limit).Find(&ivs).Error
	return
}
