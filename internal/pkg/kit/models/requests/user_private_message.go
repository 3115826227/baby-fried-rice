package requests

import "baby-fried-rice/internal/pkg/kit/constant"

type UserSendPrivateMessageReq struct {
	SendId          string                          `json:"send_id"`
	ReceiveId       string                          `json:"receive_id"`
	MessageSendType constant.SendPrivateMessageType `json:"message_send_type"`
	MessageType     int32                           `json:"message_type"`
	MessageTitle    string                          `json:"message_title"`
	MessageContent  string                          `json:"message_content"`
}

type UserPrivateMessagesReq struct {
	SendId string `json:"send_id"`
	PageCommonReq
}

type UpdatePrivateMessageStatusReq struct {
	AccountId  string   `json:"account_id"`
	MessageIds []string `json:"message_ids"`
}
