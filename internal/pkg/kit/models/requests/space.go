package requests

import "baby-fried-rice/internal/pkg/kit/constant"

type ReqAddSpace struct {
	Content     string                    `json:"content"`
	VisitorType constant.SpaceVisitorType `json:"visitor_type"`
}

type ReqQuerySpaces struct {
	PageCommonReq
}
