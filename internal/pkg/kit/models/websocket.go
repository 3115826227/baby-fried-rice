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
	WSMessageType im.SessionMessageType `json:"ws_message_type"`
	// 发送者
	Send UserBaseInfo `json:"send"`
	// 空间消息
	Space *Space `json:"space,omitempty"`
	// 会话消息
	SessionMessage *SessionMessage `json:"session_message,omitempty"`
	// 私信消息
	PrivateMessage *PrivateMessage `json:"private_message,omitempty"`
}

func (message *WSMessage) ToString() string {
	data, _ := json.Marshal(message)
	return string(data)
}

type Space struct {
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
}

type PrivateMessage struct {
}
