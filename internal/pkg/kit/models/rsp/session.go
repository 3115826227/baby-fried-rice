package rsp

import (
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"
	"encoding/json"
)

type Session struct {
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
	// 用户状态
	Users []User `json:"users"`
	// 对方正在输入
	Inputting bool `json:"inputting"`
}

type SessionDialog struct {
	// 会话id
	SessionId int64 `json:"session_id"`
	// 会话类型
	SessionType im.SessionType `json:"session_type"`
	// 会话名称
	Name string `json:"name"`
	// 会话等级
	Level im.SessionLevel `json:"level"`
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
	// 消息已读用户数
	ReadUserTotal int64 `json:"read_user_total"`
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

// 已读消息
type ReadMessage struct {
	// 已读消息id
	MessageId int64 `json:"message_id"`
	// 已读消息用户
	User string `json:"user"`
}

// 视频语音通话WebRTC消息
type WebRTC struct {
	// 邀请者id
	InviteAccount string   `json:"invite_account"`
	InviteUsers   []string `json:"invite_users"`
	Sdp           string   `json:"sdp"`
	SwapSdp       string   `json:"swap_sdp"`
	RemoteSdp     string   `json:"remote_sdp"`
	RemoteSwapSdp string   `json:"remote_swap_sdp"`
}

// 视频语音通话用户状态信息
type SessionWebRTCUserStatus struct {
	SessionId     int64                `json:"session_id"`
	AccountId     string               `json:"account_id"`
	Status        im.SessionNotifyType `json:"status"`
	Sdp           string               `json:"sdp"`
	SwapSdp       string               `json:"swap_sdp"`
	RemoteSdp     string               `json:"remote_sdp"`
	RemoteSwapSdp string               `json:"remote_swap_sdp"`
}

func (status *SessionWebRTCUserStatus) ToString() string {
	data, _ := json.Marshal(status)
	return string(data)
}

type MessageReadUsers struct {
	MessageId   int64  `json:"message_id"`
	ReadUsers   []User `json:"read_users"`
	UnreadUsers []User `json:"unread_users"`
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
	// 在线类型
	OnlineType im.OnlineStatusType `json:"online_type"`
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
