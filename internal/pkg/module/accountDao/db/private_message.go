package db

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/req"
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
	"encoding/json"
	"errors"
	"time"
)

func SendPrivateMessage(pm req.UserSendPrivateMessageReq) (err error) {
	var pmID = handle.GenerateSerialNumberByLen(constant.PrivateMessageIDDefaultLength)
	var mc []byte
	var now = time.Now()
	mc, err = json.Marshal(pm.MessageContent)
	if err != nil {
		return
	}
	var pmc = tables.UserPrivateMessageContent{
		CommonField: tables.CommonField{
			ID:        pmID,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Content:         string(mc),
		MessageSendType: int(pm.SendMessageType),
		MessageType:     pm.MessageType,
		MessageTitle:    pm.MessageTitle,
	}
	var spm = tables.UserPrivateMessage{
		MessageId:     pmID,
		SendId:        pm.SendId,
		ReceiveId:     "",
		MessageStatus: 0,
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
		}
		beans = append(beans, &rpm)
	case constant.SendGroup:
	case constant.SendGlobal:
	default:
		return errors.New("send private message type is invalid")
	}
	return GetDB().CreateMulti(beans...)
}
