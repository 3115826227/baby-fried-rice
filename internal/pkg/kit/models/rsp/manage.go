package rsp

type BaseIterativeVersion struct {
	Version         string `json:"version"`
	Content         string `json:"content"`
	UpdateTimestamp int64  `json:"update_timestamp"`
}

type ManageIterativeVersions struct {
	List     []ManageIterativeVersion `json:"list"`
	Page     int64                    `json:"page"`
	PageSize int64                    `json:"page_size"`
	Total    int64                    `json:"total"`
}

type ManageIterativeVersion struct {
	BaseIterativeVersion
	CreateTimestamp int64 `json:"create_timestamp"`
	Status          bool  `json:"status"`
}

type UserIterativeVersion struct {
	BaseIterativeVersion
}
