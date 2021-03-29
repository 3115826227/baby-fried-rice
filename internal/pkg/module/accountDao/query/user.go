package query

import (
	"baby-fried-rice/internal/pkg/module/accountDao/cache"
	"baby-fried-rice/internal/pkg/module/accountDao/db"
	"baby-fried-rice/internal/pkg/module/accountDao/log"
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
)

func IsDuplicateAccountID(accountID string) bool {
	var count int64 = 0
	if err := db.GetDB().GetDB().Model(&tables.AccountUser{}).Where("id = ?", accountID).Count(&count).Error; err != nil {
		log.Logger.Error(err.Error())
		return true
	}
	return count != 0
}

func IsDuplicateLoginNameByUser(loginName string) bool {
	var count int64 = 0
	if err := db.GetDB().GetDB().Model(&tables.AccountUser{}).Where("login_name = ?", loginName).Count(&count).Error; err != nil {
		log.Logger.Error(err.Error())
		return true
	}
	return count != 0
}

func GetUserByLogin(loginName, password string) (root tables.AccountUser, err error) {
	var query = map[string]interface{}{
		"login_name": loginName,
		"password":   password,
	}
	err = db.GetDB().GetObject(query, &root)
	return
}

func GetUser(userID string) (user tables.AccountUser, err error) {
	if user, err = cache.GetUser(userID); err != nil {
		return db.GetUser(userID)
	}
	return
}
