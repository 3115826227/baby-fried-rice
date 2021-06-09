package cache

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
	"encoding/json"
	"fmt"
)

func AddUserDetail(detail tables.AccountUserDetail) (err error) {
	key := fmt.Sprintf("%v:%v", constant.AccountUserIDPrefix, detail.AccountID)
	detailBytes, err := json.Marshal(detail)
	if err != nil {
		return
	}
	return GetCache().Add(key, string(detailBytes))
}

func GetUserDetail(accountId string) (detail tables.AccountUserDetail, err error) {
	key := fmt.Sprintf("%v:%v", constant.AccountUserIDPrefix, accountId)
	detailStr, err := GetCache().Get(key)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(detailStr), &detail)
	return
}
