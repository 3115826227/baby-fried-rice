package query

import (
	"baby-fried-rice/internal/pkg/module/accountDao/cache"
	"baby-fried-rice/internal/pkg/module/accountDao/db"
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
)

func GetUser(userID string) (user tables.AccountUser, err error) {
	if user, err = cache.GetUser(userID); err != nil {
		return db.GetUser(userID)
	}
	return
}
