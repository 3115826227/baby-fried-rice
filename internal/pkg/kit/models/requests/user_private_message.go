package requests

import "baby-fried-rice/internal/pkg/kit/constant"

type UserSendPrivateMessageReq struct {
	SendId          string                          `json:"send_id"`
	ReceiveId       string                          `json:"receive_id"`
	SendMessageType constant.SendPrivateMessageType `json:"send_message_type"`
	MessageType     int                             `json:"message_type"`
	MessageTitle    string                          `json:"message_title"`
	MessageContent  string                          `json:"message_content"`
}

type UserPrivateMessagesReq struct {
	UserId string `json:"user_id"`
	PageCommonReq
}

type UpdatePrivateMessageStatusReq struct {
	ReceiveId  string   `json:"receive_id"`
	MessageIds []string `json:"message_ids"`
}
