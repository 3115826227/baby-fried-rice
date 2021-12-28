package tables

import (
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/blog"
)

type Blogger struct {
	ID int64 `gorm:"column:id;AUTO_INCREMENT"`
	// 博主id
	Blogger string `gorm:"column:blogger;unique"`
	// 喜欢总数
	LikeTotal int64 `gorm:"column:like_total"`
	// 阅读总数
	ReadTotal int64 `gorm:"column:read_total"`
	// 粉丝总数
	FansTotal int64 `gorm:"column:fans_total"`
	// 开通时间
	Timestamp int64 `gorm:"column:timestamp"`
}

func (table *Blogger) TableName() string {
	return "baby_blog_blogger"
}

type Blog struct {
	CommonField

	// 博文标题
	Title string `gorm:"column:title"`
	// 博主
	Blogger string `gorm:"column:blogger"`
	// 博文预览内容
	PreviewContent string `gorm:"column:preview_content"`
	// 博文所属分类
	Category string `gorm:"column:category"`
	// 博文标签
	Tags string `gorm:"column:tags"`
	// 博文状态
	Status blog.BlogStatus `gorm:"column:status"`
	// 喜欢量
	LikeTotal int64 `gorm:"column:like_total"`
	// 阅读量
	ReadTotal int64 `gorm:"column:read_total"`
	// 评论数
	CommentTotal int64 `gorm:"column:comment_total"`
}

func (table *Blog) TableName() string {
	return "baby_blog"
}

// 用户浏览记录表
type BlogVisitedRelation struct {
	ID     int64  `gorm:"column:id;AUTO_INCREMENT"`
	BlogId string `gorm:"column:blog_id;unique_index:blog_account_visit"`
	// 浏览用户
	AccountId string `gorm:"column:account_id;unique_index:blog_account_visit"`
	// 浏览时间点
	VisitTimestamp int64 `gorm:"column:visit_timestamp"`
}

func (table *BlogVisitedRelation) TableName() string {
	return "baby_blog_visited_rel"
}

type BlogUserLikeRelation struct {
	ID        int64  `gorm:"column:id;AUTO_INCREMENT"`
	BlogId    string `gorm:"column:blog_id;unique_index:blog_like_account"`
	AccountId string `gorm:"column:account_id;unique_index:blog_like_account"`
}

func (table *BlogUserLikeRelation) TableName() string {
	return "baby_blog_user_like_rel"
}

type BlogCategory struct {
	ID       int64  `gorm:"column:id;AUTO_INCREMENT"`
	Blogger  string `gorm:"column:blogger;unique_index:blogger_category"`
	Category string `gorm:"column:category;unique_index:blogger_category"`
}

func (table *BlogCategory) TableName() string {
	return "baby_blog_category"
}

type BlogTag struct {
	ID      int64  `gorm:"column:id;AUTO_INCREMENT"`
	Blogger string `gorm:"column:blogger;unique_index:blogger_tag"`
	Tag     string `gorm:"column:tag;unique_index:blogger_tag"`
}

func (table *BlogTag) TableName() string {
	return "baby_blog_tag"
}

type BloggerFansRelation struct {
	ID        int64  `gorm:"column:id;AUTO_INCREMENT"`
	Blogger   string `gorm:"column:blogger;unique_index:blogger_fans"`
	Fans      string `gorm:"column:fans;unique_index:blogger_fans'"`
	Timestamp int64  `gorm:"column:timestamp"`
}

func (table *BloggerFansRelation) TableName() string {
	return "baby_blog_blogger_fans_rel"
}
