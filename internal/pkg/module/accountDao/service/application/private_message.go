package application

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/privatemessage"
	"baby-fried-rice/internal/pkg/module/accountDao/db"
	"baby-fried-rice/internal/pkg/module/accountDao/log"
	"baby-fried-rice/internal/pkg/module/accountDao/query"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PrivateMessageService struct {
}

func (service *PrivateMessageService) PrivateMessageAddDao(ctx context.Context, req *privatemessage.ReqPrivateMessageAddDao) (resp *privatemessage.RspPrivateMessageAddDao, err error) {
	var request = requests.UserSendPrivateMessageReq{
		SendId:          req.SendId,
		ReceiveId:       req.ReceiveId,
		MessageSendType: constant.SendPrivateMessageType(req.MessageSendType),
		MessageType:     req.MessageType,
		MessageTitle:    req.Title,
		MessageContent:  req.Content,
	}
	var id string
	if id, err = db.SendPrivateMessage(request); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	resp = &privatemessage.RspPrivateMessageAddDao{Id: id}
	return
}

func (service *PrivateMessageService) PrivateMessageStatusUpdateDao(ctx context.Context, req *privatemessage.ReqPrivateMessageStatusUpdateDao) (empty *emptypb.Empty, err error) {
	if err = db.UpdatePrivateMessagesStatus(req.AccountId, req.Ids); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

func (service *PrivateMessageService) PrivateMessageDeleteDao(ctx context.Context, req *privatemessage.ReqPrivateMessageDeleteDao) (empty *emptypb.Empty, err error) {
	if err = db.DeletePrivateMessage(req.AccountId, req.Ids); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

func (service *PrivateMessageService) PrivateMessageQueryDao(ctx context.Context, req *privatemessage.ReqPrivateMessageQueryDao) (resp *privatemessage.RspPrivateMessageQueryDao, err error) {
	var pms = requests.UserPrivateMessagesReq{
		AccountId: req.AccountId,
		SendId:    req.SendId,
		PageCommonReq: requests.PageCommonReq{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	var messages []tables.UserPrivateMessage
	var total int64
	messages, total, err = query.GetUserPrivateMessages(pms)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var list = make([]*privatemessage.PrivateMessageQueryDao, 0)
	for _, msg := range messages {
		list = append(list, privateMessageModelConvertPb(msg))
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
		PrivateMessage: privateMessageModelConvertPb(msg),
		Content:        detail.Content,
	}
	return
}

func privateMessageModelConvertPb(msg tables.UserPrivateMessage) *privatemessage.PrivateMessageQueryDao {
	return &privatemessage.PrivateMessageQueryDao{
		Id:              msg.Id,
		SendId:          msg.SendId,
		ReceiveId:       msg.ReceiveId,
		MessageType:     msg.MessageType,
		MessageSendType: msg.MessageSendType,
		Title:           msg.MessageTitle,
		CreateTime:      msg.ReceiveTime.String(),
		Status:          msg.MessageStatus,
	}
}
