package db

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/module/accountDao/log"
	"errors"
	"time"
)

// 用户端仅支持一对一私信发送，不支持组发送和全局发送
func SendPrivateMessage(pm requests.UserSendPrivateMessageReq) (string, error) {
	var pmID = handle.GenerateSerialNumberByLen(constant.PrivateMessageIDDefaultLength)
	var now = time.Now()
	var pmc = tables.UserPrivateMessageContent{
		Content: pm.MessageContent,
	}
	pmc.ID = pmID
	pmc.CreatedAt, pmc.UpdatedAt = now, now
	var beans = make([]interface{}, 0)
	beans = append(beans, &pmc)
	switch pm.MessageSendType {
	case constant.SendPerson:
		var rpm = tables.UserPrivateMessage{
			Id:              pmID,
			SendId:          pm.SendId,
			ReceiveId:       pm.ReceiveId,
			MessageStatus:   0,
			ReceiveTime:     now,
			MessageType:     pm.MessageType,
			MessageSendType: int32(pm.MessageSendType),
			MessageTitle:    pm.MessageTitle,
		}
		beans = append(beans, &rpm)
	case constant.SendGroup:
		return pmID, errors.New("your have no permission to send private message to group")
	case constant.SendGlobal:
		return pmID, errors.New("your have no permission to send private message to global")
	default:
		return pmID, errors.New("send private message type is invalid")
	}
	return pmID, GetDB().CreateMulti(beans...)
}

func UpdatePrivateMessagesStatus(receiveId string, messageId []string) (err error) {
	return GetDB().GetDB().Model(&tables.UserPrivateMessage{}).
		Where("receive_id = ? and id in (?)", receiveId, messageId).
		Updates(map[string]interface{}{"status": 1}).Error
}

func DeletePrivateMessage(accountId string, messageId []string) (err error) {
	var checkPrivateMessages []tables.UserPrivateMessage
	var checkIds []string
	if err = GetDB().GetDB().Where("send_id = ? or receive = ?", accountId, accountId).
		Where("id in (?)", messageId).
		Find(&checkPrivateMessages).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	for _, pm := range checkPrivateMessages {
		checkIds = append(checkIds, pm.Id)
	}
	var tx = GetDB().GetDB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		if err = tx.Commit().Error; err != nil {
			log.Logger.Error(err.Error())
		}
	}()
	if err = tx.Model(&tables.UserPrivateMessage{}).Where("id in (?)", checkIds).
		Delete(&tables.UserPrivateMessage{}).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	return tx.Model(&tables.UserPrivateMessageContent{}).Where("id in (?)", checkIds).
		Delete(&tables.UserPrivateMessageContent{}).Error
}
