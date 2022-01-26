package rsp

import "baby-fried-rice/internal/pkg/kit/constant"

type CommonResp struct {
	Code    constant.Code `json:"code"`
	Message string        `json:"message"`
	Data    interface{}   `json:"data"`
}
