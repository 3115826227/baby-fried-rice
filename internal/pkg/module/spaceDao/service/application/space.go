package application

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/space"
	"baby-fried-rice/internal/pkg/module/spaceDao/db"
	"baby-fried-rice/internal/pkg/module/spaceDao/model/tables"
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
	var s tables.Space
	s.ID = req.Id
	if err = db.GetDB().DeleteObject(&s); err != nil {
		return
	}
	empty = new(emptypb.Empty)
	return
}

func (service *SpaceService) SpacesQueryDao(ctx context.Context, req *space.ReqSpacesQueryDao) (resp *space.RspSpacesQueryDao, err error) {
	return
}
