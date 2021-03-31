package query

import (
	"baby-fried-rice/internal/pkg/module/accountDao/db"
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
)

func GetUserLoginLogs() (logs []tables.AccountUserLoginLog, err error) {
	err = db.GetDB().GetDB().Find(&logs).Error
	return
}

func GetRootLoginLogs() (logs []tables.AccountRootLoginLog, err error) {
	err = db.GetDB().GetDB().Find(&logs).Error
	return
}
