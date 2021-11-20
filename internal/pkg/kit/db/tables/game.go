package tables

import "baby-fried-rice/internal/pkg/kit/rpc/pbservices/game"

// 对局记录
type GameRecord struct {
	CommonIntField
	// 游戏类型
	GameType game.GameType `gorm:"column:game_type"`
	// 游戏状态数据
	GameStatusData string `gorm:"column:game_status_data;type:text;"`
	// 游戏进程数据
	GameProcessData string `gorm:"column:game_process_data"`
	// 游戏状态
	GameStatus game.GameStatus `gorm:"column:game_status"`
	// 结束时间
	FinishTimestamp int64 `gorm:"column:finish_timestamp"`
}

func (table *GameRecord) TableName() string {
	return "baby_game_record"
}

// 用户对局记录关系表
type GameRecordUserRelation struct {
	// 游戏对局记录id
	GameRecordId int64 `gorm:"column:game_record_id;unique_index:record_type_account"`
	// 游戏类型
	GameType game.GameType `gorm:"column:game_type;unique_index:record_type_account"`
	// 用户id
	AccountId string `gorm:"column:account_id;unique_index:record_type_account"`
	// 对局结果
	Result game.GameResult `gorm:"column:result"`
	// 用户对局角色
	UserRole game.UserRole `gorm:"column:user_role"`
}

func (table *GameRecordUserRelation) TableName() string {
	return "baby_game_record_user_rel"
}
