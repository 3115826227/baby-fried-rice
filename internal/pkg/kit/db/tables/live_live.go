package tables

import "baby-fried-rice/internal/pkg/kit/rpc/pbservices/live"

// 直播房间信息
type LiveRoom struct {
	CommonField

	// 直播房间主播
	Origin string
	// 直播房间状态
	Status live.LiveRoomStatus
}

func (table *LiveRoom) TableName() string {
	return "baby_live__room"
}

// 直播房间用户关联表
type LiveRoomUserRelation struct {
	LiveRoomID    string `gorm:"column:live_room_id;unique_index:live_room_account"`
	AccountID     string `gorm:"column:account_id;unique_index:live_room_account"`
	JoinTimestamp int64  `gorm:"column:join_timestamp"`
}

func (table *LiveRoomUserRelation) TableName() string {
	return "baby_live_room_user_rel"
}

// 直播房间消息表
type LiveRoomMessage struct {
	// 消息id
	ID int64 `gorm:"column:id;primaryKey;autoIncrement"`
	// 直播房间id
	LiveRoomID string `gorm:"column:live_room_id;primaryKey"`
	// 消息类型
	MessageType live.LiveRoomMessageType `gorm:"column:message_type;not null"`
	// 发送者id
	Send string `gorm:"column:send;not null"`
	// 消息内容
	Content string `gorm:"column:content"`
	// 消息发送时间
	SendTimestamp int64 `gorm:"column:send_timestamp;not null"`
}

func (table *LiveRoomMessage) TableName() string {
	return "baby_live_room_message"
}

// 直播房间活动信息
type LiveRoomActivity struct {
	CommonIntField

	// 直播房间ID
	LiveRoomID string
	// 直播活动主题
	Title string
	// 直播活动描述
	Describe string
	// 直播活动开始时间
	StartTimestamp int64
	// 直播活动结束时间
	EndTimestamp int64
}

func (table *LiveRoomActivity) TableName() string {
	return "baby_live_room_activity"
}
