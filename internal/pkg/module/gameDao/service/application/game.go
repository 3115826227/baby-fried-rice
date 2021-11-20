package application

import (
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/game"
	"context"
)

type GameService struct {
}

// 添加游戏对局
func (service *GameService) GameRecordAddDao(ctx context.Context, req *game.ReqGameRecordAddDao) (resp *game.RspGameRecordAddDao, err error) {
	return
}

// 添加游戏进程数据
func (service *GameService) GameProcessAddDao(ctx context.Context, req *game.ReqGameProcessAddDao) (resp *game.RspGameProcessAddDao, err error) {
	return
}

// 查询游戏状态数据
func (service *GameService) GameStatusQueryDao(ctx context.Context, req *game.ReqGameStatusQueryDao) (resp *game.RspGameStatusQueryDao, err error) {
	return
}

// 查询游戏的个人对局记录
func (service *GameService) GameRecordQueryDao(ctx context.Context, req *game.ReqGameRecordQueryDao) (resp *game.RspGameRecordQueryDao, err error) {
	return
}

// 查询游戏的详情记录
func (service *GameService) GameRecordDetailQueryDao(ctx context.Context, req *game.ReqGameRecordDetailQueryDao) (resp *game.RspGameRecordDetailQueryDao, err error) {
	return
}
