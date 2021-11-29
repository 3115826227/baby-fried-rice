package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/blog"
	"baby-fried-rice/internal/pkg/module/blogDao/db"
)

type BlogParams struct {
	QueryType blog.QueryType
	Page      int64
	PageSize  int64
}

func Blog(param BlogParams) (blogs []tables.Blog, total int64, err error) {
	var (
		offset = int((param.Page - 1) * param.PageSize)
		limit  = int(param.PageSize)
	)
	var template = db.GetDB().GetDB().Model(&tables.Blog{})
	if err = template.Count(&total).Error; err != nil {
		return
	}
	template = template.Offset(offset).Limit(limit)
	switch param.QueryType {
	case blog.QueryType_NewType:
		template = template.Order("update_time desc")
	case blog.QueryType_LikeMostType:
		template = template.Order("like_total desc")
	case blog.QueryType_ReadMostType:
		template = template.Order("read_total desc")
	}
	err = template.Find(&blogs).Error
	return
}

func Tags(blogger string) (tags []tables.BlogTag, err error) {
	err = db.GetDB().GetDB().Where("blogger = ?", blogger).Find(&tags).Error
	return
}

func Categories(blogger string) (categories []tables.BlogCategory, err error) {
	err = db.GetDB().GetDB().Where("blogger = ?", blogger).Find(&categories).Error
	return
}

func BlogDetail(blogId string) (blog tables.Blog, err error) {
	err = db.GetDB().GetDB().Where("id = ?", blogId).Find(&blog).Error
	return
}

func Blogger(id string) (blogger tables.Blogger, err error) {
	err = db.GetDB().GetDB().Where("blogger = ?", id).First(&blogger).Error
	return
}

func BlogTotalByBlogger(blogger string) (total int64, err error) {
	err = db.GetDB().GetDB().Model(&tables.Blog{}).Where("blogger = ?", blogger).Count(&total).Error
	return
}

func Fans(blogger string, page, pageSize int64) (users []tables.BloggerFansRelation, total int64, err error) {
	var (
		offset = int((page - 1) * pageSize)
		limit  = int(pageSize)
	)
	var template = db.GetDB().GetDB().Model(&tables.BloggerFansRelation{}).Where("blogger_id = ?", blogger)
	if err = template.Count(&total).Error; err != nil {
		return
	}
	err = template.Offset(offset).Limit(limit).Find(&users).Order("timestamp").Error
	return
}

func VisitedExist(blogId string, accountId string) (exist bool, err error) {
	return db.GetDB().ExistObject(map[string]interface{}{
		"blog_id":    blogId,
		"account_id": accountId,
	}, &tables.BlogVisitedRelation{})
}
