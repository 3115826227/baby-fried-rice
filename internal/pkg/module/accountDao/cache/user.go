package cache

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
	"encoding/json"
	"fmt"
)

func AddUserDetail(detail tables.AccountUserDetail) (err error) {
	var key = fmt.Sprintf("%v:%v", constant.AccountUserIDPrefix, detail.AccountID)
	var detailBytes []byte
	if detailBytes, err = json.Marshal(detail); err != nil {
		return
	}
	return GetCache().Add(key, string(detailBytes))
}

func GetUserDetail(accountId string) (detail tables.AccountUserDetail, err error) {
	var key = fmt.Sprintf("%v:%v", constant.AccountUserIDPrefix, accountId)
	var detailStr string
	if detailStr, err = GetCache().Get(key); err != nil {
		return
	}
	err = json.Unmarshal([]byte(detailStr), &detail)
	return
}
