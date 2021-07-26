package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/manage/db"
)

type CommoditiesQueryParam struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
}

func GetCommodities(param CommoditiesQueryParam) (commodities []tables.Commodity, total int64, err error) {
	var (
		offset = int((param.Page - 1) * param.PageSize)
		limit  = int(param.PageSize)
	)
	template := db.GetShopDB().GetDB().Model(&tables.Commodity{})
	if err = template.Count(&total).Error; err != nil {
		return
	}
	err = template.Order("id").Offset(offset).Limit(limit).Find(&commodities).Error
	return
}
