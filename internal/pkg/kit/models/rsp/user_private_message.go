package rsp

import (
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/privatemessage"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
)

type UserPrivateMessage struct {
	MessageId     string `json:"message_id"`
	Send          User   `json:"send"`
	ReceiveId     string `json:"receive_id"`
	MessageStatus uint32 `json:"message_status"`
	ReceiveTime   string `json:"receive_time"`
	Title         string `json:"title"`
}

type UserPrivateMessageDetailResp struct {
	UserPrivateMessage
	Content string `json:"content"`
}

func PrivateMessagePbConvertToRsp(pm *privatemessage.PrivateMessageQueryDao, detail *user.UserDao) UserPrivateMessage {
	return UserPrivateMessage{
		MessageId: pm.Id,
		Send: User{
			AccountID:  detail.Id,
			Username:   detail.Username,
			HeadImgUrl: detail.HeadImgUrl,
			IsOfficial: detail.IsOfficial,
		},
		ReceiveId:     pm.ReceiveId,
		MessageStatus: pm.Status,
		ReceiveTime:   pm.CreateTime,
		Title:         pm.Title,
	}
}
