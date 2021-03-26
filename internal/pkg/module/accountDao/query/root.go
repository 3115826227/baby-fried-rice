package query

import (
	"baby-fried-rice/internal/pkg/module/accountDao/db"
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
)

func GetRootByLogin(loginName, password string) (root tables.AccountRoot, err error) {
	var query = map[string]interface{}{
		"login_name": loginName,
		"password":   password,
	}
	err = db.GetDB().GetObject(query, &root)
	return
}
