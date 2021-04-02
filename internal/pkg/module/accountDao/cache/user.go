package cache

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
	"baby-fried-rice/internal/pkg/module/gateway/cache"
	"encoding/json"
	"fmt"
)

func AddUserDetail(detail tables.AccountUserDetail) (err error) {
	key := fmt.Sprintf("%v:%v", constant.AccountUserIDPrefix, detail.ID)
	detailBytes, err := json.Marshal(detail)
	if err != nil {
		return
	}
	return cache.GetCache().Add(key, string(detailBytes))
}

func GetUserDetail(userID string) (detail tables.AccountUserDetail, err error) {
	key := fmt.Sprintf("%v:%v", constant.AccountUserIDPrefix, userID)
	detailStr, err := cache.GetCache().Get(key)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(detailStr), &detail)
	return
}
