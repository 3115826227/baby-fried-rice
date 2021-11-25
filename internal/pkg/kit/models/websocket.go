package models

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"
	"encoding/json"
)

type UserBaseInfo struct {
	AccountId  string `json:"account_id"`
	HeadImgUrl string `json:"head_img_url"`
	Username   string `json:"username"`
	IsOfficial bool   `json:"is_official"`
}

type WSMessageNotify struct {
	// 消息通知类型
	WSMessageNotifyType constant.WSMessageNotifyType `json:"ws_message_notify_type"`
	// 消息接收方
	Receive string `json:"receive"`
	// 消息
	WSMessage WSMessage `json:"ws_message"`
	// 发送时间
	Timestamp int64 `json:"timestamp"`
}

func (notify *WSMessageNotify) ToString() string {
	data, _ := json.Marshal(notify)
	return string(data)
}

type WSMessage struct {
	// 消息类型
	WSMessageType im.SessionNotifyType `json:"ws_message_type"`
	// 发送者
	Send rsp.User `json:"send"`
	// 空间消息 主要是推送
	Space *rsp.SpaceResp `json:"space,omitempty"`
	// 会话消息 既有接收也有推送
	SessionMessage *SessionMessage `json:"session_message,omitempty"`
	// 私信消息 主要是推送
	PrivateMessage rsp.UserPrivateMessageDetailResp `json:"private_message,omitempty"`
	// 直播消息
	LiveMessage *rsp.LiveRoomMessage `json:"live_message"`
}

func (message *WSMessage) ToString() string {
	data, _ := json.Marshal(message)
	return string(data)
}

type SessionMessage struct {
	// 信息类别
	SessionMessageType constant.SessionMessageType `json:"session_message_type"`
	// 新操作信息
	Operator rsp.Operator `json:"operator"`
	// 新会话信息
	Session rsp.Session `json:"session"`
	// 新会话消息信息
	Message rsp.Message `json:"message"`
	// 已读消息
	ReadMessage rsp.ReadMessage `json:"read_message"`
	// webrtc消息
	WebRtc rsp.WebRTC `json:"web_rtc"`
}
