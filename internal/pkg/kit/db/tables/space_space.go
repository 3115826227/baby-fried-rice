package tables

import (
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/comment"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/space"
)

// 空间动态
type Space struct {
	CommonField

	// 空间动态发布人
	Origin string `gorm:"column:origin"`
	// 空间访问动态类型 0-公开 1-仅好友可见 2-仅好友可见且部分人不可见 3-指定人可见 4-私有仅自己可见
	VisitorType space.SpaceVisitorType `gorm:"column:visitor_type"`
	// 浏览数
	VisitTotal int64 `gorm:"column:visit_total"`
	// 点赞数
	LikeTotal int64 `gorm:"column:like_total"`
	// 评论数
	CommentTotal int64 `gorm:"column:comment_total"`
	// 楼层数
	FloorTotal int64 `gorm:"column:floor_total"`
	// 转发数
	ForwardTotal int64 `gorm:"column:forward_total"`
	// 审核状态 0-未审核 1-审核通过 2-审核不通过
	AuditStatus int64 `gorm:"column:audit_status"`
	// 空间动态是否转发
	Forward bool `gorm:"column:forward"`
	// 转发原空间动态Id
	OriginSpaceId string `gorm:"column:origin_space_id"`
	// 是否匿名
	Anonymity bool `gorm:"column:anonymity"`
}

func (table *Space) TableName() string {
	return "space"
}

// 空间动态内容
type SpaceDetail struct {
	ID string `gorm:"column:id;pk"`
	// 文字内容
	Content string `gorm:"column:content"`
	// 图片地址列表，用','分割
	Images string `gorm:"column:images"`
}

func (table *SpaceDetail) TableName() string {
	return "space_detail"
}

// 空间动态转发记录表
type SpaceForwardRelation struct {
	ID int64 `gorm:"column:id;AUTO_INCREMENT;pk"`
	// 原空间动态
	OriginSpaceId string `gorm:"column:origin_space_id"`
	// 转发空间动态
	ForwardSpaceId string `gorm:"column:forward_space_id"`
}

func (table *SpaceForwardRelation) TableName() string {
	return "space_forward_rel"
}

// 用户浏览记录表
type VisitedRelation struct {
	ID int64 `gorm:"column:id;AUTO_INCREMENT;pk"`
	// 业务id
	BizID string `gorm:"column:biz_id;unique_index:biz_type_account"`
	// 业务类型
	BizType comment.BizType `gorm:"column:biz_type;unique_index:biz_type_account"`
	// 浏览用户
	AccountId string `gorm:"column:account_id;unique_index:biz_type_account"`
	// 浏览时间点
	VisitTimestamp int64 `gorm:"column:visit_timestamp"`
}

func (table *VisitedRelation) TableName() string {
	return "visited_rel"
}

// 操作关系
type OperatorRelation struct {
	ID int64 `gorm:"column:id;AUTO_INCREMENT" json:"id"`
	// 业务id
	BizID string `gorm:"column:biz_id;unique_index:biz_type_host_origin_operator"`
	// 业务类型
	BizType comment.BizType `gorm:"column:biz_type;unique_index:biz_type_host_origin_operator"`
	// 宿主Id
	HostID string `gorm:"host_id;unique_index:biz_type_host_origin_operator"`
	// 操作类型
	OperatorType comment.OperatorType `gorm:"column:operator_type;unique_index:biz_type_host_origin_operator"`
	// 操作人
	Origin string `gorm:"column:origin;unique_index:biz_type_host_origin_operator"`
	// 操作时间
	CreatedTimestamp int64 `gorm:"column:create_timestamp"`
}

func (table *OperatorRelation) TableName() string {
	return "operator_rel"
}

// 评论关系
type CommentRelation struct {
	CommonField

	// 业务id
	BizID string `gorm:"column:biz_id"`
	// 业务类型
	BizType comment.BizType `gorm:"column:biz_type"`
	// 评论上一级id
	ParentId string `gorm:"column:parent_id"`
	// 楼层
	Floor int64 `gorm:"column:floor"`
	// 点赞数
	LikeTotal int64 `gorm:"column:like_total"`
	// 回复数
	ReplyTotal int64 `gorm:"column:reply_total"`
	// 操作人
	Origin string `gorm:"column:origin"`
	// 是否匿名
	Anonymity bool `gorm:"column:anonymity"`
}

func (table *CommentRelation) TableName() string {
	return "comment_rel"
}

// 空间动态评论内容
type CommentDetail struct {
	ID string `gorm:"column:id;pk"`
	// 评论内容，仅支持文字评论
	Content string `gorm:"column:content"`
}

func (table *CommentDetail) TableName() string {
	return "comment_detail"
}
