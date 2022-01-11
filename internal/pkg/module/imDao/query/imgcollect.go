package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/imDao/db"
)

func GetImgCollectRelations(accountId string, page, pageSize int64) (relations []tables.MessageImgCollectRelation, total int64, err error) {
	var (
		offset = int((page - 1) * pageSize)
		limit  = int(pageSize)
	)
	var template = db.GetDB().GetDB().Model(&tables.MessageImgCollectRelation{}).Where("account_id = ?", accountId)
	if err = template.Count(&total).Error; err != nil {
		return
	}
	err = template.Offset(offset).Limit(limit).Find(&relations).Order("id desc").Error
	return
}

func GetImgCollectRelation(accountId, img string) (relation tables.MessageImgCollectRelation, err error) {
	err = db.GetDB().GetDB().Where("account_id = ? and img = ?", accountId, img).First(&relation).Error
	return
}
