package db

import (
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
)

func AddRoot(root tables.AccountRoot) error {
	return GetDB().CreateObject(&root)
}

func AddRootLoginLog(rootLoginLog tables.AccountRootLoginLog) error {
	return GetDB().CreateObject(&rootLoginLog)
}
