package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/liveDao/db"
)

type LiveRoomParam struct {
	Page     int64
	PageSize int64
}

func GetLiveRooms(param LiveRoomParam) (liveRooms []tables.LiveRoom, err error) {
	var (
		offset = int((param.Page - 1) * param.PageSize)
		limit  = int(param.PageSize)
	)
	var template = db.GetDB().GetDB().Model(&tables.LiveRoom{})
	err = template.Offset(offset).Limit(limit).Find(&liveRooms).Order("id desc").Error
	return
}

func LiveRoomUserRelation(accountId, id string) (rel tables.LiveRoomUserRelation, err error) {
	err = db.GetDB().GetDB().Where("live_room_id = ? and account_id = ?", id, accountId).First(&rel).Error
	return
}

func LiveRoomUsers(id string, page, pageSize int64) (users []string, err error) {
	var relations []tables.LiveRoomUserRelation
	var (
		offset = int((page - 1) * pageSize)
		limit  = int(pageSize)
	)
	if err = db.GetDB().GetDB().Where("id = ?", id).Offset(offset).Limit(limit).Find(&relations).Error; err != nil {
		return
	}
	for _, rel := range relations {
		users = append(users, rel.AccountID)
	}
	return
}

func LiveRoomById(id string) (liveRoom tables.LiveRoom, err error) {
	err = db.GetDB().GetDB().Where("id = ?", id).First(&liveRoom).Error
	return
}

func LiveRoomByOrigin(origin string) (liveRoom tables.LiveRoom, err error) {
	err = db.GetDB().GetDB().Where("origin = ?", origin).First(&liveRoom).Error
	return
}

func LiveRoomUserTotal(liveRoomId string) (total int64, err error) {
	err = db.GetDB().GetDB().Model(&tables.LiveRoomUserRelation{}).Where("live_room_id = ?", liveRoomId).Count(&total).Error
	return
}
