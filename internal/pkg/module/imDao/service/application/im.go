package application

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"
	"baby-fried-rice/internal/pkg/module/imDao/db"
	"baby-fried-rice/internal/pkg/module/imDao/log"
	"baby-fried-rice/internal/pkg/module/imDao/model/tables"
	"context"
	"errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"time"
)

type IMService struct {
}

func checkFriend(origin string, users ...string) (err error) {
	var filterUsers = make([]string, 0)
	for _, u := range users {
		if origin == u {
			continue
		}
		filterUsers = append(filterUsers, u)
	}
	var count int64
	if err = db.GetDB().GetDB().Model(&tables.Friend{}).Where("origin = ? and friend in (?) and black_list = 0",
		origin, filterUsers).Count(&count).Error; err != nil {
		return
	}
	if len(filterUsers) > int(count) {
		err = errors.New("joins had user isn't your friend")
	}
	return
}

func checkUserLimit(session tables.Session, joins int) (err error) {
	// 查询会话当前人数
	var count int64
	if err = db.GetDB().GetDB().Model(&tables.SessionUserRelation{}).Where("session_id = ?", session.ID).Count(&count).Error; err != nil {
		return
	}
	// 判断人数是否超过限制
	if count+int64(joins) > int64(constant.SessionLevelUserLimitMap[session.Level]) {
		err = errors.New("user join number already exceed the session's user limit")
	}
	return
}

// 创建会话
func (service *IMService) SessionAddDao(ctx context.Context, req *im.ReqSessionAddDao) (resp *im.RspSessionAddDao, err error) {
	// 校验会话用户数目是否已经超过会话限制
	if len(req.Joins) > int(constant.SessionLevelUserLimitMap[req.Level]) {
		err = errors.New("user join number already exceed the session's user limit")
		log.Logger.Error(err.Error())
		return
	}
	// 校验会话用户是否都为创建者好友
	var joins = make([]string, 0)
	for _, j := range req.Joins {
		joins = append(joins, j.AccountId)
	}
	if err = checkFriend(req.Origin, joins...); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	now := time.Now().Unix()
	var session = tables.Session{
		Name:               req.Name,
		SessionType:        req.SessionType,
		JoinPermissionType: req.JoinPermissionType,
		Origin:             req.Origin,
		Level:              req.Level,
		UserLimit:          constant.SessionLevelUserLimitMap[req.Level],
		CreateTime:         now,
		UpdateTime:         now,
	}
	if err = db.GetDB().CreateObject(&session); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var beans = make([]interface{}, 0)
	for _, user := range req.Joins {
		var rel = tables.SessionUserRelation{
			SessionID: session.ID,
			UserID:    user.AccountId,
			Remark:    user.Remark,
			JoinTime:  now,
		}
		beans = append(beans, &rel)
	}

	if err = db.GetDB().CreateMulti(beans...); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	resp = &im.RspSessionAddDao{
		SessionId: session.ID,
	}
	return
}

