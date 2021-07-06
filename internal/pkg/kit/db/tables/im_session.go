package tables

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"
)

// 会话用户关系表
type SessionUserRelation struct {
	// 会话id
	SessionID int64 `gorm:"column:session_id;unique_index:session_user_relation" json:"session_id"`
	// 用户id
	UserID string `gorm:"column:user_id;unique_index:session_user_relation" json:"user_id"`
	// 加入时间
	JoinTime int64 `gorm:"column:join_time;" json:"join_time"`
	// 用户在会话中的备注
	Remark string `gorm:"column:remark" json:"remark"`
}

func (table *SessionUserRelation) TableName() string {
	return "baby_im_session_user_rel"
}

// 会话详情表
type Session struct {
	// 会话id
	ID int64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	// 会话名称
	Name string `gorm:"column:name;not null" json:"name"`
	// 会话类型
	SessionType im.SessionType `gorm:"column:session_type;" json:"session_type"`
	// 会话加入权限类型
	JoinPermissionType im.SessionJoinPermissionType `gorm:"column:join_permission_type" json:"join_permission_type"`
	// 会话创建者id
	Origin string `gorm:"column:origin;not null" json:"origin"`
	// 会话等级
	Level im.SessionLevel `gorm:"column:level;not null" json:"level"`
	// 会话人数限制
	UserLimit constant.SessionUserLimit `gorm:"column:user_limit;not null" json:"user_limit"`
	// 会话创建时间
	CreateTime int64 `json:"create_time"`
	// 会话信息更新时间
	UpdateTime int64 `json:"update_time"`
}

func (table *Session) TableName() string {
	return "baby_im_session"
}

// 消息用户关系表
type MessageUserRelation struct {
	// 消息id
	MessageID int64 `gorm:"column:message_id;primaryKey" json:"message_id"`
	// 会话id
	SessionID int64 `gorm:"column:session_id;" json:"session_id"`
	// 接收者id
	Receive string `gorm:"column:receive;primaryKey" json:"receive"`
	// 接收者消息读取状态
	Read bool `gorm:"column:read;not null" json:"read"`
	// 消息发送时间
	SendTimestamp int64 `gorm:"column:send_timestamp;not null" json:"send_time"`
}

func (table *MessageUserRelation) TableName() string {
	return "baby_im_message_user_rel"
}

// 消息表
type Message struct {
	// 消息id
	ID int64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	// 会话id
	SessionID int64 `gorm:"column:session_id;primaryKey" json:"session_id"`
	// 消息类型
	MessageType int32 `gorm:"column:message_type;not null" json:"message_type"`
	// 发送者id
	Send string `gorm:"column:send;not null" json:"send"`
	// 消息内容
	Content string `gorm:"column:content" json:"content"`
	// 消息发送时间
	SendTimestamp int64 `gorm:"column:send_timestamp;not null" json:"send_time"`
}

func (table *Message) TableName() string {
	return "baby_im_message"
}
