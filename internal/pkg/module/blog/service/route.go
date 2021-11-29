package service

import (
	"baby-fried-rice/internal/pkg/kit/middleware"
	"baby-fried-rice/internal/pkg/module/blog/service/handle"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	app := engine.Group("/api/blog", middleware.SetUserMeta())
	app.POST("/tag", handle.TagAddHandle)
	app.GET("/tag", handle.TagHandle)
	app.DELETE("/tag", handle.TagDeleteHandle)
	app.POST("/category", handle.CategoryAddHandle)
	app.GET("/category", handle.CategoryHandle)
	app.DELETE("/category", handle.CategoryDeleteHandle)
	app.POST("/blog", handle.BlogAddHandle)
	app.PATCH("/blog", handle.BlogUpdateHandle)
	app.GET("/blog", handle.BlogHandle)
	app.GET("/blog/detail", handle.BlogDetailHandle)
	app.DELETE("/blog", handle.BlogDeleteHandle)
	app.GET("/blogger", handle.BloggerHandle)
	app.POST("/blogger/focus_on", handle.BloggerFocusOnHandle)
	app.POST("/blog/like", handle.BlogLikeHandle)
	app.GET("/blog/fans", handle.BloggerFansHandle)
}
