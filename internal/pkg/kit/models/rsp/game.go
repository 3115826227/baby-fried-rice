package rsp

import "baby-fried-rice/internal/pkg/kit/rpc/pbservices/game"

type GameRecord struct {
	GameRecordId    int64           `json:"game_record_id"`
	GameStatus      game.GameStatus `json:"game_status"`
	GameType        game.GameType   `json:"game_type"`
	Result          game.GameResult `json:"result"`
	FinishTimestamp int64           `json:"finish_timestamp"`
	UserRole        game.UserRole   `json:"user_role"`
}

type PieceRoleType int32

const (
	NoPiece    PieceRoleType = 0
	RedPiece                 = 1
	BlackPiece               = 2
)

type PieceType int32

const (
	Default  PieceType = 0
	Soldiers           = 1 // 兵 卒
	Gun                = 2 // 炮
	Vehicle            = 3 // 车 車
	Horse              = 4 // 马
	Elephant           = 5 // 象 相
	Scholar            = 6 // 士 仕
	General            = 7 // 帅 将
)

type ChinaChessPoint struct {
	// 棋子角色类型 空/红方/黑方
	PieceRoleType PieceRoleType `json:"piece_role_type"`
	// 棋子类型
	PieceType PieceType `json:"piece_type"`
}

type ChinaChessBoard [10][9]ChinaChessPoint

type ChinaChessStatusResp struct {
	// 用户游戏角色
	UserRole game.UserRole `json:"user_role"`
	// 游戏状态
	GameStatus game.GameStatus `json:"game_status"`
	// 棋盘状态
	Board ChinaChessBoard `json:"board"`
}

type Point struct {
	Row int32 `json:"row"`
	Col int32 `json:"col"`
}

type ChinaChessProcess struct {
	ChinaChessPoint
	From Point `json:"from"`
	To   Point `json:"to"`
}
