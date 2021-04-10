package tables

import "baby-fried-rice/internal/pkg/kit/models/tables"

// 空间动态
type Space struct {
	tables.CommonField

	// 空间动态发布人
	Origin string `gorm:"column:origin"`
	// 空间动态内容 动态内容不可编辑，存储html格式
	Content string `gorm:"column:content"`
	// 空间动态类型 1-公开 2-仅好友可见 3-仅好友可见且部分人不可见 4-指定人可见 5-私有仅自己可见
	VisitorType int `gorm:"column:visitor_type"`
}

func (table *Space) TableName() string {
	return "space"
}

func (table *Space) Get() interface{} {
	return *table
}

// 空间动态操作关系
type SpaceOperatorRelation struct {
	// 操作对象 1-空间动态 2-空间动态评论
	OperatorType int `gorm:"column:operator_type;"`
	// 空间动态id
	OperatorId string `gorm:"column:operator_id;unique_index:space_account_operator"`
	// 操作类型 1-点赞
	Operator int `gorm:"column:operator;unique_index:space_account_operator"`
	// 操作人
	AccountId string `gorm:"column:account_id;unique_index:space_account_operator"`
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
	AccountId string `gorm:"column:account_id"`
	// 评论上一级id
	ParentId string `gorm:"parent_id"`
	// 评论内容
	Comment string `gorm:"comment"`
	// 评论类型 1-公开，2-私有
	CommentType int `gorm:"comment_type"`
}

func (table *SpaceCommentRelation) TableName() string {
	return "space_comment_rel"
}

func (table *SpaceCommentRelation) Get() interface{} {
	return *table
}
