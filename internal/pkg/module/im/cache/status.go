package cache

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"
	"encoding/json"
	"fmt"
	"time"
)

func GetUserOnlineStatus(accountId string) (status models.UserOnlineStatus, err error) {
	var value string
	value, err = GetCache().HGet(constant.AccountUserOnlineStatusKey, accountId)
	if err != nil {
		return
	}
	if err = json.Unmarshal([]byte(value), &status); err != nil {
		return
	}
	return
}

func UpdateUserVideoStatus(accountId string, videoStatus bool) error {
	status, err := GetUserOnlineStatus(accountId)
	if err != nil {
		return err
	}
	status.VideoStatus = videoStatus
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

func RemoveSessionWebRTCUserStatus(sessionId int64, accountId string) error {
	var key = fmt.Sprintf("%v:%v", constant.AccountUserSessionWebRTCStatusKey, sessionId)
	return GetCache().HDel(key, accountId)

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

func SetSessionWebRTCTimeInfo(sessionId int64, video bool) error {
	var key = fmt.Sprintf("%v:%v", constant.AccountUserSessionWebRTCTimeKey, sessionId)
	var timeInfo = rsp.SessionWebRTCTimeInfo{
		SessionId: sessionId,
		Video:     video,
		StartTime: time.Now().Unix(),
	}
	data, err := json.Marshal(timeInfo)
	if err != nil {
		return err
	}
	return GetCache().Add(key, string(data))
}

func GetSessionWebRTCTimeInfo(sessionId int64) (timeInfo rsp.SessionWebRTCTimeInfo, err error) {
	var key = fmt.Sprintf("%v:%v", constant.AccountUserSessionWebRTCTimeKey, sessionId)
	var value string
	if value, err = GetCache().Get(key); err != nil {
		return
	}
	if err = json.Unmarshal([]byte(value), &timeInfo); err != nil {
		return
	}
	return
}

func DeleteSessionWebRTCTimeInfo(sessionId int64) error {
	var key = fmt.Sprintf("%v:%v", constant.AccountUserSessionWebRTCTimeKey, sessionId)
	return GetCache().Del(key)
}
