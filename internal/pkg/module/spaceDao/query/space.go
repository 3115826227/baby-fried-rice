package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/spaceDao/db"
)

func SpaceQuery(page, pageSize int64) (spaces []tables.Space, err error) {
	var (
		offset = int((page - 1) * pageSize)
		limit  = int(pageSize)
	)
	if err = db.GetDB().GetDB().Offset(offset).Limit(limit).Find(&spaces).Error; err != nil {
		return
	}
	return
}

func SpaceOptQuery(spaceId string) (relations []tables.SpaceOperatorRelation, err error) {
	if err = db.GetDB().GetDB().Where("space_id = ?", spaceId).Find(&relations).Error; err != nil {
		return
	}
	return
}

func SpaceCommentQuery(spaceId string) (relations []tables.SpaceCommentRelation, err error) {
	if err = db.GetDB().GetDB().Where("space_id = ?", spaceId).Find(&relations).Error; err != nil {
		return
	}
	return
}
