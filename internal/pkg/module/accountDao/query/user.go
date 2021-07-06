package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/module/accountDao/cache"
	"baby-fried-rice/internal/pkg/module/accountDao/db"
	"baby-fried-rice/internal/pkg/module/accountDao/log"
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

func GetUserDetail(accountId string) (detail tables.AccountUserDetail, err error) {
	if detail, err = cache.GetUserDetail(accountId); err != nil {
		err = db.GetDB().GetObject(map[string]interface{}{"account_id": accountId}, &detail)
	}
	return
}

func GetUsers(ids []string) (details []tables.AccountUserDetail, err error) {
	details = make([]tables.AccountUserDetail, 0)
	for _, id := range ids {
		var detail tables.AccountUserDetail
		if detail, err = cache.GetUserDetail(id); err != nil {
			if err = db.GetDB().GetObject(map[string]interface{}{"account_id": id}, &detail); err != nil {
				return
			}
		}
		details = append(details, detail)
	}
	return
}

func GetAll() (ids []string, err error) {
	var users []tables.AccountUserDetail
	if err = db.GetDB().GetDB().Select("account_id").Find(&users).Error; err != nil {
		return
	}
	for _, user := range users {
		ids = append(ids, user.AccountID)
	}
	return
}
