package requests

import "baby-fried-rice/internal/pkg/kit/constant"

type UserSendPrivateMessageReq struct {
	SendId          string                          `json:"send_id"`
	ReceiveId       string                          `json:"receive_id; binding:required"`
	MessageSendType constant.SendPrivateMessageType `json:"message_send_type; binding:required"`
	MessageType     int32                           `json:"message_type"`
	MessageTitle    string                          `json:"message_title; binding:required"`
	MessageContent  string                          `json:"message_content; binding:required"`
}

type UserPrivateMessagesReq struct {
	AccountId string `json:"account_id"`
	SendId    string `json:"send_id"`
	PageCommonReq
}

type UpdatePrivateMessageStatusReq struct {
	MessageIds []string `json:"message_ids"`
}
