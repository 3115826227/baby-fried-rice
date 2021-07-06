package tables

import (
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"
)

// 信息操作管理
type Operator struct {
	// 操作id
	ID int64 `gorm:"column:id;primaryKey;AUTO_INCREMENT" json:"id"`
	// 操作者
	Origin string `gorm:"column:origin" json:"origin"`
	// 受众者
	Receive string `gorm:"column:receive" json:"receive"`
	// 操作类型
	OptType int64 `gorm:"column:opt_type" json:"opt_type"`
	// 信息内容
	Content string `gorm:"column:content" json:"content"`
	// 是否需要确认
	NeedConfirm bool `gorm:"column:need_confirm" json:"need_confirm"`
	// 确认 0-未确认 1-同意 2-拒绝
	Confirm int64 `gorm:"column:confirm" json:"confirm"`
	// 会话id
	SessionId int64 `gorm:"column:session_id" json:"session_id"`
	// 操作时间戳
	OptTimestamp int64 `gorm:"column:opt_timestamp" json:"opt_timestamp"`
	// 接收者读取状态
	ReceiveReadStatus bool `gorm:"column:receive_read_status" json:"receive_read_status"`
	// 操作者是否删除
	OriginDelete bool `gorm:"origin_delete" json:"origin_delete"`
	// 接收者是否删除
	ReceiveDelete bool `gorm:"receive_delete" json:"receive_delete"`
}

func (table *Operator) TableName() string {
	return "baby_im_operator"
}

type Friend struct {
	// 用户id
	Origin string `gorm:"column:origin" json:"origin"`
	// 好友id
	Friend string `gorm:"column:friend" json:"friend"`
	// 好友备注
	Remark string `gorm:"column:remark" json:"remark"`
	// 是否为黑名单 0-否 1-是
	BlackList bool `gorm:"column:black_list" json:"black_list"`
	// 成为好友的时间
	Timestamp int64 `gorm:"column:timestamp" json:"timestamp"`
	// 审核 0-不需要审核 1-审核中 2-审核通过
	Audit int64 `gorm:"column:audit" json:"audit"`
}

func (table *Friend) TableName() string {
	return "baby_im_friend"
}

type UserManage struct {
	// 用户id
	AccountId string `gorm:"column:account_id;primaryKey" json:"account_id"`
	// 好友加入权限
	AddFriendPermissionType im.AddFriendPermissionType `gorm:"column:add_permission_type" json:"add_friend_permission_type"`
	// 更新时间
	UpdateTimestamp int64 `gorm:"update_timestamp" json:"update_timestamp"`
}

func (table *UserManage) TableName() string {
	return "baby_im_user_manage"
}
