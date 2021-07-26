package service

import (
	"baby-fried-rice/internal/pkg/kit/middleware"
	"baby-fried-rice/internal/pkg/module/im/config"
	_ "baby-fried-rice/internal/pkg/module/im/docs"
	"baby-fried-rice/internal/pkg/module/im/service/handle"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Register(engine *gin.Engine) {
	handle.Init()
	url := ginSwagger.URL(fmt.Sprintf("http://localhost:%v/swagger/doc.json", config.GetConfig().Server.Port))
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	app := engine.Group("/api/im", middleware.SetUserMeta())

	app.POST("/session", handle.SessionAddHandle)
	app.PATCH("/session", handle.SessionUpdateHandle)
	app.GET("/session", handle.SessionQueryHandle)
	app.GET("/session/detail", handle.SessionDetailHandle)
	app.POST("/session/join", handle.SessionJoinHandle)
	app.POST("/session/invite", handle.SessionInviteHandle)
	app.POST("/session/remove", handle.SessionRemoveHandle)
	app.GET("/session/leave", handle.SessionLeaveHandle)
	app.DELETE("/session", handle.SessionDeleteHandle)

	app.GET("/session/smsDao", handle.SessionMessageQueryHandle)
	app.GET("/session/smsDao/read_status", handle.SessionMessageReadStatusUpdateHandle)
	app.DELETE("/session/smsDao/flush", handle.SessionMessageFlushHandle)

	app.GET("/session/manage", handle.UserManageQueryHandle)
	app.PATCH("/session/manage", handle.UserManageUpdateHandle)

	app.POST("/session/operator", handle.OperatorAddHandle)
	app.PATCH("/session/operator/confirm", handle.OperatorConfirmHandle)
	app.PATCH("/session/operator/read_status", handle.OperatorReadStatusUpdateHandle)
	app.GET("/session/operator", handle.OperatorQueryHandle)
	app.DELETE("/session/operator", handle.OperatorDeleteHandle)
	app.POST("/session/friend", handle.FriendAddHandle)
	app.GET("/session/friend", handle.FriendQueryHandle)
	app.PATCH("/session/friend/black_list", handle.FriendBlackListUpdateHandle)
	app.PATCH("/session/friend/remark", handle.FriendRemarkUpdateHandle)
	app.DELETE("/session/friend", handle.FriendDeleteHandle)
}
