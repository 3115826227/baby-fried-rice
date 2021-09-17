package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/manage/db"
)

type CommoditiesQueryParam struct {
	SellType string
	LikeName string
	Status   string
	Page     int64
	PageSize int64
}

func GetCommodities(param CommoditiesQueryParam) (commodities []tables.Commodity, total int64, err error) {
	var (
		offset = int((param.Page - 1) * param.PageSize)
		limit  = int(param.PageSize)
	)
	template := db.GetShopDB().GetDB().Model(&tables.Commodity{})
	if param.SellType != "" {
		template = template.Where("sell_type = ?", param.SellType)
	}
	if param.LikeName != "" {
		template = template.Where("name like ?%", param.LikeName)
	}
	if param.Status != "" {
		template = template.Where("status = ?", param.Status)
	}
	if err = template.Count(&total).Error; err != nil {
		return
	}
	err = template.Order("id").Offset(offset).Limit(limit).Find(&commodities).Error
	return
}
