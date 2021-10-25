package requests

import (
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/comment"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/space"
)

type ReqAddSpace struct {
	Content     string                 `json:"content"`
	Images      []string               `json:"images"`
	VisitorType space.SpaceVisitorType `json:"visitor_type"`
}

type ReqQuerySpaces struct {
	PageCommonReq
}

type ReqAddSpaceOpt struct {
	OperatorId   string               `json:"operator_id"`
	SpaceId      string               `json:"space_id"`
	OperatorType comment.OperatorType `json:"operator_type"`
}

type ReqAddSpaceComment struct {
	SpaceId     string `json:"space_id"`
	ParentId    string `json:"parent_id"`
	Comment     string `json:"comment"`
	CommentType int64  `json:"comment_type"`
}

// 空间信息审核状态更新
type ReqUpdateSpaceAudit struct {
	SpaceId string `json:"space_id"`
	Audit   *int64 `json:"audit"`
}
