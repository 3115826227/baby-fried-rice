package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/manage/db"
)

func GetUserCoins(page, pageSize int64) (userCoins []tables.AccountUserCoin, total int64, err error) {
	var (
		offset = int((page - 1) * pageSize)
		limit  = int(pageSize)
	)
	template := db.GetAccountDB().GetDB().Model(&tables.AccountUserCoin{})
	if err = template.Count(&total).Error; err != nil {
		return
	}
	err = template.Order("account_id").Offset(offset).Limit(limit).Find(&userCoins).Error
	return
}

func GetUserCoinsByIds(ids []string) (userCoins []tables.AccountUserCoin, err error) {
	err = db.GetAccountDB().GetDB().Where("account_id in (?)", ids).Find(&userCoins).Error
	return
}
