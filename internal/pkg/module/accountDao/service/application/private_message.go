package application

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/privatemessage"
	"baby-fried-rice/internal/pkg/module/accountDao/db"
	"baby-fried-rice/internal/pkg/module/accountDao/model/tables"
	"baby-fried-rice/internal/pkg/module/accountDao/query"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PrivateMessageService struct {
}

func (service *PrivateMessageService) PrivateMessageAddDao(ctx context.Context, req *privatemessage.ReqPrivateMessageAddDao) (empty *emptypb.Empty, err error) {
	var request = requests.UserSendPrivateMessageReq{
		SendId:          req.SendId,
		ReceiveId:       req.ReceiveId,
		MessageSendType: constant.SendPrivateMessageType(req.MessageSendType),
		MessageType:     req.MessageType,
		MessageTitle:    req.Title,
		MessageContent:  req.Content,
	}
	if err = db.SendPrivateMessage(request); err != nil {
		return
	}
	empty = new(emptypb.Empty)
	return
}

func (service *PrivateMessageService) PrivateMessageStatusUpdateDao(ctx context.Context, req *privatemessage.ReqPrivateMessageStatusUpdateDao) (empty *emptypb.Empty, err error) {
	if err = db.UpdatePrivateMessagesStatus(req.AccountId, req.Ids); err != nil {
		return
	}
	empty = new(emptypb.Empty)
	return
}

func (service *PrivateMessageService) PrivateMessageQueryDao(ctx context.Context, req *privatemessage.ReqPrivateMessageQueryDao) (resp *privatemessage.RspPrivateMessageQueryDao, err error) {
	var pms = requests.UserPrivateMessagesReq{
		SendId:        req.SendId,
		PageCommonReq: requests.PageCommonReq{},
	}
	var messages []tables.UserPrivateMessage
	var total int64
	messages, total, err = query.GetUserPrivateMessages(pms)
	if err != nil {
		return
	}
	var list = make([]*privatemessage.PrivateMessageQueryDao, 0)
	for _, msg := range messages {
		var pmd = &privatemessage.PrivateMessageQueryDao{
			Id:              msg.MessageId,
			SendId:          msg.SendId,
			ReceiveId:       msg.ReceiveId,
			MessageType:     msg.MessageType,
			MessageSendType: msg.MessageSendType,
			Title:           msg.MessageTitle,
			CreateTime:      msg.ReceiveTime.String(),
			Status:          msg.MessageStatus,
		}
		list = append(list, pmd)
	}
	resp = &privatemessage.RspPrivateMessageQueryDao{
		List:     list,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	}
	return
}

func (service *PrivateMessageService) PrivateMessageDetailDao(ctx context.Context, req *privatemessage.ReqPrivateMessageDetailDao) (resp *privatemessage.RspPrivateMessageDetailDao, err error) {
	var msg tables.UserPrivateMessage
	var detail tables.UserPrivateMessageContent
	if msg, detail, err = query.GetUserPrivateMessageDetail(req.AccountId, req.Id); err != nil {
		return
	}
	resp = &privatemessage.RspPrivateMessageDetailDao{
		PrivateMessage: &privatemessage.PrivateMessageQueryDao{
			Id:              msg.MessageId,
			SendId:          msg.SendId,
			ReceiveId:       msg.ReceiveId,
			MessageType:     msg.MessageType,
			MessageSendType: msg.MessageSendType,
			Title:           msg.MessageTitle,
			CreateTime:      msg.ReceiveTime.String(),
			Status:          msg.MessageStatus,
		},
		Content: detail.Content,
	}
	return
}
