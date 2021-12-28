package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/module/accountDao/db"
)

func GetUserPrivateMessages(pms requests.UserPrivateMessagesReq) (messages []tables.UserPrivateMessage, total int64, err error) {
	pms.PageCommonReq.Validate()
	var (
		offset = int((pms.Page - 1) * pms.PageSize)
		limit  = int(pms.PageSize)
	)
	template := db.GetDB().GetDB().Model(&tables.UserPrivateMessage{}).Where("receive_id = ?", pms.AccountId)
	if pms.SendId != "" {
		template = template.Where("send_id = ?", pms.SendId)
	}
	if err = template.Count(&total).Error; err != nil {
		return
	}
	if err = template.Order("receive_time").Offset(offset).Limit(limit).Find(&messages).Error; err != nil {
		return
	}
	return
}

func GetUserPrivateMessageDetail(accountId, messageId string) (msg tables.UserPrivateMessage, detail tables.UserPrivateMessageContent, err error) {
	if err = db.GetDB().GetDB().Where("id = ? and receive_id = ?", messageId, accountId).First(&msg).Error; err != nil {
		return
	}
	if err = db.GetDB().GetDB().Model(&tables.UserPrivateMessage{}).Where("id = ? and receive_id = ?", messageId, accountId).Update("message_status", 1).Error; err != nil {
		return
	}
	if err = db.GetDB().GetDB().Where("id = ?", messageId).First(&detail).Error; err != nil {
		return
	}
	return
}
