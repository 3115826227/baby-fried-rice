package rsp

type UserPrivateMessage struct {
	MessageId     string `json:"message_id"`
	Send          User   `json:"send"`
	ReceiveId     string `json:"receive_id"`
	MessageStatus uint32 `json:"message_status"`
	ReceiveTime   string `json:"receive_time"`
	Title         string `json:"title"`
}

type UserPrivateMessageDetailResp struct {
	UserPrivateMessage
	Content string `json:"content"`
}
