package rsp

type TagResp struct {
	Tags []string `json:"tags"`
}

type CategoryResp struct {
	Categories []string `json:"categories"`
}

type Blog struct {
	BlogId         string `json:"blog_id"`
	Blogger        User   `json:"blogger"`
	Title          string `json:"title"`
	Category       string `json:"category"`
	PreviewContent string `json:"preview_content"`
	LikeTotal      int64  `json:"like_total"`
	ReadTotal      int64  `json:"read_total"`
	CommentTotal   int64  `json:"comment_total"`
	Timestamp      int64  `json:"timestamp"`
}

type BlogDetailResp struct {
	Blog
	Content string `json:"content"`
}

type BloggerResp struct {
	Blogger   User     `json:"blogger"`
	Tags      []string `json:"tags"`
	BlogTotal int64    `json:"blog_total"`
	LikeTotal int64    `json:"like_total"`
	ReadTotal int64    `json:"read_total"`
	FansTotal int64    `json:"fans_total"`
}
