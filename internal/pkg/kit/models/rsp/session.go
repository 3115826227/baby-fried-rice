package rsp

import "baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"

type SessionsResp struct {
	Sessions []Session `json:"sessions"`
}

type Session struct {
	SessionId     int64          `json:"session_id"`
	SessionType   im.SessionType `json:"session_type"`
	Name          string         `json:"name"`
	Origin        string         `json:"origin"`
	Unread        int64          `json:"unread"`
	CreateTime    string         `json:"create_time"`
	LatestMessage *Message       `json:"latest_message,omitempty"`
	Joins         []string       `json:"joins"`
}

type Message struct {
	SessionId     int64                 `json:"session_id"`
	MessageId     int64                 `json:"message_id"`
	MessageType   im.SessionMessageType `json:"message_type"`
	Send          User                  `json:"send"`
	Receive       string                `json:"receive"`
	Content []byte                `json:"content"`
	SendTimestamp int64                 `json:"send_timestamp"`
	ReadStatus    bool                  `json:"read_status"`
}
