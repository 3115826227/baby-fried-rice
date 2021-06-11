package tables

import (
	"baby-fried-rice/internal/pkg/kit/models/tables"
	"time"
)

type UserPrivateMessage struct {
	// 消息id
	MessageId string `gorm:"unique_index:message_send_receive_id"`
	// 消息发送者
	SendId string `gorm:"unique_index:message_send_receive_id"`
	// 消息接收者
	ReceiveId string `gorm:"unique_index:message_send_receive_id"`
	// 接收者消息状态 0-未读 1-已读
	MessageStatus uint32
	// 接收时间
	ReceiveTime time.Time `gorm:"column:receive_time" json:"receive_time"`
	// 消息发送类型 1-person 2-group 3-global 4-notify
	MessageSendType int32
	// 消息标题
	MessageTitle string
	// 消息类型
	MessageType int32
}

func (table *UserPrivateMessage) TableName() string {
	return "baby_user_private_message"
}

func (table *UserPrivateMessage) Get() interface{} {
	return *table
}

type UserPrivateMessageContent struct {
	tables.CommonField
	// 消息内容
	Content string
}

func (table *UserPrivateMessageContent) TableName() string {
	return "baby_user_private_message_content"
}

func (table *UserPrivateMessageContent) Get() interface{} {
	return *table
}
