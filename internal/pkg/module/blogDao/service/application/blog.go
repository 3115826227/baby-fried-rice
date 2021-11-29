package application

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/blog"
	"baby-fried-rice/internal/pkg/module/blogDao/db"
	"baby-fried-rice/internal/pkg/module/blogDao/log"
	"baby-fried-rice/internal/pkg/module/blogDao/query"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

type BlogService struct {
}

// 添加标签
func (service *BlogService) TagAddDao(ctx context.Context, req *blog.ReqTagAddDao) (empty *emptypb.Empty, err error) {
	var tags = make([]tables.BlogTag, 0)
	for _, tag := range req.Tags {
		tags = append(tags, tables.BlogTag{
			Blogger: req.Origin,
			Tag:     tag,
		})
	}
	if err = db.GetDB().GetDB().Create(&tags).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 查询标签
func (service *BlogService) TagQueryDao(ctx context.Context, req *blog.ReqTagQueryDao) (resp *blog.RspTagQueryDao, err error) {
	var tags []tables.BlogTag
	if tags, err = query.Tags(req.Origin); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var list = make([]string, 0)
	for _, tag := range tags {
		list = append(list, tag.Tag)
	}
	resp = &blog.RspTagQueryDao{
		Origin: req.Origin,
		Tags:   list,
	}
	return
}

// 删除标签
func (service *BlogService) TagDeleteDao(ctx context.Context, req *blog.ReqTagDeleteDao) (empty *emptypb.Empty, err error) {
	if err = db.GetDB().GetDB().Where("blogger = ? and tag in (?)",
		req.Origin, req.Tags).Delete(&tables.BlogTag{}).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 添加分类
func (service *BlogService) CategoryAddDao(ctx context.Context, req *blog.ReqCategoryAddDao) (empty *emptypb.Empty, err error) {
	var category = tables.BlogCategory{
		Blogger:  req.Origin,
		Category: req.Category,
	}
	if err = db.GetDB().GetDB().Create(&category).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 查询分类
func (service *BlogService) CategoryQueryDao(ctx context.Context, req *blog.ReqCategoryQueryDao) (resp *blog.RspCategoryQueryDao, err error) {
	var categories []tables.BlogCategory
	if categories, err = query.Categories(req.Origin); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var list = make([]string, 0)
	for _, c := range categories {
		list = append(list, c.Category)
	}
	resp = &blog.RspCategoryQueryDao{
		Origin:   req.Origin,
		Category: list,
	}
	return
}

// 删除分类
func (service *BlogService) CategoryDeleteDao(ctx context.Context, req *blog.ReqCategoryDeleteDao) (empty *emptypb.Empty, err error) {
	if err = db.GetDB().GetDB().Where("blogger = ? and category = ?",
		req.Origin, req.Category).Delete(&tables.BlogCategory{}).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 关注/取关博主
func (service *BlogService) FocusAddDao(ctx context.Context, req *blog.ReqFocusAddDao) (empty *emptypb.Empty, err error) {
	var tx = db.GetDB().GetDB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()
	var count int64
	if err = tx.Model(&tables.BloggerFansRelation{}).Where("blogger = ? and fans = ?", req.Origin, req.Account).Count(&count).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	if req.Focus {
		var fans = tables.BloggerFansRelation{
			Blogger:   req.Origin,
			Fans:      req.Account,
			Timestamp: time.Now().Unix(),
		}
		if err = tx.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&fans).Error; err != nil {
			log.Logger.Error(err.Error())
			return
		}
		if count == 0 {
			if err = tx.Model(&tables.Blogger{}).Where("id = ? ", req.Origin).Update("fans_total", gorm.Expr("fans_total + ?", 1)).Error; err != nil {
				log.Logger.Error(err.Error())
				return
			}
		}
	} else {
		if err = tx.Where("blogger = ? and fans = ?",
			req.Origin, req.Account).Delete(&tables.BloggerFansRelation{}).Error; err != nil {
			log.Logger.Error(err.Error())
			return
		}
		if count != 0 {
			if err = tx.Model(&tables.Blogger{}).Where("id = ? ", req.Origin).Update("fans_total", gorm.Expr("fans_total + ?", -1)).Error; err != nil {
				log.Logger.Error(err.Error())
				return
			}
		}
	}
	empty = new(emptypb.Empty)
	return
}

// 喜欢/取消喜欢博文
func (service *BlogService) BlogLikeAddDao(ctx context.Context, req *blog.ReqBlogLikeAddDao) (empty *emptypb.Empty, err error) {
	var tx = db.GetDB().GetDB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()
	var b tables.Blog
	if err = tx.Where("blog_id = ?", req.BlogId).First(&b).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var count int64
	if err = tx.Model(&tables.BlogUserLikeRelation{}).Where("blog_id = ? and account_id = ?", req.BlogId, req.Account).Count(&count).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	if req.Like {
		var rel = tables.BlogUserLikeRelation{
			BlogId:    req.BlogId,
			AccountId: req.Account,
		}
		if err = tx.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&rel).Error; err != nil {
			log.Logger.Error(err.Error())
			return
		}
		if count == 0 {
			if err = tx.Model(&tables.Blog{}).Where("id = ? ", req.BlogId).Update("like_total", gorm.Expr("like_total + ?", 1)).Error; err != nil {
				log.Logger.Error(err.Error())
				return
			}
			if err = tx.Model(&tables.Blogger{}).Where("blogger = ? ", b.Blogger).Update("like_total", gorm.Expr("like_total + ?", 1)).Error; err != nil {
				log.Logger.Error(err.Error())
				return
			}
		}
	} else {
		if err = tx.Where("blog_id = ? and account_id = ?",
			req.BlogId, req.Account).Delete(&tables.BlogUserLikeRelation{}).Error; err != nil {
			log.Logger.Error(err.Error())
			return
		}
		if count != 0 {
			if err = tx.Model(&tables.Blog{}).Where("id = ? ", req.BlogId).Update("like_total", gorm.Expr("like_total + ?", -1)).Error; err != nil {
				log.Logger.Error(err.Error())
				return
			}
			if err = tx.Model(&tables.Blogger{}).Where("blogger = ? ", b.Blogger).Update("like_total", gorm.Expr("like_total + ?", -1)).Error; err != nil {
				log.Logger.Error(err.Error())
				return
			}
		}
	}
	empty = new(emptypb.Empty)
	return
}

// 查询粉丝
func (service *BlogService) FansQueryDao(ctx context.Context, req *blog.ReqFansQueryDao) (resp *blog.RspFansQueryDao, err error) {
	var fans []tables.BloggerFansRelation
	var total int64
	if fans, total, err = query.Fans(req.Blogger, req.Page, req.PageSize); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var list = make([]string, 0)
	for _, f := range fans {
		list = append(list, f.Fans)
	}
	resp = &blog.RspFansQueryDao{
		Users:    list,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	}
	return
}

// 添加博文
func (service *BlogService) BlogAddDao(ctx context.Context, req *blog.ReqBlogAddDao) (empty *emptypb.Empty, err error) {
	var b = tables.Blog{
		Title:          req.Title,
		Blogger:        req.Origin,
		PreviewContent: req.PreviewContent,
		Category:       req.Category,
		Tags:           strings.Join(req.Tags, ","),
	}
	var now = time.Now()
	b.CreatedAt, b.UpdatedAt = now, now
	if err = db.GetDB().GetDB().Create(&b).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 更新博文
func (service *BlogService) BlogUpdateDao(ctx context.Context, req *blog.ReqBlogUpdateDao) (empty *emptypb.Empty, err error) {
	var b = tables.Blog{
		Title:          req.Title,
		Blogger:        req.Origin,
		PreviewContent: req.PreviewContent,
		Category:       req.Category,
		Tags:           strings.Join(req.Tags, ","),
		Status:         req.Status,
	}
	b.UpdatedAt = time.Now()
	if err = db.GetDB().GetDB().Save(&b).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}

// 查询博主信息
func (service *BlogService) BloggerQueryDao(ctx context.Context, req *blog.ReqBloggerQueryDao) (resp *blog.RspBloggerQueryDao, err error) {
	var blogger tables.Blogger
	if blogger, err = query.Blogger(req.Blogger); err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Logger.Error(err.Error())
			return
		}
		var newBlogger = tables.Blogger{
			Blogger:   req.Blogger,
			Timestamp: time.Now().Unix(),
		}
		if err = db.GetDB().GetDB().Create(&newBlogger).Error; err != nil {
			log.Logger.Error(err.Error())
			return
		}
		blogger = newBlogger
		err = nil
	}
	var blogTotal int64
	if blogTotal, err = query.BlogTotalByBlogger(req.Blogger); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	resp = &blog.RspBloggerQueryDao{
		Blogger:   blogger.Blogger,
		LikeTotal: blogger.LikeTotal,
		BlogTotal: blogTotal,
		ReadTotal: blogger.ReadTotal,
		FansTotal: blogger.FansTotal,
	}
	return
}

// 查询博文列表
func (service *BlogService) BlogQueryDao(ctx context.Context, req *blog.ReqBlogQueryDao) (resp *blog.RspBlogQueryDao, err error) {
	var blogs []tables.Blog
	var total int64
	var param = query.BlogParams{
		QueryType: req.QueryType,
		Page:      req.Page,
		PageSize:  req.PageSize,
	}
	if blogs, total, err = query.Blog(param); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var list = make([]*blog.BlogQueryDao, 0)
	for _, b := range blogs {
		list = append(list, &blog.BlogQueryDao{
			Blogger:        b.Blogger,
			BlogId:         b.ID,
			Title:          b.Title,
			Tags:           strings.Split(b.Tags, ","),
			Category:       b.Category,
			PreviewContent: b.PreviewContent,
			LikeTotal:      b.LikeTotal,
			ReadTotal:      b.ReadTotal,
			CommentTotal:   b.CommentTotal,
			Timestamp:      b.UpdatedAt.Unix(),
		})
	}
	resp = &blog.RspBlogQueryDao{
		List:     list,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	}
	return
}

// 查询博文详情
func (service *BlogService) BlogDetailQueryDao(ctx context.Context, req *blog.ReqBlogDetailQueryDao) (resp *blog.RspBlogDetailQueryDao, err error) {
	var b tables.Blog
	if b, err = query.BlogDetail(req.BlogId); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var exist bool
	if exist, err = query.VisitedExist(req.BlogId, req.AccountId); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var tx = db.GetDB().GetDB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()
	if !exist {
		var visit = tables.BlogVisitedRelation{
			BlogId:         req.BlogId,
			AccountId:      req.AccountId,
			VisitTimestamp: time.Now().Unix(),
		}
		if err = tx.Create(&visit).Error; err != nil {
			log.Logger.Error(err.Error())
			return
		}
		if err = tx.Model(&tables.Blog{}).Where("id = ? ", req.BlogId).Update("read_total", gorm.Expr("read_total + ?", 1)).Error; err != nil {
			log.Logger.Error(err.Error())
			return
		}
		if err = tx.Model(&tables.Blogger{}).Where("id = ? ", b.Blogger).Update("read_total", gorm.Expr("read_total + ?", 1)).Error; err != nil {
			log.Logger.Error(err.Error())
			return
		}
	}
	resp = &blog.RspBlogDetailQueryDao{
		Blog: &blog.BlogQueryDao{
			Blogger:        b.Blogger,
			BlogId:         b.ID,
			Title:          b.Title,
			Tags:           strings.Split(b.Tags, ","),
			Category:       b.Category,
			PreviewContent: b.PreviewContent,
			LikeTotal:      b.LikeTotal,
			ReadTotal:      b.ReadTotal,
			CommentTotal:   b.CommentTotal,
			Timestamp:      b.UpdatedAt.Unix(),
		},
	}
	return
}

// 删除博文
func (service *BlogService) BlogDeleteDao(ctx context.Context, req *blog.ReqBlogDeleteDao) (empty *emptypb.Empty, err error) {
	if err = db.GetDB().GetDB().Where("blogger = ? and id in (?)", req.Origin, req.BlogIds).Delete(&tables.Blog{}).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	empty = new(emptypb.Empty)
	return
}
