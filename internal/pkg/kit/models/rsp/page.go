package rsp

type CommonListResp struct {
	List     []interface{} `json:"list"`
	Page     int64         `json:"page"`
	PageSize int64         `json:"page_size"`
	Total    int64         `json:"total"`
}
