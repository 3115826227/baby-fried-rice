package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"
	"baby-fried-rice/internal/pkg/module/im/log"
	"time"
)

// 邀请加入会话通知
func sendInviteNotify(session rsp.Session, send rsp.User, accountId string) {
	var notify = models.WSMessageNotify{
		WSMessageNotifyType: constant.SessionMessageNotify,
		Receive:             accountId,
		WSMessage: models.WSMessage{
			WSMessageType: im.SessionNotifyType_InviteNotify,
			Send:          send,
			SessionMessage: &models.SessionMessage{
				SessionMessageType: constant.SessionMessage,
				Session:            session,
			},
		},
		Timestamp: time.Now().Unix(),
	}
	if err := mq.Send(topic, notify.ToString()); err != nil {
		log.Logger.Error(err.Error())
		return
	}
}

// 消息已读通知
func sendMessageReadNotify(readMessage rsp.ReadMessage, send rsp.User, accountId string) {
	var notify = models.WSMessageNotify{
		WSMessageNotifyType: constant.SessionMessageNotify,
		Receive:             accountId,
		WSMessage: models.WSMessage{
			Send:          send,
			WSMessageType: im.SessionNotifyType_UserReadMessage,
			SessionMessage: &models.SessionMessage{
				SessionMessageType: constant.SessionMessage,
				ReadMessage:        readMessage,
			},
		},
		Timestamp: time.Now().Unix(),
	}
	if err := mq.Send(topic, notify.ToString()); err != nil {
		log.Logger.Error(err.Error())
		return
	}
}

// 用户撤回发送消息通知
func sendWithDrawnMessageNotify(message rsp.Message, accountId string) {
	var notify = models.WSMessageNotify{
		WSMessageNotifyType: constant.SessionMessageNotify,
		Receive:             accountId,
		WSMessage: models.WSMessage{
			WSMessageType: im.SessionNotifyType_UserWithDrawn,
			SessionMessage: &models.SessionMessage{
				SessionMessageType: constant.SessionMessage,
				Message:            message,
			},
		},
		Timestamp: time.Now().Unix(),
	}
	if err := mq.Send(topic, notify.ToString()); err != nil {
		log.Logger.Error(err.Error())
		return
	}
}
