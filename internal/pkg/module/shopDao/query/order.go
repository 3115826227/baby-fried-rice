package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/shopDao/db"
)

func GetCommodityOrders(page, pageSize int64, accountId string) (commodityOrders []tables.CommodityOrder, total int64, err error) {
	var (
		offset = int((page - 1) * pageSize)
		limit  = int(pageSize)
	)
	var template = db.GetDB().GetDB().Model(&tables.CommodityOrder{}).Where("account_id = ?", accountId)
	if err = template.Count(&total).Error; err != nil {
		return
	}
	err = template.Offset(offset).Limit(limit).Find(&commodityOrders).Order("create_time desc").Error
	return
}

func GetCommodityOrderById(id, accountId string) (commodityOrder tables.CommodityOrder, err error) {
	err = db.GetDB().GetObject(map[string]interface{}{
		"id":         id,
		"account_id": accountId,
	}, &commodityOrder)
	return
}
