package cache

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"fmt"
)

func SetCommodityCart(accountId string, relation tables.CommodityCartRel) (err error) {
	key := fmt.Sprintf("%v:%v", constant.AccountUserCommodityCartPrefix, accountId)
	return GetCache().HSet(key, relation.CommodityId, relation)
}

func GetCommodityCart(accountId string) (relations []tables.CommodityCartRel, err error) {
	return
}
