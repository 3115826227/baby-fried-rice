package cache

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func UpdateUserOnlineStatus(accountId string, onlineType im.OnlineStatusType) error {
	var status = models.UserOnlineStatus{
		AccountId:  accountId,
		OnlineType: onlineType,
	}
	return GetCache().HMSet(constant.AccountUserOnlineStatusKey, map[string]interface{}{
		accountId: status.ToString(),
	})
}

func UpdateSessionWebRTCUserStatus(status rsp.SessionWebRTCUserStatus) error {
	var key = fmt.Sprintf("%v:%v", constant.AccountUserSessionWebRTCStatusKey, status.SessionId)
	return GetCache().HMSet(key, map[string]interface{}{
		status.AccountId: status.ToString(),
	})
}

func GetSessionWebRTC(sessionId int64) (statusMap map[string]rsp.SessionWebRTCUserStatus, err error) {
	var key = fmt.Sprintf("%v:%v", constant.AccountUserSessionWebRTCStatusKey, sessionId)
	var mp = make(map[string]string)
	statusMap = make(map[string]rsp.SessionWebRTCUserStatus)
	mp, err = GetCache().HGetAll(key)
	if err != nil {
		return
	}
	for k, v := range mp {
		var status rsp.SessionWebRTCUserStatus
		if err = json.Unmarshal([]byte(v), &status); err != nil {
			return
		}
		statusMap[k] = status
	}
	return
}

// 判断视频邀请是否有人接受
func JudgeSessionReceiveWebRTC(statusMap map[string]rsp.SessionWebRTCUserStatus) (ok bool, err error) {
	for _, v := range statusMap {
		if v.Status == im.SessionNotifyType_ReceiveVideoMessage {
			ok = true
			return
		}
	}
	ok = false
	return
}

func DeleteSessionWebRTC(sessionId int64) error {
	var key = fmt.Sprintf("%v:%v", constant.AccountUserSessionWebRTCStatusKey, sessionId)
	return GetCache().Del(key)
}

func SetSessionWebRTCStartTime(sessionId int64) error {
	var key = fmt.Sprintf("%v:%v", constant.AccountUserSessionWebRTCTimeKey, sessionId)
	return GetCache().Add(key, fmt.Sprintf("%v", time.Now().Unix()))
}

func GetSessionWebRTCStartTime(sessionId int64) (startTime int64, err error) {
	var key = fmt.Sprintf("%v:%v", constant.AccountUserSessionWebRTCTimeKey, sessionId)
	var value string
	if value, err = GetCache().Get(key); err != nil {
		return
	}
	startTime, err = strconv.ParseInt(value, 10, 64)
	return
}

func DeleteSessionWebRTCStartTime(sessionId int64) error {
	var key = fmt.Sprintf("%v:%v", constant.AccountUserSessionWebRTCTimeKey, sessionId)
	return GetCache().Del(key)
}
