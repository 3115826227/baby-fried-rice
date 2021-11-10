package cache

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"
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
