package rsp

import "baby-fried-rice/internal/pkg/kit/rpc/pbservices/live"

type UpdateOriginLiveRoomResp struct {
	LiveRoom LiveRoom `json:"live_room"`
	SwapSdp  string   `json:"swap_sdp"`
}

type LiveRoom struct {
	LiveRoomId string `json:"live_room_id"`
	Origin     User   `json:"origin"`
	Status     live.LiveRoomStatus
	UserTotal  int64 `json:"user_total"`
}

type LiveRoomUserResp struct {
	List     []User `json:"list"`
	Page     int64  `json:"page"`
	PageSize int64  `json:"page_size"`
}

type LiveRoomDetailResp struct {
	LiveRoom
	OnlineTime int64 `json:"online_time"`
}

type UpdateUserLiveRoomOptResp struct {
	RemoteSwapSdp string `json:"remote_swap_sdp"`
}

type LiveRoomMessage struct {
	MessageId     int64                    `json:"message_id"`
	MessageType   live.LiveRoomMessageType `json:"message_type"`
	Send          User                     `json:"send"`
	Content       string                   `json:"content"`
	SendTimestamp int64                    `json:"send_timestamp"`
}
