package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/module/manage/cache"
	"baby-fried-rice/internal/pkg/module/manage/db"
	"baby-fried-rice/internal/pkg/module/manage/log"
	"baby-fried-rice/internal/pkg/module/manage/query"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 系统发放积分
func SystemGiveawayUserCoinHandle(c *gin.Context) {
	var req requests.UserCoinGiveawayReq
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	userCoins, err := query.GetUserCoinsByIds(req.Ids)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var userCoinMap = make(map[string]tables.AccountUserCoin)
	for _, uc := range userCoins {
		userCoinMap[uc.AccountID] = uc
	}
	var updateSql = fmt.Sprintf("insert into %v (account_id, coin, coin_total, update_timestamp) values ",
		(&tables.AccountUserCoin{}).TableName())
	var now = time.Now().Unix()
	for index, id := range req.Ids {
		if index != 0 {
			updateSql = fmt.Sprintf("%v ,", updateSql)
		}
		var newUserCoin tables.AccountUserCoin
		if uc, exist := userCoinMap[id]; exist {
			newUserCoin = tables.AccountUserCoin{
				AccountID:       uc.AccountID,
				Coin:            uc.Coin + req.Coin,
				CoinTotal:       uc.CoinTotal + req.Coin,
				UpdateTimestamp: now,
			}
		} else {
			newUserCoin = tables.AccountUserCoin{
				AccountID:       id,
				Coin:            req.Coin,
				CoinTotal:       req.Coin,
				UpdateTimestamp: now,
			}
		}
		userCoinMap[id] = newUserCoin
		updateSql = fmt.Sprintf("%v ('%v', %v, %v, %v)", updateSql,
			newUserCoin.AccountID, newUserCoin.Coin, newUserCoin.CoinTotal, newUserCoin.UpdateTimestamp)
	}
	updateSql = fmt.Sprintf("%v ON DUPLICATE KEY UPDATE coin = VALUES(coin) , coin_total = VALUES(coin_total) , update_timestamp = VALUES(update_timestamp);", updateSql)
	var tx = db.GetAccountDB().GetDB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		// 删除缓存，更新缓存的rank
		for _, coin := range userCoinMap {
			if err = cache.DeleteUserCoin(coin); err != nil {
				log.Logger.Error(err.Error())
			}
			if err = cache.SetUserCoinRank(coin); err != nil {
				log.Logger.Error(err.Error())
			}
		}
		tx.Commit()
	}()
	if err = tx.Exec(updateSql).Error; err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var coinLogs = make([]tables.AccountUserCoinLog, 0)
	for _, accountId := range req.Ids {
		var coinLog = tables.AccountUserCoinLog{
			AccountID: accountId,
			Coin:      req.Coin,
			CoinType:  constant.SystemGiveawayCoinType,
			Timestamp: now,
		}
		coinLogs = append(coinLogs, coinLog)
	}
	if err = tx.Create(&coinLogs).Error; err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 用户积分列表查询
func UserCoinHandle(c *gin.Context) {
	reqPage, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var (
		details   []tables.AccountUserDetail
		userCoins []tables.AccountUserCoin
		total     int64
	)
	details, total, err = query.GetUsers(reqPage.Page, reqPage.PageSize)
	var ids = make([]string, 0)
	for _, d := range details {
		ids = append(ids, d.AccountID)
	}
	userCoins, err = query.GetUserCoinsByIds(ids)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var userCoinMap = make(map[string]tables.AccountUserCoin)
	for _, uc := range userCoins {
		userCoinMap[uc.AccountID] = uc
	}
	var list = make([]rsp.UserCoin, 0)
	for _, d := range details {
		var userCoin = rsp.UserCoin{
			User: rsp.User{
				AccountID:  d.AccountID,
				Username:   d.Username,
				HeadImgUrl: d.HeadImgUrl,
			},
		}
		if uc, exist := userCoinMap[d.AccountID]; exist {
			userCoin.Coin = uc.Coin
			userCoin.CoinTotal = uc.CoinTotal
			userCoin.UpdateTimestamp = uc.UpdateTimestamp
		}
		list = append(list, userCoin)
	}
	var response = rsp.UserCoinResp{
		List:     list,
		Page:     reqPage.Page,
		PageSize: reqPage.PageSize,
		Total:    total,
	}
	handle.SuccessResp(c, "", response)
}
