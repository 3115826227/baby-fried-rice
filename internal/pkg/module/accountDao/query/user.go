package query

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/module/accountDao/cache"
	"baby-fried-rice/internal/pkg/module/accountDao/db"
	"baby-fried-rice/internal/pkg/module/accountDao/log"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"strings"
)

func IsDuplicateAccountID(accountID string) bool {
	if _, err := cache.GetUserDetail(accountID); err != nil {
		var count int64 = 0
		if err = db.GetDB().GetDB().Model(&tables.AccountUser{}).Where("account_id = ?", accountID).Count(&count).Error; err != nil {
			log.Logger.Error(err.Error())
			return true
		}
		return count != 0
	}
	return false

}

// 校验用户登录名是否重复
func IsDuplicateLoginNameByUser(loginName string) bool {
	var count int64 = 0
	if err := db.GetDB().GetDB().Model(&tables.AccountUser{}).Where("login_name = ?", loginName).Count(&count).Error; err != nil {
		log.Logger.Error(err.Error())
		return true
	}
	return count != 0
}

func GetUserByLogin(loginName string) (root tables.AccountUser, err error) {
	var query = map[string]interface{}{
		"login_name": loginName,
	}
	err = db.GetDB().GetObject(query, &root)
	return
}

func GetUserDetail(accountId string) (detail tables.AccountUserDetail, err error) {
	if detail, err = cache.GetUserDetail(accountId); err != nil {
		err = db.GetDB().GetObject(map[string]interface{}{"account_id": accountId}, &detail)
		if err != nil {
			return
		}
		go cache.SetUserDetail(detail)
	}
	return
}

func GetUserDetails(ids []string) (details []tables.AccountUserDetail, err error) {
	var detailMap = make(map[string]tables.AccountUserDetail, 0)
	var failedIds = make([]string, 0)
	var cacheDetails []tables.AccountUserDetail
	// 获取缓存中有效用户信息
	cacheDetails, err = cache.GetUserByIds(ids)
	for _, detail := range cacheDetails {
		detailMap[detail.AccountID] = detail
	}
	details = make([]tables.AccountUserDetail, 0)
	// 过滤找出缓存未命中的用户信息
	for _, id := range ids {
		if _, exist := detailMap[id]; !exist {
			failedIds = append(failedIds, id)
		}
	}
	// 从数据库中批量查找未命中的用户信息，并更新到缓存中
	var failedDetails []tables.AccountUserDetail
	if err = db.GetDB().GetDB().Where("account_id in (?)", failedIds).Find(&failedDetails).Error; err != nil {
		return
	}
	if err = cache.SetUserDetails(failedDetails); err != nil {
		return
	}
	for _, detail := range failedDetails {
		detailMap[detail.AccountID] = detail
	}
	for _, id := range ids {
		details = append(details, detailMap[id])
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

func GetUserCoin(accountId string) (coin tables.AccountUserCoin, err error) {
	if coin, err = cache.GetUserCoin(accountId); err != nil {
		err = db.GetDB().GetObject(map[string]interface{}{"account_id": accountId}, &coin)
		if err != nil {
			return
		}
		go cache.SetUserCoin(coin)
	}
	return
}

func GetUserCoinLog(accountId string, pageReq requests.PageCommonReq) (logs []tables.AccountUserCoinLog, total int64, err error) {
	var (
		offset = int((pageReq.Page - 1) * pageReq.PageSize)
		limit  = int(pageReq.PageSize)
	)
	template := db.GetDB().GetDB().Model(&tables.AccountUserCoinLog{}).Where("account_id = ?", accountId)
	if err = template.Count(&total).Error; err != nil {
		return
	}
	err = template.Order("timestamp desc").Offset(offset).Limit(limit).Find(&logs).Error
	return
}

// 获取用户排名
func GetUserCoinRank(accountId string) (rank int64, sameCoinUsers int64, coin int64, err error) {
	if rank, err = cache.GetRedisClient().ZRevRank(constant.AccountUserCoinRankKey, accountId).Result(); err != nil {
		return
	}
	rank += 1
	var userCoin tables.AccountUserCoin
	if userCoin, err = GetUserCoin(accountId); err != nil {
		return
	}
	coin = userCoin.Coin
	sameCoinUsers, err = cache.GetRedisClient().ZCount(constant.AccountUserCoinRankKey, fmt.Sprintf("%v", coin), fmt.Sprintf("%v", coin+1)).Result()
	return
}

// 获取topN排名榜单
func GetUserCoinRankBoard(topN int64) (usersMap []map[string]interface{}, err error) {
	var members []redis.Z
	members, err = cache.GetRedisClient().ZRevRangeWithScores(constant.AccountUserCoinRankKey, 0, topN-1).Result()
	if err != nil {
		return
	}
	usersMap = make([]map[string]interface{}, 0)
	for index, member := range members {
		var accountId = member.Member.(string)
		var userMap = make(map[string]interface{})
		userMap["account_id"] = accountId
		userMap["rank"] = int64(index + 1)
		userMap["timestamp"], userMap["coin"], err = ParseUserCoinFromScore(member.Score)
		if err != nil {
			return
		}
		usersMap = append(usersMap, userMap)
	}
	return
}

func ParseUserCoinFromScore(score float64) (timestamp int64, coin int64, err error) {
	scoreStr := strconv.FormatFloat(score, 'f', 12, 64)
	slice := strings.Split(scoreStr, ".")
	coin, err = strconv.ParseInt(slice[0], 10, 64)
	if err != nil {
		return
	}
	var ts int64
	ts, err = strconv.ParseInt(slice[1], 10, 64)
	if err != nil {
		return
	}
	timestamp = constant.MaxTimestamp - ts
	return
}

// 查询用户最近一条签到信息
func GetUserLatestSignIn(accountId string) (signInLog tables.AccountUserSignInLog, err error) {
	if signInLog, err = cache.GetUserSignInLatestLog(accountId); err != nil {
		if err = db.GetDB().GetDB().Where("account_id = ?", accountId).Order("timestamp desc").First(&signInLog).Error; err != nil {
			return
		}
		go cache.SetUserSignInLatestLog(signInLog)
	}
	return
}
