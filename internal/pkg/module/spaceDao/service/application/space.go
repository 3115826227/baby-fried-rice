package application

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/comment"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/space"
	"baby-fried-rice/internal/pkg/module/spaceDao/db"
	"baby-fried-rice/internal/pkg/module/spaceDao/log"
	"baby-fried-rice/internal/pkg/module/spaceDao/query"
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"strings"
	"time"
)

type SpaceService struct {
}

// 空间动态添加
func (service *SpaceService) SpaceAddDao(ctx context.Context, req *space.ReqSpaceAddDao) (resp *space.RspSpaceAddDao, err error) {
	var s = tables.Space{
		Origin:      req.Origin,
		VisitorType: req.VisitorType,
	}
	now := time.Now()
	s.CreatedAt, s.UpdatedAt = now, now
	s.ID = handle.GenerateSerialNumberByLen(10)
	var tx = db.GetDB().GetDB().Begin()
	var images string
	if len(req.Images) != 0 {
		images = strings.Join(req.Images, ",")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	if err = tx.Create(&s).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var detail = tables.SpaceDetail{
		ID:      s.ID,
		Content: req.Content,
		Images:  images,
	}
	if err = tx.Create(&detail).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	resp = &space.RspSpaceAddDao{
		Id: s.ID,
	}
	return
}

// 空间动态删除
func (service *SpaceService) SpaceDeleteDao(ctx context.Context, req *space.ReqSpaceDeleteDao) (empty *emptypb.Empty, err error) {
	var tx = db.GetDB().GetDB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	if err = tx.Where("id = ? and origin = ?",
		req.Id, req.Origin).Delete(&tables.Space{}).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	if err = tx.Where("space_id = ?", req.Id).Delete(&tables.SpaceDetail{}).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 空间动态列表查询
func (service *SpaceService) SpaceQueryDao(ctx context.Context, req *space.ReqSpaceQueryDao) (resp *space.RspSpaceQueryDao, err error) {
	var spaces []tables.Space
	var queryReq = query.SpaceQueryParams{
		SpaceId: req.SpaceId,
		Origin:  req.Origin,
	}
	if req.CommonSearchReq != nil {
		queryReq.Page = req.CommonSearchReq.Page
		queryReq.PageSize = req.CommonSearchReq.PageSize
	}
	if spaces, err = query.SpaceQuery(queryReq); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var ids []string
	for _, s := range spaces {
		ids = append(ids, s.ID)
	}
	var details []tables.SpaceDetail
	if details, err = query.SpaceDetailQuery(ids); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var detailMap = make(map[string]tables.SpaceDetail)
	for _, d := range details {
		detailMap[d.ID] = d
	}
	var likeMap = make(map[string]struct{})
	if likeMap, err = query.SpaceLikedQuery(req.Origin, ids); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var querySpaces []*space.SpaceQueryDao
	for _, s := range spaces {
		var querySpace = &space.SpaceQueryDao{
			Id:           s.ID,
			Origin:       s.Origin,
			Content:      detailMap[s.ID].Content,
			VisitorType:  s.VisitorType,
			VisitTotal:   s.VisitTotal,
			LikeTotal:    s.LikeTotal,
			CommentTotal: s.CommentTotal,
			FloorTotal:   s.FloorTotal,
			CreateTime:   s.CreatedAt.Unix(),
		}
		if detailMap[s.ID].Images != "" {
			querySpace.Images = strings.Split(detailMap[s.ID].Images, ",")
		} else {
			querySpace.Images = make([]string, 0)
		}
		if _, exist := likeMap[s.ID]; exist {
			querySpace.OriginLiked = true
		}
		querySpaces = append(querySpaces, querySpace)
	}
	resp = &space.RspSpaceQueryDao{
		Spaces: querySpaces,
	}
	if req.CommonSearchReq != nil {
		resp.Page = req.CommonSearchReq.Page
		resp.PageSize = req.CommonSearchReq.PageSize
	}
	return
}

// 动态空间增量更新
func (service *SpaceService) SpaceIncrUpdateDao(ctx context.Context, req *space.ReqSpaceIncrUpdateDao) (empty *emptypb.Empty, err error) {
	var s tables.Space
	if s, err = query.GetSpace(req.Id); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	s.VisitTotal += req.VisitIncrement
	s.LikeTotal += req.LikeIncrement
	s.CommentTotal += req.CommentIncrement
	s.FloorTotal += req.FloorIncrement
	if err = db.GetDB().UpdateObject(&s); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 浏览记录添加
func (service *SpaceService) VisitAddDao(ctx context.Context, req *comment.ReqVisitAddDao) (resp *comment.RspVisitAddDao, err error) {
	var exist bool
	if exist, err = query.VisitedExist(req.BizId, req.BizType, req.AccountId); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	if exist {
		resp = &comment.RspVisitAddDao{Result: false}
		return
	}
	var rel = tables.VisitedRelation{
		BizID:          req.BizId,
		BizType:        req.BizType,
		AccountId:      req.AccountId,
		VisitTimestamp: time.Now().Unix(),
	}
	if err = db.GetDB().CreateObject(&rel); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	resp = &comment.RspVisitAddDao{Result: true}
	return
}

// 评论添加
func (service *SpaceService) CommentAddDao(ctx context.Context, req *comment.ReqCommentAddDao) (resp *comment.RspCommentAddDao, err error) {
	var s = tables.CommentRelation{
		BizID:    req.BizId,
		BizType:  req.BizType,
		ParentId: req.ParentId,
		Origin:   req.Origin,
		Floor:    req.Floor,
	}
	if req.ParentId != "" {
		var parentComment tables.CommentRelation
		if err = db.GetDB().GetObject(map[string]interface{}{
			"biz_id":   req.BizId,
			"biz_type": req.BizType,
			"id":       req.ParentId,
		}, &parentComment); err != nil {
			log.Logger.Error(err.Error())
			return
		}
		s.Floor = parentComment.Floor
	}
	s.CreatedAt = time.Now()
	s.UpdatedAt = s.CreatedAt
	s.ID = handle.GenerateSerialNumberByLen(10)
	var tx = db.GetDB().GetDB()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	if err = db.GetDB().CreateObject(&s); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var detail = tables.CommentDetail{
		ID:      s.ID,
		Content: req.Content,
	}
	if err = db.GetDB().CreateObject(&detail); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	resp = &comment.RspCommentAddDao{Id: s.ID}
	return
}

func CommentReplyQuery(params query.CommentQueryParams) (replies []*comment.CommentReplyDao, total int64, err error) {
	var comments []tables.CommentRelation
	if comments, total, err = query.ReplyQuery(params); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var ids = make([]string, 0)
	for _, c := range comments {
		ids = append(ids, c.ID)
	}
	var detailMap map[string]tables.CommentDetail
	if detailMap, err = query.CommentDetailQuery(ids); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var optMap map[string]tables.OperatorRelation
	var optParams = query.OperatorLikedQueryParams{
		BizId:   params.BizId,
		HostIds: ids,
		Origin:  params.Origin,
	}
	if optMap, err = query.OperatorLikedQuery(optParams); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	for _, c := range comments {
		var reply = &comment.CommentReplyDao{
			Id:              c.ID,
			ParentId:        c.ParentId,
			Content:         detailMap[c.ID].Content,
			Origin:          c.Origin,
			CreateTimestamp: c.CreatedAt.Unix(),
			LikeTotal:       c.LikeTotal,
		}
		if _, exist := optMap[c.ID]; exist {
			reply.OriginLiked = true
		}
		replies = append(replies, reply)
	}
	return
}

func CommentReplyRecursionQuery(params query.CommentQueryParams) (replies []*comment.CommentReplyDao, total int64, err error) {
	var comments []tables.CommentRelation
	var commentTotal int64
	if comments, commentTotal, err = query.CommentQuery(params); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	if commentTotal == 0 {
		return
	}
	total += commentTotal
	var ids = make([]string, 0)
	for _, c := range comments {
		ids = append(ids, c.ID)
	}
	var detailMap map[string]tables.CommentDetail
	if detailMap, err = query.CommentDetailQuery(ids); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var optMap map[string]tables.OperatorRelation
	var optParams = query.OperatorLikedQueryParams{
		BizId:   params.BizId,
		HostIds: ids,
		Origin:  params.Origin,
	}
	if optMap, err = query.OperatorLikedQuery(optParams); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	for _, c := range comments {
		params.HostId = ""
		params.ParentId = c.ID
		var childReplies []*comment.CommentReplyDao
		var replyTotal int64
		if childReplies, replyTotal, err = CommentReplyRecursionQuery(params); err != nil {
			log.Logger.Error(err.Error())
			return
		}
		var reply = &comment.CommentReplyDao{
			Id:              c.ID,
			ParentId:        c.ParentId,
			Content:         detailMap[c.ID].Content,
			Origin:          c.Origin,
			CreateTimestamp: c.CreatedAt.Unix(),
			LikeTotal:       c.LikeTotal,
			Reply:           childReplies,
		}
		if _, exist := optMap[c.ID]; exist {
			reply.OriginLiked = true
		}
		total += replyTotal
		replies = append(replies, reply)
	}
	return
}

// 评论查询
func (service *SpaceService) CommentQueryDao(ctx context.Context, req *comment.ReqCommentQueryDao) (resp *comment.RspCommentQueryDao, err error) {
	// 查询一级评论列表
	var params = query.CommentQueryParams{
		Page:     req.Page,
		PageSize: req.PageSize,
		BizId:    req.BizId,
		BizType:  req.BizType,
		Origin:   req.Origin,
	}
	var comments []tables.CommentRelation
	var total int64
	if comments, total, err = query.CommentQuery(params); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var ids = make([]string, 0)
	var replyMap = make(map[string][]*comment.CommentReplyDao)
	var replyTotalMap = make(map[string]int64)
	for _, c := range comments {
		ids = append(ids, c.ID)
		// 查询回复列表
		var replyParams = query.CommentQueryParams{
			Page:     1,
			PageSize: 4,
			BizId:    req.BizId,
			BizType:  req.BizType,
			Origin:   req.Origin,
			Floor:    c.Floor,
		}
		var replies []*comment.CommentReplyDao
		var replyTotal int64
		if replies, replyTotal, err = CommentReplyQuery(replyParams); err != nil {
			log.Logger.Error(err.Error())
			return
		}
		replyMap[c.ID] = replies
		replyTotalMap[c.ID] = replyTotal
	}
	var detailMap map[string]tables.CommentDetail
	if detailMap, err = query.CommentDetailQuery(ids); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var optMap map[string]tables.OperatorRelation
	var optParams = query.OperatorLikedQueryParams{
		BizId:   params.BizId,
		HostIds: ids,
		Origin:  params.Origin,
	}
	if optMap, err = query.OperatorLikedQuery(optParams); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var list = make([]*comment.CommentDao, 0)
	for _, c := range comments {
		var cmmt = &comment.CommentDao{
			BizId:           c.BizID,
			BizType:         c.BizType,
			Id:              c.ID,
			Content:         detailMap[c.ID].Content,
			Origin:          c.Origin,
			Floor:           c.Floor,
			CreateTimestamp: c.CreatedAt.Unix(),
			ReplyTotal:      replyTotalMap[c.ID],
			LikeTotal:       c.LikeTotal,
			Reply:           replyMap[c.ID],
		}
		if _, exist := optMap[c.ID]; exist {
			cmmt.OriginLiked = true
		}
		list = append(list, cmmt)
	}
	resp = &comment.RspCommentQueryDao{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	return
}

// 评论回复查询
func (service *SpaceService) CommentReplyQueryDao(ctx context.Context, req *comment.ReqCommentReplyQueryDao) (resp *comment.RspCommentReplyQueryDao, err error) {
	var replyParams = query.CommentQueryParams{
		Page:     req.Page,
		PageSize: req.PageSize,
		BizId:    req.BizId,
		BizType:  req.BizType,
		Origin:   req.Origin,
		ParentId: req.ParentId,
		Floor:    req.Floor,
	}
	var replies []*comment.CommentReplyDao
	var total int64
	if req.Recursion {
		if replies, total, err = CommentReplyRecursionQuery(replyParams); err != nil {
			log.Logger.Error(err.Error())
			return
		}
	} else {
		if replies, total, err = CommentReplyRecursionQuery(replyParams); err != nil {
			log.Logger.Error(err.Error())
			return
		}
	}
	resp = &comment.RspCommentReplyQueryDao{
		List:     replies,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	return
}

// 个人评论查询
func (service *SpaceService) CommentPersonQueryDao(ctx context.Context, req *comment.ReqCommentPersonQueryDao) (resp *comment.RspCommentPersonQueryDao, err error) {
	return
}

func commentDeleteDao(tx *gorm.DB, req *comment.ReqCommentDeleteDao) (total int64, err error) {
	var template = tx.Where("biz_id = ? and biz_type = ? and parent_id = ?",
		req.BizId, req.BizType, req.Id)
	if req.Origin != "" {
		template = template.Where("origin = ?", req.Origin)
	}
	var relations []tables.CommentRelation
	if err = template.Find(&relations).Error; err != nil {
		return
	}
	var ids = make([]string, 0)
	for _, rel := range relations {
		ids = append(ids, rel.ID)
	}
	if err = template.Delete(&tables.CommentRelation{}).Error; err != nil {
		return
	}
	if err = tx.Where("comment_id in (?)", ids).Delete(&tables.CommentDetail{}).Error; err != nil {
		return
	}
	total = int64(len(ids))
	for _, id := range ids {
		var childReq = &comment.ReqCommentDeleteDao{
			BizId:   req.BizId,
			BizType: req.BizType,
			Id:      id,
		}
		var count int64
		count, err = commentDeleteDao(tx, childReq)
		if err != nil {
			return
		}
		total += count
	}
	return
}

// 评论删除
func (service *SpaceService) CommentDeleteDao(ctx context.Context, req *comment.ReqCommentDeleteDao) (resp *comment.RspCommentDeleteDao, err error) {
	var tx = db.GetDB().GetDB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	var total int64
	total, err = commentDeleteDao(tx, req)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	resp = &comment.RspCommentDeleteDao{Total: total}
	return
}

// 业务下的评论清空
func (service *SpaceService) CommentClearDao(ctx context.Context, req *comment.ReqCommentClearDao) (empty *emptypb.Empty, err error) {
	return
}

// 添加操作，默认不会重复，重复判断逻辑放到上层做
func (service *SpaceService) OperatorAddDao(ctx context.Context, req *comment.ReqOperatorAddDao) (resp *comment.RspOperatorAddDao, err error) {
	// 查询点赞记录
	var likeExist bool
	likeExist, err = db.GetDB().ExistObject(map[string]interface{}{
		"biz_id":        req.BizId,
		"biz_type":      req.BizType,
		"host_id":       req.HostId,
		"operator_type": comment.OperatorType_Like,
		"origin":        req.Origin,
	}, &tables.OperatorRelation{})
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	// 添加操作记录
	var opt = tables.OperatorRelation{
		BizID:            req.BizId,
		BizType:          req.BizType,
		HostID:           req.HostId,
		OperatorType:     req.OptType,
		Origin:           req.Origin,
		CreatedTimestamp: time.Now().Unix(),
	}
	var tx = db.GetDB().GetDB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	switch req.OptType {
	case comment.OperatorType_Like:
		// 点赞操作
		if likeExist {
			// 重复点赞，返回失败
			err = fmt.Errorf("alreay liked")
			log.Logger.Error(err.Error())
			return
		} else {
			if err = tx.Create(&opt).Error; err != nil {
				log.Logger.Error(err.Error())
				return
			}
		}
	case comment.OperatorType_CancelLike:
		// 取消点赞操作
		if likeExist {
			if err = tx.Where("biz_id = ? and biz_type = ? and host_id = ? and origin = ? and operator_type = ?",
				req.BizId, req.BizType, req.HostId, req.Origin, comment.OperatorType_Like).
				Delete(&tables.OperatorRelation{}).Error; err != nil {
				log.Logger.Error(err.Error())
				return
			}
		} else {
			// 重复点赞，返回失败
			err = fmt.Errorf("alreay cancel like")
			log.Logger.Error(err.Error())
			return
		}
	}
	if req.HostId != req.BizId {
		// 处理宿主和业务id不相同的情况，即评论
		var c tables.CommentRelation
		if err = db.GetDB().GetObject(map[string]interface{}{
			"id":       req.HostId,
			"biz_id":   req.BizId,
			"biz_type": req.BizType,
		}, &c); err != nil {
			log.Logger.Error(err.Error())
			return
		}
		switch req.OptType {
		case comment.OperatorType_Like:
			// 点赞
			c.LikeTotal += 1
		case comment.OperatorType_CancelLike:
			// 取消点赞
			c.LikeTotal -= 1
		}
		if err = tx.Save(&c).Error; err != nil {
			log.Logger.Error(err.Error())
			return
		}
	}
	resp = &comment.RspOperatorAddDao{Result: true}
	return
}

// 查询操作
func (service *SpaceService) OperatorQueryDao(ctx context.Context, req *comment.ReqOperatorQueryDao) (resp *comment.RspOperatorQueryDao, err error) {
	var params = query.OperatorQueryParams{
		Page:         req.Page,
		PageSize:     req.PageSize,
		BizId:        req.Params.BizId,
		HostId:       req.Params.HostId,
		OperatorType: req.Params.OptType,
		Origin:       req.Params.Origin,
	}
	var opts []tables.OperatorRelation
	var total int64
	if opts, total, err = query.OperatorQuery(params); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var list = make([]*comment.OperatorDao, 0)
	for _, opt := range opts {
		list = append(list, &comment.OperatorDao{
			BizId:   opt.BizID,
			BizType: opt.BizType,
			HostId:  opt.HostID,
			Origin:  opt.Origin,
			OptType: opt.OperatorType,
		})
	}
	resp = &comment.RspOperatorQueryDao{
		List:     list,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	}
	return
}
