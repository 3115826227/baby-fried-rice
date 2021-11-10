package cache

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"
	"encoding/json"
)

func GetUserOnlineStatus(accountIds []string) (statusMap map[string]im.OnlineStatusType, err error) {
	var values []interface{}
	values, err = GetCache().HMGet(constant.AccountUserOnlineStatusKey, accountIds...)
	if err != nil {
		return
	}
	statusMap = make(map[string]im.OnlineStatusType)
	for _, value := range values {
		if value != nil {
			var status models.UserOnlineStatus
			if err = json.Unmarshal([]byte(value.(string)), &status); err != nil {
				return
			}
			statusMap[status.AccountId] = status.OnlineType
		}
	}
	return
}
