package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/liveDao/db"
)

type LiveRoomParamMessage struct {
	StartTimestamp int64
	Page           int64
	PageSize       int64
}

func GetLiveRoomMessage(param LiveRoomParamMessage) (messages []tables.LiveRoomMessage, err error) {
	var (
		offset = int((param.Page - 1) * param.PageSize)
		limit  = int(param.PageSize)
	)
	var template = db.GetDB().GetDB().Model(&tables.LiveRoomMessage{})
	err = template.Offset(offset).Limit(limit).Where("send_timestamp > ?", param.StartTimestamp).Find(&messages).Order("send_timestamp desc").Error
	return
}
