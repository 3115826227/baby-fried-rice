package requests

// 迭代版本添加
type ReqAddIterativeVersion struct {
	Version string `json:"version" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// 迭代版本更新
type ReqUpdateIterativeVersion struct {
	Version string  `json:"version" binding:"required"`
	Content *string `json:"content"`
	Status  *bool   `json:"status"`
}
