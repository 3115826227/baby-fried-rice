package rsp

type UserPrivateMessage struct {
	MessageId     string `json:"message_id"`
	Send          User   `json:"send"`
	ReceiveId     string `json:"receive_id"`
	MessageStatus uint32 `json:"message_status"`
	ReceiveTime   string `json:"receive_time"`
	Title         string `json:"title"`
}

type UserPrivateMessagesResp struct {
	List     []UserPrivateMessage `json:"list"`
	Page     int64                `json:"page"`
	PageSize int64                `json:"page_size"`
	Total    int64                `json:"total"`
}

type UserPrivateMessageDetailResp struct {
	UserPrivateMessage
	Content string `json:"content"`
}
