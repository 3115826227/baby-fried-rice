package rsp

import "baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"

type SessionsResp struct {
	Sessions []Session `json:"sessions"`
}

type Session struct {
	// 会话id
	SessionId int64 `json:"session_id"`
	// 会话类型
	SessionType im.SessionType `json:"session_type"`
	// 会话名称
	Name string `json:"name"`
	// 未读消息数
	Unread int64 `json:"unread"`
	// 最近一条消息
	LatestMessage *Message `json:"latest_message,omitempty"`
}

type SessionDetail struct {
	// 会话id
	SessionId int64 `json:"session_id"`
	// 会话类型
	SessionType im.SessionType `json:"session_type"`
	// 会话名称
	Name string `json:"name"`
	// 会话等级
	Level im.SessionLevel `json:"level"`
	// 会话创建者
	Origin string `json:"origin"`
	// 会话加入人员
	Joins []User `json:"joins"`
	// 会话加入权限
	JoinPermissionType im.SessionJoinPermissionType `json:"join_permission_type"`
	// 会话创建时间
	CreateTime int64 `json:"create_time"`
}

type Message struct {
	// 会话id
	SessionId int64 `json:"session_id"`
	// 消息id
	MessageId int64 `json:"message_id"`
	// 消息类型
	MessageType im.SessionMessageType `json:"message_type"`
	// 发送者
	Send User `json:"send"`
	// 接收者id
	Receive string `json:"receive"`
	// 消息内容
	Content string `json:"content"`
	// 发送时间
	SendTimestamp int64 `json:"send_timestamp"`
	// 读取状态
	ReadStatus bool `json:"read_status"`
}

type Messages []Message

func (messages Messages) Len() int {
	return len(messages)
}

func (messages Messages) Swap(i, j int) {
	messages[i], messages[j] = messages[j], messages[i]
}

// 根据message_id顺序排序
func (messages Messages) Less(i, j int) bool {
	return messages[i].MessageId < messages[j].MessageId
}

type SessionMessageResp struct {
	Messages []Message `json:"messages"`
	Page     int64     `json:"page"`
	PageSize int64     `json:"page_size"`
}

type UserManageResp struct {
	// 用户id
	AccountId string `json:"account_id"`
	// 用户添加为好友设置的权限
	AddFriendPermissionType int32 `json:"add_friend_permission_type"`
	// 更新时间
	UpdateTimestamp int64 `json:"update_timestamp"`
}

type FriendResp struct {
	Friends []Friend `json:"friends"`
}

type Friend struct {
	// 用户id
	AccountId string `json:"account_id"`
	// 好友备注
	Remark string `json:"remark"`
	// 是否在黑名单中
	BlackList bool `json:"black_list"`
	// 成为好友的时间
	Timestamp int64 `json:"timestamp"`
}

type Friends []Friend

func (friends Friends) Len() int {
	return len(friends)
}

func (friends Friends) Swap(i, j int) {
	friends[i], friends[j] = friends[j], friends[i]
}

// 根据好友备注的字典序顺序排序
func (friends Friends) Less(i, j int) bool {
	return friends[i].Remark < friends[j].Remark
}

type OperatorResp struct {
	Operators []Operator `json:"operators"`
	Page      int64      `json:"page"`
	PageSize  int64      `json:"page_size"`
	Total     int64      `json:"total"`
}

type Operator struct {
	// 操作id
	OperatorId int64 `json:"operator_id"`
	// 操作用户id
	Origin User `json:"origin"`
	// 操作接收用户id
	Receive string `json:"receive"`
	// 操作类型
	OptType im.OptType `json:"opt_type"`
	// 操作内容
	Content string `json:"content"`
	// 是否需要确认
	NeedConfirm bool `json:"need_confirm"`
	// 确认情况
	Confirm int64 `json:"confirm"`
	// 操作时间
	OptTimestamp int64 `json:"opt_timestamp"`
	// 操作接收用户读取状态 0-不可见 1-未读 2-已读
	ReceiveReadStatus int `json:"receive_read_status,omitempty"`
}
