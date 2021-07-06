package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/backend/db"
)

func GetRootByLogin(loginName, password string) (root tables.AccountRoot, err error) {
	var query = map[string]interface{}{
		"login_name": loginName,
		"password":   password,
	}
	err = db.GetAccountDB().GetObject(query, &root)
	return
}
