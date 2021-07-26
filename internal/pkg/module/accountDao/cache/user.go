package cache

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
)

func SetUserDetail(detail tables.AccountUserDetail) (err error) {
	var key = fmt.Sprintf("%v:%v", constant.AccountUserIDPrefix, detail.AccountID)
	var bytes []byte
	if bytes, err = json.Marshal(detail); err != nil {
		return
	}
	return GetCache().Add(key, string(bytes))
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

func SetUserCoin(coin tables.AccountUserCoin) (err error) {
	var key = fmt.Sprintf("%v:%v", constant.AccountUserCoinIDPrefix, coin.AccountID)
	var bytes []byte
	if bytes, err = json.Marshal(coin); err != nil {
		return
	}
	return GetCache().Add(key, string(bytes))
}

func GetUserCoin(accountId string) (coin tables.AccountUserCoin, err error) {
	var key = fmt.Sprintf("%v:%v", constant.AccountUserCoinIDPrefix, accountId)
	var coinStr string
	if coinStr, err = GetCache().Get(key); err != nil {
		return
	}
	err = json.Unmarshal([]byte(coinStr), &coin)
	return
}

// 更新用户积分排名
func SetUserCoinRank(userCoinLog tables.AccountUserCoinLog, coin int64) (err error) {
	if userCoinLog.Coin < 0 {
		return
	}
	var newCoinTotal = userCoinLog.Coin + coin
	var scoreStr = fmt.Sprintf("%v.%v", newCoinTotal, constant.MaxTimestamp-userCoinLog.Timestamp)
	var score float64
	if score, err = strconv.ParseFloat(scoreStr, 64); err != nil {
		return
	}
	err = GetRedisClient().ZAdd(constant.AccountUserCoinRankKey, redis.Z{
		Score:  score,
		Member: userCoinLog.AccountID,
	}).Err()
	return
}
