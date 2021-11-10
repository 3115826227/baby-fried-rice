package cache

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

func getDialogSessionKey(accountId string) string {
	return fmt.Sprintf("%v:%v", constant.AccountDialogSessionUserIDPrefix, accountId)
}

// 查询对话框的会话列表
func GetSessionDialog(accountId string, page, pageSize int64) (sessionIds []int64, total int64, err error) {
	var (
		start = (page - 1) * pageSize
		stop  = start + pageSize
	)
	var results []string
	results, err = GetCache().ZRevRange(getDialogSessionKey(accountId), start, stop)
	for _, res := range results {
		var sessionId int64
		sessionId, err = strconv.ParseInt(res, 10, 64)
		if err != nil {
			return
		}
		sessionIds = append(sessionIds, sessionId)
	}
	return
}

// 设置对话框中的会话，如果会话已经存在，则获取当前时间戳，放到最前面
func SetSessionDialog(accountId string, sessionId int64) error {
	return GetCache().ZSet(getDialogSessionKey(accountId), redis.Z{
		Score:  float64(time.Now().UnixNano()),
		Member: sessionId,
	})
}

// 删除对话框的会话
func DeleteSessionDialog(accountId string, sessionId int64) error {
	return GetCache().ZRem(getDialogSessionKey(accountId), sessionId)
}
