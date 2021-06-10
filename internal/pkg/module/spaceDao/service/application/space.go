package application

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/space"
	"baby-fried-rice/internal/pkg/module/spaceDao/db"
	"baby-fried-rice/internal/pkg/module/spaceDao/model/tables"
	"baby-fried-rice/internal/pkg/module/spaceDao/query"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type SpaceService struct {
}

func (service *SpaceService) SpaceAddDao(ctx context.Context, req *space.ReqSpaceAddDao) (empty *emptypb.Empty, err error) {
	var s = tables.Space{
		Origin:      req.Origin,
		Content:     req.Content,
		VisitorType: req.VisitorType,
	}
	now := time.Now()
	s.CreatedAt, s.UpdatedAt = now, now
	s.ID = handle.GenerateSerialNumberByLen(10)
	if err = db.GetDB().CreateObject(&s); err != nil {
		return
	}
	empty = new(emptypb.Empty)
	return
}

func (service *SpaceService) SpaceDeleteDao(ctx context.Context, req *space.ReqSpaceDeleteDao) (empty *emptypb.Empty, err error) {
	if err = db.GetDB().GetDB().Where("id = ? and origin = ?",
		req.Id, req.Origin).Delete(&tables.Space{}).Error; err != nil {
		return
	}
	empty = new(emptypb.Empty)
	return
}

func SpaceCommentConvert(relations []tables.SpaceCommentRelation, parentId string) (comments []*space.SpaceCommentDao) {
	comments = make([]*space.SpaceCommentDao, 0)
	for _, rel := range relations {
		if rel.ParentId != parentId {
			continue
		}
		var comment = &space.SpaceCommentDao{
			Id:          rel.ID,
			SpaceId:     rel.SpaceId,
			Content:     rel.Comment,
			CommentType: rel.CommentType,
			Origin:      rel.Origin,
			CreateTime:  rel.CreatedAt.String(),
			ReplyList:   SpaceCommentConvert(relations, rel.ID),
		}
		comments = append(comments, comment)
	}
	return
}

func (service *SpaceService) SpacesQueryDao(ctx context.Context, req *space.ReqSpacesQueryDao) (resp *space.RspSpacesQueryDao, err error) {
	var spaces []tables.Space
	if spaces, err = query.SpaceQuery(req.CommonSearchReq.Page, req.CommonSearchReq.PageSize); err != nil {
		return
	}
	var querySpaces []*space.SpaceQueryDao
	for _, s := range spaces {
		var optRelations []tables.SpaceOperatorRelation
		var likes = make([]string, 0)
		if optRelations, err = query.SpaceOptQuery(s.ID); err != nil {
			return
		}
		for _, rel := range optRelations {
			likes = append(likes, rel.Origin)
		}
		var commentRelations []tables.SpaceCommentRelation
		if commentRelations, err = query.SpaceCommentQuery(s.ID); err != nil {
			return
		}
		var querySpace = &space.SpaceQueryDao{
			Id:          s.ID,
			Origin:      s.Origin,
			Content:     s.Content,
			VisitorType: s.VisitorType,
			CreateTime:  s.CreatedAt.String(),
			Other: &space.SpaceOtherDao{
				Id:        s.ID,
				Liked:     int64(len(optRelations)),
				Likes:     likes,
				Commented: int64(len(commentRelations)),
				Comments:  SpaceCommentConvert(commentRelations, ""),
			},
		}
		querySpaces = append(querySpaces, querySpace)
	}
	resp = &space.RspSpacesQueryDao{
		Spaces:   querySpaces,
		Page:     req.CommonSearchReq.Page,
		PageSize: req.CommonSearchReq.PageSize,
	}
	return
}

func (service *SpaceService) SpaceOptAddDao(ctx context.Context, req *space.ReqSpaceOptAddDao) (empty *emptypb.Empty, err error) {
	var s = tables.SpaceOperatorRelation{
		Origin:         req.Origin,
		OperatorObject: req.OperatorObject,
		OperatorType:   req.OperatorType,
		SpaceId:        req.SpaceId,
	}
	s.CreatedAt = time.Now()
	s.OperatorId = handle.GenerateSerialNumberByLen(10)
	if err = db.GetDB().CreateObject(&s); err != nil {
		return
	}
	empty = new(emptypb.Empty)
	return
}

func (service *SpaceService) SpaceOptCancelDao(ctx context.Context, req *space.ReqSpaceOptCancelDao) (empty *emptypb.Empty, err error) {
	if err = db.GetDB().GetDB().Where("operator_id = ? and space_id = ? and origin = ?",
		req.OperatorId, req.SpaceId, req.Origin).Delete(&tables.SpaceOperatorRelation{}).Error; err != nil {
		return
	}
	empty = new(emptypb.Empty)
	return
}

func (service *SpaceService) SpaceCommentAddDao(ctx context.Context, req *space.ReqSpaceCommentAddDao) (empty *emptypb.Empty, err error) {
	var s = tables.SpaceCommentRelation{
		Origin:      req.Origin,
		ParentId:    req.ParentId,
		Comment:     req.Comment,
		SpaceId:     req.SpaceId,
		CommentType: req.CommentType,
	}
	s.CreatedAt = time.Now()
	s.ID = handle.GenerateSerialNumberByLen(10)
	if err = db.GetDB().CreateObject(&s); err != nil {
		return
	}
	empty = new(emptypb.Empty)
	return
}

func (service *SpaceService) SpaceCommentDeleteDao(ctx context.Context, req *space.ReqSpaceCommentDeleteDao) (empty *emptypb.Empty, err error) {
	if err = db.GetDB().GetDB().Where("id = ? and space_id = ? and origin = ?",
		req.Id, req.SpaceId, req.Origin).Delete(&tables.SpaceCommentRelation{}).Error; err != nil {
		return
	}
	empty = new(emptypb.Empty)
	return
}
