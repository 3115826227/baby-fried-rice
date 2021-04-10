package application

import (
	"baby-fried-rice/internal/pkg/kit/grpc/pbservices/common"
	"baby-fried-rice/internal/pkg/kit/grpc/pbservices/space"
	"context"
)

type SpaceService struct {
}

func (service *SpaceService) SpaceAddDao(ctx context.Context, req *space.ReqSpaceAddDao) (resp *common.CommonResponse, err error) {
	return
}

func (service *SpaceService) SpaceDeleteDao(ctx context.Context, req *space.ReqSpaceDeleteDao) (resp *common.CommonResponse, err error) {
	return
}

func (service *SpaceService) SpacesQueryDao(ctx context.Context, req *space.ReqSpacesQueryDao) (resp *common.CommonListResponse, err error) {
	return
}
