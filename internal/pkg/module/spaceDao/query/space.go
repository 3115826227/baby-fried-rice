package query

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/comment"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/space"
	"baby-fried-rice/internal/pkg/module/spaceDao/db"
)

type SpaceQueryParams struct {
	Page        int64
	PageSize    int64
	SpaceId     string
	Origin      string
	VisitorType space.SpaceVisitorType
}

func SpaceQuery(params SpaceQueryParams) (spaces []tables.Space, err error) {
	var template = db.GetDB().GetDB()
	if params.SpaceId != "" {
		template = template.Where("id = ?", params.SpaceId)
	}
	if params.VisitorType != space.SpaceVisitorType_Public {
		template = template.Where("visitor_type = ?", params.VisitorType)
	}
	template = template.Where("audit_status != 2").Order("create_time desc")
	if params.Page != 0 && params.PageSize != 0 {
		var (
			offset = int((params.Page - 1) * params.PageSize)
			limit  = int(params.PageSize)
		)
		template = template.Offset(offset).Limit(limit)
	}
	if err = template.Find(&spaces).Error; err != nil {
		return
	}
	return
}

func GetSpace(id string) (space tables.Space, err error) {
	err = db.GetDB().GetObject(map[string]interface{}{"id": id}, &space)
	return
}

func SpaceDetailQuery(ids []string) (details []tables.SpaceDetail, err error) {
	err = db.GetDB().GetDB().Where("space_id in (?)", ids).Find(&details).Error
	return
}

func SpaceLikedQuery(origin string, ids []string) (likeMap map[string]struct{}, err error) {
	var opts []tables.OperatorRelation
	likeMap = make(map[string]struct{})
	if err = db.GetDB().GetDB().Where("origin = ? and host_id in (?)", origin, ids).Find(&opts).Error; err != nil {
		return
	}
	for _, opt := range opts {
		if opt.OperatorType == comment.OperatorType_Like {
			likeMap[opt.HostID] = struct{}{}
		}
		if opt.OperatorType == comment.OperatorType_CancelLike {
			delete(likeMap, opt.HostID)
		}
	}
	return
}

func SpaceOptQuery(spaceId string) (relations []tables.OperatorRelation, err error) {
	if err = db.GetDB().GetDB().Where("space_id = ?", spaceId).Find(&relations).Error; err != nil {
		return
	}
	return
}

func SpaceCommentQuery(spaceId string) (relations []tables.CommentRelation, err error) {
	if err = db.GetDB().GetDB().Where("space_id = ?", spaceId).Find(&relations).Error; err != nil {
		return
	}
	return
}
