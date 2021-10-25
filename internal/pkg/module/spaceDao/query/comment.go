package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/comment"
	"baby-fried-rice/internal/pkg/module/spaceDao/db"
)

type CommentQueryParams struct {
	Page     int64           `json:"page"`
	PageSize int64           `json:"page_size"`
	BizId    string          `json:"biz_id"`
	BizType  comment.BizType `json:"biz_type"`
	HostId   string          `json:"host_id"`
	ParentId string          `json:"parent_id"`
	Origin   string          `json:"origin"`
}

func CommentQuery(params CommentQueryParams) (comments []tables.CommentRelation, total int64, err error) {
	var (
		offset = int((params.Page - 1) * params.PageSize)
		limit  = int(params.PageSize)
	)
	var template = db.GetDB().GetDB().Model(&tables.CommentRelation{})
	template = template.Where("biz_id = ? and biz_type = ? and parent_id = ?",
		params.BizId, params.BizType, params.ParentId)
	if params.HostId != "" {
		template = template.Where("id = ?", params.HostId)
	}
	if err = template.Count(&total).Error; err != nil {
		return
	}
	err = template.Offset(offset).Limit(limit).Order("floor desc").Find(&comments).Error
	return
}

func CommentDetailQuery(ids []string) (detailMap map[string]tables.CommentDetail, err error) {
	var details []tables.CommentDetail
	if err = db.GetDB().GetDB().Where("comment_id in (?)", ids).Find(&details).Error; err != nil {
		return
	}
	detailMap = make(map[string]tables.CommentDetail)
	for _, d := range details {
		detailMap[d.CommentID] = d
	}
	return
}

type OperatorQueryParams struct {
	Page         int64  `json:"page"`
	PageSize     int64  `json:"page_size"`
	BizId        string `json:"biz_id"`
	HostId       string `json:"host_id"`
	Origin       string `json:"origin"`
	OperatorType comment.OperatorType
}

func OperatorQuery(params OperatorQueryParams) (opts []tables.OperatorRelation, total int64, err error) {
	var (
		offset = int((params.Page - 1) * params.PageSize)
		limit  = int(params.PageSize)
	)
	var template = db.GetDB().GetDB().Model(&tables.OperatorRelation{})
	if params.BizId != "" {
		template = template.Where("biz_id = ?", params.BizId)
	}
	if params.HostId != "" {
		template = template.Where("host_id = ?", params.HostId)
	}
	if params.Origin != "" {
		template = template.Where("origin = ?", params.Origin)
	}
	if params.OperatorType != comment.OperatorType_Default {
		template = template.Where("operator_type = ?", params.OperatorType)
	}
	if err = template.Count(&total).Error; err != nil {
		return
	}
	err = template.Offset(offset).Limit(limit).Order("create_timestamp desc").Find(&opts).Error
	return
}

type OperatorLikedQueryParams struct {
	BizId   string   `json:"biz_id"`
	HostIds []string `json:"host_ids"`
	Origin  string   `json:"origin"`
}

func OperatorLikedQuery(params OperatorLikedQueryParams) (optMap map[string]tables.OperatorRelation, err error) {
	var template = db.GetDB().GetDB().Model(&tables.OperatorRelation{})
	template = template.Where("biz_id = ? and biz_type = ? and operator_type = ? and origin = ?",
		params.BizId, comment.BizType_Space, comment.OperatorType_Like, params.Origin)
	if len(params.HostIds) != 0 {
		template = template.Where("host_id in (?)", params.HostIds)
	}
	var opts []tables.OperatorRelation
	if err = template.Find(&opts).Error; err != nil {
		return
	}
	optMap = make(map[string]tables.OperatorRelation)
	for _, opt := range opts {
		optMap[opt.HostID] = opt
	}
	return
}
