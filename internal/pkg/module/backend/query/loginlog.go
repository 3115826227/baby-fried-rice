package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/backend/db"
)

func GetUserLoginLogs() (logs []tables.AccountUserLoginLog, err error) {
	err = db.GetAccountDB().GetDB().Find(&logs).Error
	return
}

func GetRootLoginLogs() (logs []tables.AccountRootLoginLog, err error) {
	err = db.GetAccountDB().GetDB().Find(&logs).Error
	return
}
