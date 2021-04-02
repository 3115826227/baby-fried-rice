package rsp

type CommonListResp struct {
	List     []interface{} `json:"list"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
	Total    int64           `json:"total"`
}
