package models

import (
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"
	"encoding/json"
)

type UserOnlineStatus struct {
	AccountId   string              `json:"account_id"`
	VideoStatus bool                `json:"video_status"`
	OnlineType  im.OnlineStatusType `json:"online_type"`
}

func (status *UserOnlineStatus) ToString() string {
	data, _ := json.Marshal(status)
	return string(data)
}

type UserPhoneCode struct {
	AccountId string `json:"account_id"`
	Phone     string `json:"phone"`
	Code      string `json:"code"`
}
