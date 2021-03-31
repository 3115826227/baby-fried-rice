package req

import "baby-fried-rice/internal/pkg/kit/constant"

type UserSendPrivateMessageReq struct {
	SendId          string                          `json:"send_id"`
	ReceiveId       string                          `json:"receive_id"`
	SendMessageType constant.SendPrivateMessageType `json:"send_message_type"`
	MessageType     int                             `json:"message_type"`
	MessageTitle    string                          `json:"message_title"`
	MessageContent  interface{}                     `json:"message_content"`
}

type UserPrivateMessagesReq struct {
	UserId string `json:"user_id"`
	PageCommonReq
}
