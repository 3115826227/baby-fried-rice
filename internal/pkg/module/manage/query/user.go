package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/manage/cache"
	"baby-fried-rice/internal/pkg/module/manage/db"
)

func GetUsers(page, pageSize int64) (details []tables.AccountUserDetail, total int64, err error) {
	var (
		offset = int((page - 1) * pageSize)
		limit  = int(pageSize)
	)
	template := db.GetAccountDB().GetDB().Model(&tables.AccountUserDetail{})
	if err = template.Count(&total).Error; err != nil {
		return
	}
	err = template.Order("account_id").Offset(offset).Limit(limit).Find(&details).Error
	return
}

func GetUsersByIds(ids []string) (details []tables.AccountUserDetail, err error) {
	details = make([]tables.AccountUserDetail, 0)
	for _, id := range ids {
		var detail tables.AccountUserDetail
		if detail, err = cache.GetUserDetail(id); err != nil {
			if err = db.GetAccountDB().GetObject(map[string]interface{}{"account_id": id}, &detail); err != nil {
				return
			}
		}
		details = append(details, detail)
	}
	return
}
