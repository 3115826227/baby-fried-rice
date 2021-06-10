package db

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
	"errors"
	"time"
)

func SendPrivateMessage(pm requests.UserSendPrivateMessageReq) (err error) {
	var pmID = handle.GenerateSerialNumberByLen(constant.PrivateMessageIDDefaultLength)
	var now = time.Now()
	var pmc = tables.UserPrivateMessageContent{
		Content:         pm.MessageContent,
		MessageSendType: int(pm.SendMessageType),
		MessageType:     pm.MessageType,
		MessageTitle:    pm.MessageTitle,
	}
	pmc.ID = pmID
	pmc.CreatedAt = now
	pmc.UpdatedAt = now
	var spm = tables.UserPrivateMessage{
		MessageId:     pmID,
		SendId:        pm.SendId,
		ReceiveId:     "",
		MessageStatus: 0,
		ReceiveTime:   now,
	}
	var beans = make([]interface{}, 0)
	beans = append(beans, &pmc)
	beans = append(beans, &spm)
	switch pm.SendMessageType {
	case constant.SendPerson:
		var rpm = tables.UserPrivateMessage{
			MessageId:     spm.MessageId,
			SendId:        spm.ReceiveId,
			ReceiveId:     spm.SendId,
			MessageStatus: 0,
			ReceiveTime:   now,
		}
		beans = append(beans, &rpm)
	case constant.SendGroup:
	case constant.SendGlobal:
	default:
		return errors.New("send private message type is invalid")
	}
	return GetDB().CreateMulti(beans...)
}

func UpdatePrivateMessagesStatus(receiveId string, messageId []string) (err error) {
	return GetDB().GetDB().Model(&tables.UserPrivateMessage{}).
		Where("receive_id = ? and message_id in (?)", receiveId, messageId).
		Updates(map[string]interface{}{"status": 1}).Error
}
