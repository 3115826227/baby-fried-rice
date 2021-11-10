package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/imDao/db"
)

func GetMessage(id, sessionId int64) (message tables.Message, err error) {
	err = db.GetDB().GetDB().Where("id = ? and session_id = ?", id, sessionId).First(&message).Error
	return
}

func GetMessages(ids []int64, sessionId int64) (messages []tables.Message, err error) {
	err = db.GetDB().GetDB().Where("id in (?) and session_id = ?", ids, sessionId).Find(&messages).Error
	return
}

func GetMessageRelation(id, sessionId int64) (relations []tables.MessageUserRelation, err error) {
	err = db.GetDB().GetDB().Where("message_id = ? and session_id = ?", id, sessionId).Find(&relations).Error
	return
}

func GetMessageReadUserTotal(id, sessionId int64, accountId string) (count int64, err error) {
	err = db.GetDB().GetDB().Model(&tables.MessageUserRelation{}).Where("message_id = ? and session_id = ? and receive != ?", id, sessionId, accountId).Count(&count).Error
	return
}
