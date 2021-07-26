package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/manage/db"
)

func GetAdminByLogin(loginName string) (admin tables.AccountAdmin, err error) {
	var query = map[string]interface{}{
		"login_name": loginName,
	}
	err = db.GetAccountDB().GetObject(query, &admin)
	return
}
