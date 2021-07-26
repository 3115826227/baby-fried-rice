package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/shopDao/db"
)

func GetCommodities(page, pageSize int64, id string) (commodities []tables.Commodity, total int64, err error) {
	var (
		offset = int((page - 1) * pageSize)
		limit  = int(pageSize)
	)
	var template = db.GetDB().GetDB().Model(&tables.Commodity{})
	if id != "" {
		template = template.Where("id = ?", id)
	}
	if err = template.Count(&total).Error; err != nil {
		return
	}
	err = template.Offset(offset).Limit(limit).Find(&commodities).Order("update_time desc").Error
	return
}

func GetCommodity(id string) (commodity tables.Commodity, err error) {
	err = db.GetDB().GetObject(map[string]interface{}{"id": id}, &commodity)
	return
}

func GetCommoditiesByIds(ids []string) (commodities []tables.Commodity, err error) {
	err = db.GetDB().GetDB().Where("id in (?)", ids).Find(commodities).Error
	return
}

func GetCommodityImageRelation(commodityId string) (relations []tables.CommodityImageRel, err error) {
	err = db.GetDB().GetDB().Where("commodity_id = ?", commodityId).Find(&relations).Error
	return
}
