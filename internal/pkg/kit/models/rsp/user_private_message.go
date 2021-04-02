package rsp

type UserPrivateMessagesResp struct {
	MessageId     string `json:"message_id"`
	SendId        string `json:"send_id"`
	SendName      string `json:"send_name"`
	ReceiveId     string `json:"receive_id"`
	MessageStatus int    `json:"message_status"`
	ReceiveTime   string `json:"receive_time"`
}
