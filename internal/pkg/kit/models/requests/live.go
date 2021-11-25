package requests

import "baby-fried-rice/internal/pkg/kit/rpc/pbservices/live"

type ReqUpdateOriginLiveRoom struct {
	Status live.LiveRoomStatus `json:"status"`
	Sdp    string              `json:"sdp"`
}

type ReqUpdateUserLiveRoomOpt struct {
	LiveRoomId string                   `json:"live_room_id"`
	Opt        live.LiveRoomUserOptType `json:"opt"`
	RemoteSdp  string                   `json:"remote_sdp"`
}
