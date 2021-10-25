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
	var bytes []byte
	if bytes, err = json.Marshal(detail); err != nil {
		return
	}
	return GetCache().HSet(constant.AccountUserIDPrefix, detail.AccountID, string(bytes))
}

func GetUserDetail(accountId string) (detail tables.AccountUserDetail, err error) {
	var detailStr string
	if detailStr, err = GetCache().HGet(constant.AccountUserIDPrefix, accountId); err != nil {
		return
	}
	err = json.Unmarshal([]byte(detailStr), &detail)
	return
}

func DeleteUserDetail(ids ...string) error {
	return GetCache().HDel(constant.AccountUserIDPrefix, ids...)
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

func DeleteUserCoin(accountId string) error {
	var key = fmt.Sprintf("%v:%v", constant.AccountUserCoinIDPrefix, accountId)
	return GetCache().Del(key)
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

func SetUserDetails(details []tables.AccountUserDetail) error {
	if len(details) == 0 {
		return nil
	}
	var mp = make(map[string]interface{})
	for _, detail := range details {
		data, err := json.Marshal(detail)
		if err != nil {
			return err
		}
		mp[detail.AccountID] = string(data)
	}
	return GetCache().HMSet(constant.AccountUserIDPrefix, mp)
}

func GetUserByIds(ids []string) (details []tables.AccountUserDetail, err error) {
	var resps []interface{}
	if resps, err = GetCache().HMGet(constant.AccountUserIDPrefix, ids...); err != nil {
		return
	}
	for _, resp := range resps {
		var detail tables.AccountUserDetail
		var data []byte
		if data, err = json.Marshal(resp); err != nil {
			return
		}
		if err = json.Unmarshal(data, &detail); err != nil {
			return
		}
		details = append(details, detail)
	}
	return
}

func SetUserSignInLatestLog(signInLog tables.AccountUserSignInLog) error {
	key := fmt.Sprintf("%v:%v", constant.AccountUserSignInPrefix, signInLog.AccountId)
	data, err := json.Marshal(signInLog)
	if err != nil {
		return err
	}
	return GetCache().Add(key, string(data))
}

func GetUserSignInLatestLog(accountId string) (signInLog tables.AccountUserSignInLog, err error) {
	key := fmt.Sprintf("%v:%v", constant.AccountUserSignInPrefix, accountId)
	var signInLogStr string
	if signInLogStr, err = GetCache().Get(key); err != nil {
		return
	}
	err = json.Unmarshal([]byte(signInLogStr), &signInLog)
	return
}

func DeleteUserSignInLatestLog(accountId string) error {
	key := fmt.Sprintf("%v:%v", constant.AccountUserSignInPrefix, accountId)
	return GetCache().Del(key)
}
