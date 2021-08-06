package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/shopDao/db"
)

func GetCommodityCartById(accountId string) (relations []tables.CommodityCartRel, err error) {
	err = db.GetDB().GetDB().Where("account_id = ?", accountId).Find(&relations).Error
	return
}
