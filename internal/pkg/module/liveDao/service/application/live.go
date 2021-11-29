package application

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/live"
	"baby-fried-rice/internal/pkg/module/liveDao/db"
	"baby-fried-rice/internal/pkg/module/liveDao/log"
	"baby-fried-rice/internal/pkg/module/liveDao/query"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"time"
)

type LiveService struct {
}

func (service *LiveService) LiveRoomQueryDao(ctx context.Context, req *live.ReqLiveRoomQueryDao) (resp *live.RspLiveRoomQueryDao, err error) {
	var params = query.LiveRoomParam{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	var liveRooms []tables.LiveRoom
	if liveRooms, err = query.GetLiveRooms(params); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var list = make([]*live.LiveRoom, 0)
	for _, lr := range liveRooms {
		list = append(list, &live.LiveRoom{
			LiveRoomId: lr.ID,
			Origin:     lr.Origin,
			Status:     lr.Status,
		})
	}
	resp = &live.RspLiveRoomQueryDao{
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     list,
	}
	return
}

func (service *LiveService) LiveRoomDetailQueryDao(ctx context.Context, req *live.ReqLiveRoomDetailQueryDao) (resp *live.RspLiveRoomDetailQueryDao, err error) {
	var liveRoom tables.LiveRoom
	if liveRoom, err = query.LiveRoomById(req.LiveRoomId); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var total int64
	if total, err = query.LiveRoomUserTotal(req.LiveRoomId); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	resp = &live.RspLiveRoomDetailQueryDao{
		LiveRoomId: liveRoom.ID,
		Origin:     liveRoom.Origin,
		Status:     liveRoom.Status,
		UserTotal:  total,
	}
	return
}

func (service *LiveService) LiveRoomUserQueryDao(ctx context.Context, req *live.ReqLiveRoomUserQueryDao) (resp *live.RspLiveRoomUserQueryDao, err error) {
	var users []string
	if users, err = query.LiveRoomUsers(req.LiveRoomId, req.Page, req.PageSize); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	resp = &live.RspLiveRoomUserQueryDao{
		Page:     req.Page,
		PageSize: req.PageSize,
		Users:    users,
	}
	return
}

func (service *LiveService) LiveRoomStatusUpdateDao(ctx context.Context, req *live.ReqLiveRoomStatusUpdateDao) (resp *live.RspLiveRoomDetailQueryDao, err error) {
	var liveRoom tables.LiveRoom
	var exist bool
	if liveRoom, err = query.LiveRoomByOrigin(req.Origin); err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Logger.Error(err.Error())
			return
		}
		err = nil
		exist = false
	} else {
		exist = true
	}
	var now = time.Now()
	if !exist {
		liveRoom.ID = handle.GenerateSerialNumberByLen(9)
		liveRoom.CreatedAt, liveRoom.UpdatedAt = now, now
		liveRoom.Origin = req.Origin
		liveRoom.Status = req.Status
		if err = db.GetDB().CreateObject(&liveRoom); err != nil {
			log.Logger.Error(err.Error())
			return
		}
	} else {
		liveRoom.Status = req.Status
		liveRoom.UpdatedAt = now
		if err = db.GetDB().GetDB().Where("origin = ?", req.Origin).Updates(&liveRoom).Error; err != nil {
			log.Logger.Error(err.Error())
			return
		}
	}
	var total int64
	if total, err = query.LiveRoomUserTotal(liveRoom.ID); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	resp = &live.RspLiveRoomDetailQueryDao{
		LiveRoomId: liveRoom.ID,
		Origin:     liveRoom.Origin,
		Status:     liveRoom.Status,
		UserTotal:  total,
	}
	return
}

func (service *LiveService) LiveRoomUserOptAddDao(ctx context.Context, req *live.ReqLiveRoomUserOptAddDao) (resp *live.RspLiveRoomUserOptAddDao, err error) {
	var liveRoom tables.LiveRoom
	if liveRoom, err = query.LiveRoomById(req.LiveRoomId); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	switch req.Opt {
	case live.LiveRoomUserOptType_EnterOptType:
		// 进入直播房间
		var rel = tables.LiveRoomUserRelation{
			LiveRoomID:    req.LiveRoomId,
			AccountID:     req.AccountId,
			JoinTimestamp: time.Now().Unix(),
		}
		if err = db.GetDB().CreateObject(&rel); err != nil {
			log.Logger.Error(err.Error())
			return
		}
	case live.LiveRoomUserOptType_QuitOptType:
		// 退出直播房间
		if err = db.GetDB().GetDB().Where("live_room_id = ? and account_id = ?", req.LiveRoomId, req.AccountId).Delete(&tables.LiveRoomUserRelation{}).Error; err != nil {
			log.Logger.Error(err.Error())
			return
		}
	}
	resp = &live.RspLiveRoomUserOptAddDao{
		LiveRoomId: liveRoom.ID,
		Origin:     liveRoom.Origin,
	}
	return
}

func (service *LiveService) LiveRoomMessageAddDao(ctx context.Context, req *live.ReqLiveRoomMessageAddDao) (empty *emptypb.Empty, err error) {
	var message = tables.LiveRoomMessage{
		LiveRoomID:    req.LiveRoomId,
		MessageType:   req.MessageType,
		Send:          req.Send,
		Content:       req.Content,
		SendTimestamp: req.SendTimestamp,
	}
	if err = db.GetDB().CreateObject(&message); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

func (service *LiveService) LiveRoomMessageQueryDao(ctx context.Context, req *live.ReqLiveRoomMessageQueryDao) (resp *live.RspLiveRoomMessageQueryDao, err error) {
	var rel tables.LiveRoomUserRelation
	if rel, err = query.LiveRoomUserRelation(req.AccountId, req.LiveRoomId); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var param = query.LiveRoomParamMessage{
		StartTimestamp: rel.JoinTimestamp,
		Page:           req.Page,
		PageSize:       req.PageSize,
	}
	var messages []tables.LiveRoomMessage
	if messages, err = query.GetLiveRoomMessage(param); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var list = make([]*live.LiveRoomMessage, 0)
	for _, msg := range messages {
		list = append(list, &live.LiveRoomMessage{
			MessageId:     msg.ID,
			MessageType:   msg.MessageType,
			Send:          msg.Send,
			Content:       msg.Content,
			SendTimestamp: msg.SendTimestamp,
		})
	}
	resp = &live.RspLiveRoomMessageQueryDao{
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     list,
	}
	return
}
