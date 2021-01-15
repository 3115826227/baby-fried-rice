package service

import (
	"github.com/3115826227/baby-fried-rice/module/im/src/middleware"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/handle"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {

	engine.GET("/api/user/friend/chat", handle.ChatHandle)
	engine.Use(middleware.MiddlewareSetUserMeta())
	app := engine.Group("/api/im")

	app.GET("/chat/message/list")
	app.GET("/chat/message/unread", handle.GetUnreadMessageHandle)
	app.GET("/chat/message/new")
	app.GET("/chat/message/history", handle.GetHistoryMessageHandle)

	app.POST("/friend", handle.FriendAdd)
	app.PATCH("/friend/remark", handle.FriendRemarkUpdate)
	app.GET("/friend", handle.Friends)
	app.DELETE("/friend", handle.FriendDelete)

	app.POST("/friend/category", handle.FriendCategoryAdd)
	app.PATCH("/friend/category", handle.FriendCategoryUpdate)
	app.GET("/friend/category", handle.FriendCategory)
	app.DELETE("/friend/category", handle.FriendCategoryDelete)

	app.POST("/official/group", handle.OfficialGroupAdd)
	app.POST("/group", handle.GroupAdd)
	app.GET("/group/list", handle.GroupListGet)
}
