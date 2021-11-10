package tables

import "baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"

// 用户与管理沟通表
type Communication struct {
	CommonIntField
	// 沟通标题
	Title  string `gorm:"column:title;not null"`
	Origin string `gorm:"column:origin;not null"`
	// 是否回复
	Reply bool `gorm:"column:reply"`
	// 沟通类型
	CommunicationType user.CommunicationType `gorm:"column:communication_type"`
	// 是否删除
	Delete bool `gorm:"column:delete"`
}

func (table *Communication) TableName() string {
	return "baby_manage_communication"
}

type CommunicationDetail struct {
	CommunicationId int64  `gorm:"column:communication_id;pk"`
	Content         string `gorm:"column:content;not null"`
	Images          string
	ReplyAccountId  string
	ReplyContent    string `gorm:"column:reply_content"`
	ReplyTimestamp  int64  `gorm:"column:reply_timestamp"`
}

func (table *CommunicationDetail) TableName() string {
	return "baby_manage_communication_detail"
}