// 更新会话信息
func (service *IMService) SessionUpdateDao(ctx context.Context, req *im.ReqSessionUpdateDao) (empty *emptypb.Empty, err error) {
	var session = tables.Session{
		Name:               req.Name,
		SessionType:        req.SessionType,
		JoinPermissionType: req.JoinPermissionType,
		UpdateTime:         time.Now().Unix(),
	}
	if err = db.GetDB().GetDB().Where("id = ? and origin = ?", req.SessionId, req.AccountId).Updates(&session).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 修改会话备注
func (service *IMService) SessionRemarkUpdateDao(ctx context.Context, req *im.ReqSessionRemarkUpdateDao) (empty *emptypb.Empty, err error) {
	if err = db.GetDB().GetDB().Model(&tables.SessionUserRelation{}).Where("session_id = ? and user_id = ?",
		req.SessionId, req.AccountId).Update("remark", req.Remark).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 查询会话列表
func (service *IMService) SessionQueryDao(ctx context.Context, req *im.ReqSessionQueryDao) (resp *im.RspSessionQueryDao, err error) {
	var relations []tables.SessionUserRelation
	if err = db.GetDB().GetDB().Model(&tables.SessionUserRelation{}).Where("user_id = ?", req.AccountId).Find(&relations).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var sessionIds = make([]int64, 0)
	for _, rel := range relations {
		sessionIds = append(sessionIds, rel.SessionID)
	}
	var sessions []tables.Session
	if err = db.GetDB().GetDB().Where("id in (?)", sessionIds).Find(&sessions).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var sessionDaos = make([]*im.SessionQueryDao, 0)
	for _, session := range sessions {
		if err = db.GetDB().GetDB().Where("session_id = ?", session.ID).Find(&relations).Error; err != nil {
			log.Logger.Error(err.Error())
			return
		}
		sessionDao := &im.SessionQueryDao{
			SessionId:   session.ID,
			SessionType: session.SessionType,
			Name:        session.Name,
			Level:       session.Level,
		}
		if sessionDao.SessionType == im.SessionType_DoubleSession && session.Name == "" {
			for _, rel := range relations {
				if req.AccountId != rel.UserID {
					sessionDao.Name = rel.Remark
					break
				}
			}
		}
		var template = db.GetDB().GetDB().Model(&tables.MessageUserRelation{}).
			Where("session_id = ? and receive = ?", session.ID, req.AccountId)
		var latestRel = new(tables.MessageUserRelation)
		err = template.Order("send_timestamp desc").First(&latestRel).Error
		if err == nil {
			var message tables.Message
			if err = db.GetDB().GetObject(map[string]interface{}{
				"session_id": latestRel.SessionID,
				"id":         latestRel.MessageID,
			}, &message); err != nil {
				log.Logger.Error(err.Error())
				return
			}
			var rel tables.SessionUserRelation
			if err = db.GetDB().GetObject(map[string]interface{}{
				"session_id": message.SessionID,
				"user_id":    message.Send,
			}, &rel); err != nil {
				log.Logger.Error(err.Error())
				return
			}
			sessionDao.Latest = &im.SessionMessageDao{
				SessionId:     message.SessionID,
				MessageId:     message.ID,
				MessageType:   im.SessionMessageType(message.MessageType),
				Receive:       latestRel.Receive,
				Content:       message.Content,
				SendTimestamp: message.SendTimestamp,
				ReadStatus:    latestRel.Read,
				Send: &im.JoinRemarkDao{
					AccountId: message.Send,
					Remark:    rel.Remark,
				},
			}
		} else {
			if err != gorm.ErrRecordNotFound {
				log.Logger.Error(err.Error())
				return
			}
			err = nil
		}
		if err = template.Where("`read` = 0").Count(&sessionDao.Unread).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				log.Logger.Error(err.Error())
				return
			}
			err = nil
		}
		sessionDaos = append(sessionDaos, sessionDao)
	}
	resp = &im.RspSessionQueryDao{
		Sessions: sessionDaos,
	}
	return
}

// 查询会话详细信息
func (service *IMService) SessionDetailQueryDao(ctx context.Context, req *im.ReqSessionDetailQueryDao) (resp *im.RspSessionDetailQueryDao, err error) {
	var relations []tables.SessionUserRelation
	var session tables.Session
	if err = db.GetDB().GetObject(map[string]interface{}{"id": req.SessionId}, &session); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	if err = db.GetDB().GetDB().Where("session_id = ?", req.SessionId).Find(&relations).Error; err != nil {
		log.Logger.Error(err.Error())
	}
	var joins = make([]*im.JoinRemarkDao, 0)
	for _, rel := range relations {
		joins = append(joins, &im.JoinRemarkDao{
			AccountId: rel.UserID,
			Remark:    rel.Remark,
		})
	}
	var template = db.GetDB().GetDB().Model(&tables.MessageUserRelation{}).
		Where("session_id = ? and receive = ?", req.SessionId, req.AccountId)
	var unread int64
	if err = template.Count(&unread).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	resp = &im.RspSessionDetailQueryDao{
		SessionId:          session.ID,
		SessionType:        session.SessionType,
		Name:               session.Name,
		Level:              session.Level,
		Origin:             session.Origin,
		CreateTime:         session.CreateTime,
		Joins:              joins,
		JoinPermissionType: session.JoinPermissionType,
	}
	return
}

