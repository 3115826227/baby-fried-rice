package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/manage/db"
)

type OrderQueryParam struct {
	OrderId   string `json:"order_id"`
	AccountId string `json:"account_id"`
	Page      int64  `json:"page"`
	PageSize  int64  `json:"page_size"`
}

func GetOrders(param OrderQueryParam) (orders []tables.CommodityOrder, total int64, err error) {
	var (
		offset = int((param.Page - 1) * param.PageSize)
		limit  = int(param.PageSize)
	)
	template := db.GetShopDB().GetDB().Model(&tables.CommodityOrder{})
	if param.OrderId != "" {
		template = template.Where("id = ?", param.OrderId)
	}
	if param.AccountId != "" {
		template = template.Where("account_id = ?", param.AccountId)
	}
	if err = template.Count(&total).Error; err != nil {
		return
	}
	err = template.Order("id").Offset(offset).Limit(limit).Find(&orders).Error
	return
}
