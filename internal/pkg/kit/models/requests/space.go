package requests

import "baby-fried-rice/internal/pkg/kit/constant"

type ReqAddSpace struct {
	Content     string                    `json:"content"`
	VisitorType constant.SpaceVisitorType `json:"visitor_type"`
}

type ReqQuerySpaces struct {
	PageCommonReq
}

type ReqAddSpaceOpt struct {
	SpaceId        string `json:"space_id"`
	OperatorObject int32  `json:"operator_object"`
	OperatorType   int32  `json:"operator_type"`
}

type ReqAddSpaceComment struct {
	SpaceId     string `json:"space_id"`
	ParentId    string `json:"parent_id"`
	Comment     string `json:"comment"`
	CommentType int32  `json:"comment_type"`
}
