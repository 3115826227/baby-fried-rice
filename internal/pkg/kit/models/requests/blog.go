package requests

import "baby-fried-rice/internal/pkg/kit/rpc/pbservices/blog"

type ReqAddTag struct {
	Tag string `json:"tag"`
}

type ReqAddCategory struct {
	Category string `json:"category"`
}

type ReqAddBlog struct {
	Title    string   `json:"title"`
	Tags     []string `json:"tags"`
	Category string   `json:"category"`
	Content  string   `json:"content"`
}

type ReqUpdateBlog struct {
	BlogId   string          `json:"blog_id"`
	Title    string          `json:"title"`
	Tags     []string        `json:"tags"`
	Category string          `json:"category"`
	Content  string          `json:"content"`
	Status   blog.BlogStatus `json:"status"`
}

type ReqFocusOnBlogger struct {
	Blogger string `json:"blogger"`
	FocusOn bool   `json:"focus_on"`
}

type ReqBlogLikeBlogger struct {
	BlogId string `json:"blog_id"`
	Like   bool   `json:"like"`
}
