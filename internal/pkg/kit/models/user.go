package models

import (
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"
	"encoding/json"
)

type UserOnlineStatus struct {
	AccountId  string
	OnlineType im.OnlineStatusType
}

func (status *UserOnlineStatus) ToString() string {
	data, _ := json.Marshal(status)
	return string(data)
}
