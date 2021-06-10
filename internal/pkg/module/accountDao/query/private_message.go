package query

import (
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/module/accountDao/db"
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
)

func GetUserPrivateMessages(pms requests.UserPrivateMessagesReq) (messages []tables.UserPrivateMessage, err error) {
	pms.PageCommonReq.Validate()
	var (
		offset = int((pms.Page - 1) * pms.PageSize)
		limit  = int(pms.PageSize)
	)
	err = db.GetDB().GetDB().Where("receive_id = ?", pms.UserId).Order("receive_time").Offset(offset).Limit(limit).Find(&messages).Error
	return
}
