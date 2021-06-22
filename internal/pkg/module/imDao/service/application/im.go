package application

import (
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

// 创建会话
func (service *IMService) SessionAddDao(ctx context.Context, req *im.ReqSessionAddDao) (empty *emptypb.Empty, err error) {
	now := time.Now().Unix()
	var session = tables.Session{
		Name:               req.Name,
		SessionType:        int32(req.SessionType),
		JoinPermissionType: int32(req.JoinPermissionType),
		Origin:             req.Origin,
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
			UserID:    user,
			JoinTime:  now,
		}
		beans = append(beans, &rel)
	}

	if err = db.GetDB().CreateMulti(beans...); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 更新会话信息
func (service *IMService) SessionUpdateDao(ctx context.Context, req *im.ReqSessionUpdateDao) (empty *emptypb.Empty, err error) {
	var session = tables.Session{
		ID:                 req.SessionId,
		Name:               req.Name,
		SessionType:        int32(req.SessionType),
		JoinPermissionType: int32(req.JoinPermissionType),
		Origin:             req.AccountId,
		UpdateTime:         time.Now().Unix(),
	}
	if err = db.GetDB().GetDB().Where("id = ? and origin = ?", req.SessionId, req.AccountId).Updates(&session).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

type SessionMessageUnread struct {
	SessionId int64 `json:"session_id"`
	Unread    int64 `json:"unread"`
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
		var joins = make([]string, 0)
		for _, rel := range relations {
			joins = append(joins, rel.UserID)
		}
		sessionDao := &im.SessionQueryDao{
			SessionId:   session.ID,
			SessionType: im.SessionType(session.SessionType),
			Origin:      session.Origin,
			Name:        session.Name,
			CreateTime:  time.Unix(session.CreateTime, 64).String(),
			Joins:       joins,
		}
		var template = db.GetDB().GetDB().Model(&tables.MessageUserRelation{}).
			Where("session_id = ? and receive = ? and `read` = 0", session.ID, req.AccountId)
		if err = template.Count(&sessionDao.Unread).Error; err != nil {
			log.Logger.Error(err.Error())
			return
		}
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
			sessionDao.Latest = &im.SessionMessageDao{
				SessionId:     message.SessionID,
				MessageId:     message.ID,
				MessageType:   im.SessionMessageType(message.MessageType),
				Send:          message.Send,
				Receive:       latestRel.Receive,
				Content:       message.Content,
				SendTimestamp: message.SendTimestamp,
				ReadStatus:    latestRel.Read,
			}
		} else {
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
	var joins = make([]string, 0)
	for _, rel := range relations {
		joins = append(joins, rel.UserID)
	}
	var template = db.GetDB().GetDB().Model(&tables.MessageUserRelation{}).
		Where("session_id = ? and receive = ?", req.SessionId, req.AccountId)
	var unread int64
	if err = template.Count(&unread).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	resp = &im.RspSessionDetailQueryDao{
		Session: &im.SessionQueryDao{
			SessionId:   session.ID,
			SessionType: im.SessionType(session.SessionType),
			Origin:      session.Origin,
			Name:        session.Name,
			Unread:      unread,
			CreateTime:  time.Unix(session.CreateTime, 64).String(),
			Joins:       joins,
		},
		JoinPermissionType: im.SessionJoinPermissionType(session.JoinPermissionType),
	}
	var latestRel = new(tables.MessageUserRelation)
	if err = template.Order("send_timestamp desc").First(latestRel).Error; err == nil {
		var message tables.Message
		if err = db.GetDB().GetObject(map[string]interface{}{"id": latestRel.SessionID}, &message); err != nil {
			log.Logger.Error(err.Error())
			return
		}
		resp.Session.Latest = &im.SessionMessageDao{
			SessionId:     message.SessionID,
			MessageId:     message.ID,
			MessageType:   im.SessionMessageType(message.MessageType),
			Send:          message.Send,
			Receive:       latestRel.Receive,
			Content:       []byte(message.Content),
			SendTimestamp: message.SendTimestamp,
			ReadStatus:    latestRel.Read,
		}
	} else if err == gorm.ErrRecordNotFound {
		log.Logger.Error(err.Error())
		err = nil
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
	switch im.SessionJoinPermissionType(session.JoinPermissionType) {
	case im.SessionJoinPermissionType_NoneLimit:
		var rel = tables.SessionUserRelation{
			SessionID: req.SessionId,
			UserID:    req.AccountId,
			JoinTime:  time.Now().Unix(),
		}
		if err = db.GetDB().CreateObject(&rel); err != nil {
			log.Logger.Error(err.Error())
			return
		}
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

// 删除会话
func (service *IMService) SessionDeleteDao(ctx context.Context, req *im.ReqSessionDeleteDao) (empty *emptypb.Empty, err error) {
	var session tables.Session
	var exist bool
	if exist, err = db.GetDB().ExistObject(map[string]interface{}{"origin": req.AccountId, "id": req.SessionId}, &session); err != nil {
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
	if err = db.GetDB().GetDB().Where("session_id = ?", req.SessionId).Delete(&tables.SessionUserRelation{}).Error; err != nil {
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
			Read:          false,
			SendTimestamp: message.SendTimestamp,
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
		req.SessionId, req.AccountId).Order("send_timestamp desc").Limit(limit).Offset(offset).Find(&relations).Error; err != nil {
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
	var messageDaos = make([]*im.SessionMessageDao, 0)
	for _, message := range messages {
		var messageDao = &im.SessionMessageDao{
			SessionId:     message.SessionID,
			MessageId:     message.ID,
			MessageType:   im.SessionMessageType(message.MessageType),
			Send:          message.Send,
			Receive:       req.AccountId,
			Content:       []byte(message.Content),
			SendTimestamp: message.SendTimestamp,
			ReadStatus:    readMap[message.ID],
		}
		messageDaos = append(messageDaos, messageDao)
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
		return
	}
	return
}

// 会话消息读取状态更新
func (service *IMService) SessionMessageReadStatusUpdateDao(ctx context.Context, req *im.ReqSessionMessageReadStatusUpdateDao) (empty *emptypb.Empty, err error) {
	if err = db.GetDB().GetDB().Model(&tables.MessageUserRelation{}).Where("session_id = ? and receive = ? and message_id in (?)",
		req.SessionId, req.AccountId, req.MessageIds).Update("read", true).Error; err != nil {
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
