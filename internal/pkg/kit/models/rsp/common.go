package rsp

type CommonResp struct {
	Code    int         `json:"code"`
	Message string      `json:"smsDao"`
	Data    interface{} `json:"data"`
}
