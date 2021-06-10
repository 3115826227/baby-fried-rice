package tables

import (
	"baby-fried-rice/internal/pkg/kit/models/tables"
	"time"
)

// 空间动态
type Space struct {
	tables.CommonField

	// 空间动态发布人
	Origin string `gorm:"column:origin"`
	// 空间动态内容 动态内容不可编辑，存储html格式
	Content string `gorm:"column:content"`
	// 空间动态类型 1-公开 2-仅好友可见 3-仅好友可见且部分人不可见 4-指定人可见 5-私有仅自己可见
	VisitorType int32 `gorm:"column:visitor_type"`
}

func (table *Space) TableName() string {
	return "space"
}

func (table *Space) Get() interface{} {
	return *table
}

// 空间动态操作关系
type SpaceOperatorRelation struct {
	// 空间操作id(动态/动态评论)
	OperatorId string `gorm:"column:operator_id;unique_index:space_origin_operator"`
	// 操作类型 1-点赞
	OperatorType int32 `gorm:"column:operator_type;unique_index:space_origin_operator"`
	// 操作人
	Origin string `gorm:"column:origin;unique_index:space_origin_operator"`
	// 操作对象 1-空间动态 2-空间动态评论
	OperatorObject int32 `gorm:"column:operator_object;"`
	// 空间动态id
	SpaceId string `gorm:"column:space_id" json:"space_id"`
	// 操作时间
	CreatedAt time.Time `gorm:"column:create_time" json:"created_at"`
}

func (table *SpaceOperatorRelation) TableName() string {
	return "space_operator_rel"
}

func (table *SpaceOperatorRelation) Get() interface{} {
	return *table
}

// 空间动态评论关系
type SpaceCommentRelation struct {
	tables.CommonField

	// 空间动态id
	SpaceId string `gorm:"column:space_id"`
	// 操作人
	Origin string `gorm:"column:origin"`
	// 评论上一级id
	ParentId string `gorm:"parent_id"`
	// 评论内容
	Comment string `gorm:"comment"`
	// 评论类型 1-公开，2-私有
	CommentType int32 `gorm:"comment_type"`
}

func (table *SpaceCommentRelation) TableName() string {
	return "space_comment_rel"
}

func (table *SpaceCommentRelation) Get() interface{} {
	return *table
}
