package req

import "baby-fried-rice/internal/pkg/kit/constant"

type PageCommonReq struct {
	PageSize int `json:"page_size"`
	Page     int `json:"page"`
}

func (req *PageCommonReq) Validate() {
	if req.PageSize <= 0 {
		req.PageSize = constant.DefaultPageSize
	}
	if req.Page <= 0 {
		req.Page = constant.DefaultPage
	}
}
