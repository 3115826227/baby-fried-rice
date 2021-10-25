package cache

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

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

func DeleteUserCoin(coin tables.AccountUserCoin) (err error) {
	var key = fmt.Sprintf("%v:%v", constant.AccountUserCoinIDPrefix, coin.AccountID)
	return GetCache().Del(key)
}

func SetUserCoinRank(coin tables.AccountUserCoin) (err error) {
	scoreStr := fmt.Sprintf("%v.%v", coin.Coin, constant.MaxTimestamp-time.Now().Unix())
	var score float64
	if score, err = strconv.ParseFloat(scoreStr, 64); err != nil {
		return err
	}
	err = GetRedisClient().ZAdd(constant.AccountUserCoinRankKey, redis.Z{
		Score:  score,
		Member: coin.AccountID,
	}).Err()
	return
}
