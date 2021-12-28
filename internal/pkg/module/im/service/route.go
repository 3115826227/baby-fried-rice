package service

import (
	"baby-fried-rice/internal/pkg/kit/middleware"
	_ "baby-fried-rice/internal/pkg/module/im/docs"
	"baby-fried-rice/internal/pkg/module/im/service/handle"
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	handle.Init()
	app := engine.Group("/api/im", middleware.SetUserMeta())

	app.POST("/session", handle.SessionAddHandle)
	app.PATCH("/session", handle.SessionUpdateHandle)
	app.GET("/session", handle.SessionQueryHandle)
	app.DELETE("/session/dialog", handle.SessionDialogDeleteHandle)
	app.GET("/session/dialog", handle.SessionDialogQueryHandle)
	app.GET("/session/friend", handle.SessionByFriendQueryHandle)
	app.GET("/session/detail", handle.SessionDetailHandle)
	app.POST("/session/join", handle.SessionJoinHandle)
	app.POST("/session/invite", handle.SessionInviteHandle)
	app.PATCH("/session/remark", handle.SessionRemarkUpdateHandle)
	app.POST("/session/remove", handle.SessionRemoveHandle)
	app.GET("/session/leave", handle.SessionLeaveHandle)
	app.DELETE("/session", handle.SessionDeleteHandle)

	app.POST("/session/video/invite", handle.InviteVideoHandle)
	app.POST("/session/video/join", handle.JoinVideoHandle)
	app.PATCH("/session/video/return", handle.ReturnVideoHandle)
	app.POST("/session/video/swap", handle.SwapWebRTCSdpHandle)
	app.DELETE("/session/video/hangup", handle.HangupVideoHandle)
	app.GET("/session/video/status", handle.VideoStatusHandle)

	app.POST("/session/message", handle.SessionMessageSendHandle)
	app.GET("/session/message", handle.SessionMessageQueryHandle)
	app.GET("/session/message/read_users", handle.SessionMessageReadUsersQueryHandle)
	app.GET("/session/message/read_status", handle.SessionMessageReadStatusUpdateHandle)
	app.GET("/session/message/read_status/single", handle.SessionSingleMessageReadStatusUpdateHandle)
	app.GET("/session/message/with_drawn", handle.SessionMessageWithDrawnHandle)
	app.DELETE("/session/message", handle.SessionMessageDeleteHandle)
	app.DELETE("/session/message/flush", handle.SessionMessageFlushHandle)

	app.GET("/session/manage", handle.UserManageQueryHandle)
	app.PATCH("/session/manage", handle.UserManageUpdateHandle)

	app.POST("/session/operator", handle.OperatorAddHandle)
	app.PATCH("/session/operator/confirm", handle.OperatorConfirmHandle)
	app.PATCH("/session/operator/read_status", handle.OperatorReadStatusUpdateHandle)
	app.GET("/session/operator", handle.OperatorQueryHandle)
	app.DELETE("/session/operator", handle.OperatorDeleteHandle)
	app.POST("/session/friend", handle.FriendAddHandle)
	app.GET("/session/friends", handle.FriendQueryHandle)
	app.PATCH("/session/friend/black_list", handle.FriendBlackListUpdateHandle)
	app.PATCH("/session/friend/remark", handle.FriendRemarkUpdateHandle)
	app.DELETE("/session/friend", handle.FriendDeleteHandle)
}
