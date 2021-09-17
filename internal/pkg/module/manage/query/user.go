package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/manage/cache"
	"baby-fried-rice/internal/pkg/module/manage/db"
)

type UserQueryParam struct {
	AccountId    string
	LikeUsername string
	Page         int64
	PageSize     int64
}

func GetUsers(param UserQueryParam) (details []tables.AccountUserDetail, total int64, err error) {
	var (
		offset = int((param.Page - 1) * param.PageSize)
		limit  = int(param.PageSize)
	)
	template := db.GetAccountDB().GetDB().Model(&tables.AccountUserDetail{})
	if param.AccountId != "" {
		template = template.Where("account_id = ?", param.AccountId)
	}
	if param.LikeUsername != "" {
		template = template.Where("username like ?%", param.LikeUsername)
	}
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
