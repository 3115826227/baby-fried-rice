package query

import (
	"baby-fried-rice/internal/pkg/kit/models/req"
	"baby-fried-rice/internal/pkg/module/accountDao/db"
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
)

func GetUserPrivateMessages(pms req.UserPrivateMessagesReq) (messages []tables.UserPrivateMessage, err error) {
	pms.PageCommonReq.Validate()
	err = db.GetDB().GetDB().Where("receive_id = ?", pms.UserId).Find(&messages).Error
	return
}