// 加入会话
func (service *IMService) SessionJoinDao(ctx context.Context, req *im.ReqSessionJoinDao) (empty *emptypb.Empty, err error) {
	var session tables.Session
	if err = db.GetDB().GetObject(map[string]interface{}{"id": req.SessionId}, &session); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	// 校验是否可以加入会话
	if err = checkUserLimit(session, 1); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	switch session.JoinPermissionType {
	case im.SessionJoinPermissionType_NoneLimit:
	case im.SessionJoinPermissionType_OriginAudit:
		// 判断是否有审核通过的情况
		if req.OperatorId != 0 {
			var opt tables.Operator
			if err = db.GetDB().GetObject(map[string]interface{}{"id": req.OperatorId}, &opt); err != nil {
				log.Logger.Error(err.Error())
				return
			}
			/*
				1、操作类型为申请加入会话，且会话id一致
				2、操作者为会话申请加入者
				3、接收者为会话创建者
				4、接收者确认结果为同意
			*/
			if opt.OptType != int64(im.OptType_JoinSession) || opt.Origin != req.AccountId || opt.SessionId != req.SessionId ||
				opt.Receive != session.Origin || opt.Confirm == 1 {
				err = errors.New("operator condition is invalid")
				log.Logger.Error(err.Error())
				return
			}
		} else {
			err = errors.New("no permission join session")
			log.Logger.Error(err.Error())
			return
		}
	default:
		err = errors.New("no permission join session")
		log.Logger.Error(err.Error())
		return
	}
	var rel = tables.SessionUserRelation{
		SessionID: req.SessionId,
		UserID:    req.AccountId,
		JoinTime:  time.Now().Unix(),
	}
	if err = db.GetDB().CreateObject(&rel); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 离开会话
func (service *IMService) SessionLeaveDao(ctx context.Context, req *im.ReqSessionLeaveDao) (empty *emptypb.Empty, err error) {
	if err = db.GetDB().GetDB().Where("session_id = ? and user_id = ?",
		req.SessionId, req.AccountId).Delete(&tables.SessionUserRelation{}).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 邀请加入会话
func (service *IMService) SessionInviteJoinDao(ctx context.Context, req *im.ReqSessionInviteJoinDao) (empty *emptypb.Empty, err error) {
	var session tables.Session
	if err = db.GetDB().GetObject(map[string]interface{}{
		"id":     req.SessionId,
		"origin": req.Origin,
	}, &session); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	if err = checkUserLimit(session, 1); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var rel = tables.SessionUserRelation{
		SessionID: session.ID,
		UserID:    req.AccountId,
		JoinTime:  time.Now().Unix(),
	}
	if err = db.GetDB().CreateObject(&rel); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 从会话中移除
func (service *IMService) SessionRemoveDao(ctx context.Context, req *im.ReqSessionRemoveDao) (empty *emptypb.Empty, err error) {
	var session tables.Session
	if err = db.GetDB().GetObject(map[string]interface{}{
		"id":     req.SessionId,
		"origin": req.Origin,
	}, &session); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	if err = db.GetDB().GetDB().Where("session_id = ? and user_id = ?",
		session.ID, req.AccountId).Delete(&tables.SessionUserRelation{}).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	return
}

// 删除会话
func (service *IMService) SessionDeleteDao(ctx context.Context, req *im.ReqSessionDeleteDao) (empty *emptypb.Empty, err error) {
	var session tables.Session
	var exist bool
	if exist, err = db.GetDB().ExistObject(map[string]interface{}{
		"id":     req.SessionId,
		"origin": req.AccountId,
	}, &session); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	if !exist {
		err = errors.New("only origin can delete session, you have no permission")
		log.Logger.Error(err.Error())
		return
	}
	if err = db.GetDB().DeleteObject(&tables.Session{ID: req.SessionId}); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	if err = db.GetDB().GetDB().Where("session_id = ?", req.SessionId).
		Delete(&tables.SessionUserRelation{}).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 会话消息添加
func (service *IMService) SessionMessageAddDao(ctx context.Context, req *im.ReqSessionMessageAddDao) (empty *emptypb.Empty, err error) {
	var relations []tables.SessionUserRelation
	if err = db.GetDB().GetDB().Where("session_id = ?", req.SessionId).Find(&relations).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var message = tables.Message{
		SessionID:     req.SessionId,
		MessageType:   int32(req.MessageType),
		Send:          req.Send,
		Content:       req.Content,
		SendTimestamp: req.SendTimestamp,
	}
	if err = db.GetDB().CreateObject(&message); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var beans = make([]interface{}, 0)
	for _, rel := range relations {
		var msgRel = tables.MessageUserRelation{
			MessageID:     message.ID,
			SessionID:     message.SessionID,
			Receive:       rel.UserID,
			SendTimestamp: message.SendTimestamp,
		}
		if rel.UserID == message.Send {
			msgRel.Read = true
		}
		beans = append(beans, &msgRel)
	}
	if err = db.GetDB().CreateMulti(beans...); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 会话消息查询
func (service *IMService) SessionMessageQueryDao(ctx context.Context, req *im.ReqSessionMessageQueryDao) (resp *im.RspSessionMessageQueryDao, err error) {
	var relations []tables.MessageUserRelation
	var (
		offset = int((req.Page - 1) * req.PageSize)
		limit  = int(req.PageSize)
	)
	if err = db.GetDB().GetDB().Model(&tables.MessageUserRelation{}).Where("session_id = ? and receive = ?",
		req.SessionId, req.AccountId).Order("send_timestamp DESC").Offset(offset).Limit(limit).Find(&relations).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var messageIds []int64
	var readMap = make(map[int64]bool)
	for _, rel := range relations {
		messageIds = append(messageIds, rel.MessageID)
		readMap[rel.MessageID] = rel.Read
	}
	var messages []tables.Message
	if err = db.GetDB().GetDB().Where("session_id = ? and id in (?)", req.SessionId, messageIds).Find(&messages).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var messageMap = make(map[int64]tables.Message)
	for _, m := range messages {
		messageMap[m.ID] = m
	}
	var messageDaos = make([]*im.SessionMessageDao, 0)
	var sendMap = make(map[string]tables.SessionUserRelation)
	for _, rel := range relations {
		var message = messageMap[rel.MessageID]
		sendMap[message.Send] = tables.SessionUserRelation{
			SessionID: message.SessionID,
			UserID:    message.Send,
		}
		var messageDao = &im.SessionMessageDao{
			SessionId:   message.SessionID,
			MessageId:   message.ID,
			MessageType: im.SessionMessageType(message.MessageType),
			Receive:     req.AccountId,
			Send: &im.JoinRemarkDao{
				AccountId: message.Send,
			},
			Content:       message.Content,
			SendTimestamp: message.SendTimestamp,
			ReadStatus:    readMap[message.ID],
		}
		messageDaos = append(messageDaos, messageDao)
	}
	for key, rel := range sendMap {
		if err = db.GetDB().GetObject(map[string]interface{}{
			"session_id": rel.SessionID,
			"user_id":    rel.UserID,
		}, &rel); err != nil {
			log.Logger.Error(err.Error())
			return
		}
		sendMap[key] = rel
	}
	for index := range messageDaos {
		messageDaos[index].Send.Remark = sendMap[messageDaos[index].Send.AccountId].Remark
	}
	resp = &im.RspSessionMessageQueryDao{
		Messages: messageDaos,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	_, err = service.SessionMessageReadStatusUpdateDao(ctx, &im.ReqSessionMessageReadStatusUpdateDao{
		AccountId:  req.AccountId,
		SessionId:  req.SessionId,
		MessageIds: messageIds,
	})
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	return
}

// 会话消息读取状态更新
func (service *IMService) SessionMessageReadStatusUpdateDao(ctx context.Context, req *im.ReqSessionMessageReadStatusUpdateDao) (empty *emptypb.Empty, err error) {
	if err = db.GetDB().GetDB().Model(&tables.MessageUserRelation{}).Where("session_id = ? and receive = ?",
		req.SessionId, req.AccountId).Update("read", true).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 会话消息删除
func (service *IMService) SessionMessageDeleteDao(ctx context.Context, req *im.ReqSessionMessageDeleteDao) (empty *emptypb.Empty, err error) {
	if err = db.GetDB().GetDB().Model(tables.MessageUserRelation{}).
		Where("session_id = ? and receive = ? and message_id in (?)",
			req.SessionId, req.AccountId, req.MessageIds).
		Delete(&tables.MessageUserRelation{}).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 会话消息清空
func (service *IMService) SessionMessageFlushDao(ctx context.Context, req *im.ReqSessionMessageFlushDao) (empty *emptypb.Empty, err error) {
	if err = db.GetDB().GetDB().Model(tables.MessageUserRelation{}).
		Where("session_id = ? and receive = ?", req.SessionId, req.AccountId).
		Delete(&tables.MessageUserRelation{}).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 添加操作
func (service *IMService) OperatorAddDao(ctx context.Context, req *im.ReqOperatorAddDao) (resp *im.RspOperatorAddDao, err error) {
	var opt = tables.Operator{
		Origin:       req.Origin,
		Receive:      req.Receive,
		OptType:      int64(req.OptType),
		Content:      req.Content,
		NeedConfirm:  req.NeedConfirm,
		SessionId:    req.SessionId,
		OptTimestamp: time.Now().Unix(),
	}
	if err = db.GetDB().CreateObject(&opt); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	resp = &im.RspOperatorAddDao{OperatorId: opt.ID}
	return
}

// 确认操作
func (service *IMService) OperatorConfirmDao(ctx context.Context, req *im.ReqOperatorConfirmDao) (empty *emptypb.Empty, err error) {
	// 校验是否已经确认过
	var opt tables.Operator
	if err = db.GetDB().GetObject(map[string]interface{}{"id": req.OperatorId}, &opt); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	if opt.Confirm != 0 {
		err = errors.New("operator already confirm")
		log.Logger.Error(err.Error())
		return
	}
	var confirm int64
	if req.Confirm {
		// 同意
		confirm = 1
	} else {
		// 拒绝
		confirm = 2
	}
	if err = db.GetDB().GetDB().Model(&tables.Operator{}).Where("id = ? and receive = ?",
		req.OperatorId, req.AccountId).Update("confirm", confirm).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 更新操作读取状态
func (service *IMService) OperatorReadStatusUpdateDao(ctx context.Context, req *im.ReqOperatorReadStatusUpdateDao) (empty *emptypb.Empty, err error) {
	if err = db.GetDB().GetDB().Model(&tables.Operator{}).Where("id in (?) and receive = ?",
		req.OperatorIds, req.AccountId).Update("receive_read_status", true).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 查询操作列表
func (service *IMService) OperatorsQueryDao(ctx context.Context, req *im.ReqOperatorsQueryDao) (resp *im.RspOperatorsQueryDao, err error) {
	var (
		operators []tables.Operator
		offset    = int((req.Page - 1) * req.PageSize)
		limit     = int(req.PageSize)
	)
	template := db.GetDB().GetDB().Model(&tables.Operator{}).Where("origin = ? and origin_delete = 0", req.AccountId).
		Or("receive = ? and receive_delete = 0", req.AccountId)
	var count int64
	if err = template.Count(&count).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	if err = template.Order("opt_timestamp desc").Limit(limit).Offset(offset).Find(&operators).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var list = make([]*im.OperatorDao, 0)
	for _, opt := range operators {
		var operator = im.OperatorDao{
			Id:           opt.ID,
			Origin:       opt.Origin,
			Receive:      opt.Receive,
			OptType:      im.OptType(opt.OptType),
			Content:      opt.Content,
			NeedConfirm:  opt.NeedConfirm,
			Confirm:      opt.Confirm,
			SessionId:    opt.SessionId,
			OptTimestamp: opt.OptTimestamp,
		}
		if opt.Receive == req.AccountId {
			operator.ReceiveReadStatus = opt.ReceiveReadStatus
		}
		list = append(list, &operator)
	}
	resp = &im.RspOperatorsQueryDao{
		List:     list,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    count,
	}
	return
}

// 查询单个操作
func (service *IMService) OperatorSingleQueryDao(ctx context.Context, req *im.ReqOperatorSingleQueryDao) (resp *im.OperatorDao, err error) {
	var opt tables.Operator
	if err = db.GetDB().GetDB().Where("id = ? and ((origin = ? and origin_delete = 0 ) or (receive = ? and receive_delete = 0))", req.OperatorId, req.AccountId, req.AccountId).First(&opt).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	resp = &im.OperatorDao{
		Id:           opt.ID,
		Origin:       opt.Origin,
		Receive:      opt.Receive,
		OptType:      im.OptType(opt.OptType),
		Content:      opt.Content,
		NeedConfirm:  opt.NeedConfirm,
		Confirm:      opt.Confirm,
		SessionId:    opt.SessionId,
		OptTimestamp: opt.OptTimestamp,
	}
	return
}

// 删除操作
func (service *IMService) OperatorDeleteDao(ctx context.Context, req *im.ReqOperatorDeleteDao) (empty *emptypb.Empty, err error) {
	var opt tables.Operator
	if err = db.GetDB().GetObject(map[string]interface{}{"id": req.OperatorId}, &opt); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var updateColumn string
	switch req.AccountId {
	case opt.Origin:
		updateColumn = "origin_delete"
	case opt.Receive:
		updateColumn = "receive_delete"
	}
	if err = db.GetDB().GetDB().Model(&tables.Operator{}).Where("id = ?",
		req.OperatorId).Update(updateColumn, true).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 添加好友
func (service *IMService) FriendAddDao(ctx context.Context, req *im.ReqFriendAddDao) (empty *emptypb.Empty, err error) {
	// 校验好友是否有验证权限
	var um *im.RspUserManageQueryDao
	um, err = service.UserManageQueryDao(ctx, &im.ReqUserManageQueryDao{AccountId: req.AccountId})
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	switch um.AddFriendPermissionType {
	case im.AddFriendPermissionType_NoLimit:
	case im.AddFriendPermissionType_Confirm:
		// 判断是否有审核通过的情况
		if req.OperatorId != 0 {
			var opt tables.Operator
			if err = db.GetDB().GetObject(map[string]interface{}{"id": req.OperatorId}, &opt); err != nil {
				log.Logger.Error(err.Error())
				return
			}
			/*
				1、操作类型为申请添加好友
				2、操作者为申请者
				3、接收者为好友
				4、接收者确认结果为同意
			*/
			if opt.OptType != int64(im.OptType_AddFriend) || opt.Origin != req.Origin ||
				opt.Receive != req.AccountId || opt.Confirm != 1 {
				err = errors.New("operator condition is invalid")
				log.Logger.Error(err.Error())
				return
			}
		} else {
			err = errors.New("no permission join session")
			log.Logger.Error(err.Error())
			return
		}
	default:
		err = constant.NeedApplyAddFriendError
		log.Logger.Error(err.Error())
		return
	}
	now := time.Now().Unix()
	var beans = make([]interface{}, 0)
	var friend = tables.Friend{
		Origin:    req.Origin,
		Friend:    req.AccountId,
		Remark:    req.Remark,
		Timestamp: now,
	}
	var oriFriend = tables.Friend{
		Origin:    req.AccountId,
		Friend:    req.Origin,
		Remark:    req.OriRemark,
		Timestamp: now,
	}
	beans = append(beans, &friend)
	beans = append(beans, oriFriend)
	if err = db.GetDB().CreateMulti(beans...); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 好友列表查询
func (service *IMService) FriendQueryDao(ctx context.Context, req *im.ReqFriendQueryDao) (resp *im.RspFriendQueryDao, err error) {
	var friends []tables.Friend
	template := db.GetDB().GetDB()
	if req.RemarkLike != "" {
		template = template.Where("remark like ?%", req.RemarkLike)
	}
	if err = template.Where("origin = ? and black_list = ?", req.Origin, req.BlackList).Find(&friends).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var list = make([]*im.FriendDao, 0)
	for _, f := range friends {
		var friend = im.FriendDao{
			AccountId: f.Friend,
			Remark:    f.Remark,
			BlackList: f.BlackList,
			Timestamp: f.Timestamp,
		}
		list = append(list, &friend)
	}
	resp = &im.RspFriendQueryDao{List: list}
	return
}

// 好友黑名单操作
func (service *IMService) FriendBlackListDao(ctx context.Context, req *im.ReqFriendBlackListDao) (empty *emptypb.Empty, err error) {
	var friend = tables.Friend{
		Origin: req.Origin,
		Friend: req.Friend,
	}
	if err = db.GetDB().GetObject(map[string]interface{}{
		"origin": req.Origin,
		"friend": req.Friend,
	}, &friend); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	friend.BlackList = req.BlackList
	friend.Timestamp = time.Now().Unix()
	if err = db.GetDB().GetDB().Where("origin = ? and friend = ?",
		req.Origin, req.Friend).Save(&friend).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 好友备注修改
func (service *IMService) FriendRemarkDao(ctx context.Context, req *im.ReqFriendRemarkDao) (empty *emptypb.Empty, err error) {
	var friend = tables.Friend{
		Origin:    req.Origin,
		Friend:    req.Friend,
		Remark:    req.Remark,
		Timestamp: time.Now().Unix(),
	}
	if err = db.GetDB().GetDB().Where("origin = ? and friend = ?",
		req.Origin, req.Friend).Save(&friend).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 删除好友
func (service *IMService) FriendDeleteDao(ctx context.Context, req *im.ReqFriendDeleteDao) (empty *emptypb.Empty, err error) {
	if err = db.GetDB().GetDB().Where("(origin = ? and friend = ?) or (origin = ? and friend = ?)",
		req.Origin, req.Friend, req.Friend, req.Origin).Delete(&tables.Friend{}).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 用户管理更新
func (service *IMService) UserManageUpdateDao(ctx context.Context, req *im.ReqUserManageUpdateDao) (empty *emptypb.Empty, err error) {
	var userManage = tables.UserManage{
		AccountId:               req.AccountId,
		AddFriendPermissionType: req.AddFriendPermissionType,
		UpdateTimestamp:         time.Now().Unix(),
	}
	if err = db.GetDB().GetDB().Where("account_id = ?", req.AccountId).Save(&userManage).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 用户管理查询
func (service *IMService) UserManageQueryDao(ctx context.Context, req *im.ReqUserManageQueryDao) (resp *im.RspUserManageQueryDao, err error) {
	var um tables.UserManage
	if err = db.GetDB().GetObject(map[string]interface{}{"account_id": req.AccountId}, &um); err != nil {
		if err == gorm.ErrRecordNotFound {
			um = tables.UserManage{
				AccountId:       req.AccountId,
				UpdateTimestamp: time.Now().Unix(),
			}
			if err = db.GetDB().CreateObject(&um); err != nil {
				log.Logger.Error(err.Error())
				return
			}
		} else {
			log.Logger.Error(err.Error())
			return
		}
	}
	resp = &im.RspUserManageQueryDao{
		AccountId:               um.AccountId,
		AddFriendPermissionType: um.AddFriendPermissionType,
		UpdateTimestamp:         um.UpdateTimestamp,
	}
	return
}
