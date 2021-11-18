package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/common"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/im/grpc"
	"baby-fried-rice/internal/pkg/module/im/log"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// SessionAddHandle 会话创建接口
// @Summary 会话创建接口
// @Description 会话创建接口
// @Tags 会话相关接口
// @Accept application/json
// @Produce application/json
// @Param accountId header string true "用户id"
// @Param username header string true "用户名"
// @Param session body ReqAddSession true "会话"
// @Success 200 {string} rsp.CommonResp
// @Router /session [post]
func SessionAddHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var err error
	var req requests.ReqAddSession
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	userClient, err := grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var resp *user.RspUserDaoById
	resp, err = userClient.UserDaoById(context.Background(), &user.ReqUserDaoById{Ids: req.Joins})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var joins = make([]*im.JoinRemarkDao, 0)
	for _, u := range resp.Users {
		var join = &im.JoinRemarkDao{
			AccountId: u.Id,
			Remark:    u.Username,
		}
		joins = append(joins, join)
	}
	var reqSession = &im.ReqSessionAddDao{
		SessionType:        req.SessionType,
		JoinPermissionType: req.JoinPermissionType,
		Name:               req.Name,
		Origin:             userMeta.AccountId,
		Joins:              joins,
		Level:              req.SessionLevel,
	}
	_, err = imClient.SessionAddDao(context.Background(), reqSession)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// SessionQueryHandle 会话列表查询接口
// @Summary 会话列表查询接口
// @Description 会话列表查询接口
// @Tags 会话相关接口
// @Accept application/json
// @Param accountId header string true "用户id"
// @Param username header string true "用户名"
// @Success 200 {string} rsp.CommonResp
// @Router /session [get]
// 会话列表查询
func SessionQueryHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	reqPage, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var sessionId int
	sessionIdStr := c.Query("session_id")
	if sessionIdStr != "" {
		sessionId, err = strconv.Atoi(sessionIdStr)
		if err != nil {
			log.Logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, handle.ParamErrResponse)
			return
		}
	}
	var imClient im.DaoImClient
	imClient, err = grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var reqSession = &im.ReqSessionQueryDao{
		AccountId: userMeta.AccountId,
		SessionId: int64(sessionId),
		CommonSearchReq: &common.CommonSearchRequest{
			Page:     reqPage.Page,
			PageSize: reqPage.PageSize,
		},
	}
	var resp *im.RspSessionQueryDao
	resp, err = imClient.SessionQueryDao(context.Background(), reqSession)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var list = make([]interface{}, 0)
	for _, s := range resp.Sessions {
		list = append(list, rsp.Session{
			SessionId:   s.SessionId,
			SessionType: s.SessionType,
			Name:        s.Name,
			Origin:      s.Origin,
		})
	}
	handle.SuccessListResp(c, "", list, resp.Total, reqPage.Page, reqPage.PageSize)
}

// 对话框会话列表查询
func SessionDialogQueryHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	reqPage, err := handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var imClient im.DaoImClient
	imClient, err = grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var reqSession = &im.ReqSessionDialogQueryDao{
		AccountId: userMeta.AccountId,
		Page:      reqPage.Page,
		PageSize:  reqPage.PageSize,
	}
	var resp *im.RspSessionDialogQueryDao
	resp, err = imClient.SessionDialogQueryDao(context.Background(), reqSession)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var sessions = make([]interface{}, 0)
	for _, s := range resp.Sessions {
		var session = rsp.SessionDialog{
			SessionId:   s.SessionId,
			SessionType: s.SessionType,
			Name:        s.Name,
			Level:       s.Level,
			Unread:      s.Unread,
		}
		if s.Latest != nil {
			lm := s.Latest
			session.LatestMessage = &rsp.Message{
				SessionId:   lm.SessionId,
				MessageId:   lm.MessageId,
				MessageType: lm.MessageType,
				Send: rsp.User{
					AccountID: lm.Send.AccountId,
					Remark:    lm.Send.Remark,
				},
				Receive:       lm.Receive,
				Content:       lm.Content,
				SendTimestamp: lm.SendTimestamp,
				ReadStatus:    lm.ReadStatus,
			}
		}
		sessions = append(sessions, session)
	}
	handle.SuccessListResp(c, "", sessions, resp.Total, reqPage.Page, reqPage.PageSize)
}

// 根据好友查询会话
func SessionByFriendQueryHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var reqSession = &im.ReqSessionByFriendQueryDao{
		AccountId: userMeta.AccountId,
		Friend:    c.Query("friend"),
	}
	var resp *im.RspSessionByFriendQueryDao
	resp, err = imClient.SessionByFriendQueryDao(context.Background(), reqSession)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", resp.SessionId)
}

// 对话框会话删除
func SessionDialogDeleteHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	sessionId, err := strconv.Atoi(c.Query("session_id"))
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var reqSession = &im.ReqSessionDialogDao{
		AccountId: userMeta.AccountId,
		SessionId: int64(sessionId),
	}
	_, err = imClient.SessionDialogDeleteDao(context.Background(), reqSession)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

func SessionDetailHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	sessionId, err := strconv.Atoi(c.Query("session_id"))
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var reqSession = &im.ReqSessionDetailQueryDao{
		AccountId: userMeta.AccountId,
		SessionId: int64(sessionId),
	}
	var resp *im.RspSessionDetailQueryDao
	resp, err = imClient.SessionDetailQueryDao(context.Background(), reqSession)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var ids = make([]string, 0)
	for _, u := range resp.Joins {
		ids = append(ids, u.AccountId)
	}
	userClient, err := grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var userResp *user.RspUserDaoById
	userResp, err = userClient.UserDaoById(context.Background(), &user.ReqUserDaoById{Ids: ids})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var userMap = make(map[string]*user.UserDao)
	for _, u := range userResp.Users {
		userMap[u.Id] = u
	}
	var joins = make([]rsp.User, 0)
	for _, u := range resp.Joins {
		var join = rsp.User{
			AccountID:  u.AccountId,
			Username:   userMap[u.AccountId].Username,
			HeadImgUrl: userMap[u.AccountId].HeadImgUrl,
			Remark:     u.Remark,
			OnlineType: u.OnlineType,
		}
		joins = append(joins, join)
	}
	var res = rsp.SessionDetail{
		SessionId:          resp.SessionId,
		SessionType:        resp.SessionType,
		Name:               resp.Name,
		Level:              resp.Level,
		Origin:             resp.Origin,
		Joins:              joins,
		JoinPermissionType: resp.JoinPermissionType,
		CreateTime:         resp.CreateTime,
	}
	handle.SuccessResp(c, "", res)
}

func SessionUpdateHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var err error
	var req requests.ReqUpdateSession
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var reqSession = &im.ReqSessionUpdateDao{
		SessionId:          req.SessionId,
		JoinPermissionType: req.JoinPermissionType,
		Name:               req.Name,
		AccountId:          userMeta.AccountId,
	}
	_, err = imClient.SessionUpdateDao(context.Background(), reqSession)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 加入会话
func SessionJoinHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	sessionId, err := strconv.Atoi(c.Query("session_id"))
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var reqSession = &im.ReqSessionDetailQueryDao{
		AccountId: userMeta.AccountId,
		SessionId: int64(sessionId),
	}
	var resp *im.RspSessionDetailQueryDao
	resp, err = imClient.SessionDetailQueryDao(context.Background(), reqSession)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	switch resp.JoinPermissionType {
	case im.SessionJoinPermissionType_NoneLimit:
		var reqJoinSession = &im.ReqSessionJoinDao{
			AccountId: userMeta.AccountId,
			SessionId: int64(sessionId),
		}
		if _, err = imClient.SessionJoinDao(context.Background(), reqJoinSession); err != nil {
			log.Logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
			return
		}
	case im.SessionJoinPermissionType_InviteJoin:
		err = constant.NeedInviteJoinSessionError
		log.Logger.Error(err.Error())
		handle.ErrorResp(c, http.StatusOK, handle.CodeNeedInviteJoinSession, handle.CodeNeedInviteJoinSessionMsg)
		return
	case im.SessionJoinPermissionType_OriginAudit:
		var reqOperator = &im.ReqOperatorAddDao{
			Origin:      userMeta.AccountId,
			Receive:     resp.Origin,
			OptType:     im.OptType_JoinSession,
			Content:     constant.JoinSessionOptReqContent,
			NeedConfirm: true,
		}
		if _, err = imClient.OperatorAddDao(context.Background(), reqOperator); err != nil {
			log.Logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
			return
		}
		log.Logger.Info(constant.CodeNeedOriginAuditSessionMsg)
		handle.ErrorResp(c, http.StatusOK, handle.CodeNeedOriginAuditSession, handle.CodeNeedOriginAuditSessionMsg)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 邀请加入会话
func SessionInviteHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var req requests.ReqInviteJoinSession
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var resp *im.RspSessionDetailQueryDao
	resp, err = imClient.SessionDetailQueryDao(context.Background(), &im.ReqSessionDetailQueryDao{
		SessionId: req.SessionId,
		AccountId: userMeta.AccountId,
	})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var session = rsp.Session{
		SessionId:   resp.SessionId,
		SessionType: resp.SessionType,
		Name:        resp.Name,
		Origin:      resp.Origin,
	}
	var reqSession = &im.ReqSessionInviteJoinDao{
		Origin:    userMeta.AccountId,
		SessionId: req.SessionId,
		AccountId: req.AccountId,
	}
	_, err = imClient.SessionInviteJoinDao(context.Background(), reqSession)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	go sendInviteNotify(session, userMeta.GetUser(), req.AccountId)
	handle.SuccessResp(c, "", nil)
}

// 从会话中移除
func SessionRemoveHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var req requests.ReqRemoveFromSession
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var reqSession = &im.ReqSessionRemoveDao{
		Origin:    userMeta.AccountId,
		SessionId: req.SessionId,
		AccountId: req.AccountId,
	}
	_, err = imClient.SessionRemoveDao(context.Background(), reqSession)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 离开会话
func SessionLeaveHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	sessionId, err := strconv.Atoi(c.Query("session_id"))
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var reqSession = &im.ReqSessionLeaveDao{
		SessionId: int64(sessionId),
		AccountId: userMeta.AccountId,
	}
	_, err = imClient.SessionLeaveDao(context.Background(), reqSession)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 会话删除
func SessionDeleteHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	sessionId, err := strconv.Atoi(c.Query("session_id"))
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var reqSession = &im.ReqSessionDeleteDao{
		AccountId: userMeta.AccountId,
		SessionId: int64(sessionId),
	}
	_, err = imClient.SessionDeleteDao(context.Background(), reqSession)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 会话消息查询
func SessionMessageQueryHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	sessionId, err := strconv.Atoi(c.Query("session_id"))
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var req requests.PageCommonReq
	req, err = handle.PageHandle(c)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var reqSession = &im.ReqSessionMessageQueryDao{
		AccountId: userMeta.AccountId,
		SessionId: int64(sessionId),
		Page:      req.Page,
		PageSize:  req.PageSize,
	}
	var resp *im.RspSessionMessageQueryDao
	resp, err = imClient.SessionMessageQueryDao(context.Background(), reqSession)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var idMap = make(map[string]struct{})
	var ids = make([]string, 0)
	for _, u := range resp.Messages {
		if _, exist := idMap[u.Send.AccountId]; !exist {
			idMap[u.Send.AccountId] = struct{}{}
			ids = append(ids, u.Send.AccountId)
		}
	}
	userClient, err := grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var userResp *user.RspUserDaoById
	userResp, err = userClient.UserDaoById(context.Background(), &user.ReqUserDaoById{Ids: ids})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var userMap = make(map[string]*user.UserDao)
	for _, u := range userResp.Users {
		userMap[u.Id] = u
	}
	var msgs = make([]rsp.Message, 0)
	for _, m := range resp.Messages {
		var msg = rsp.Message{
			SessionId:   m.SessionId,
			MessageId:   m.MessageId,
			MessageType: m.MessageType,
			Send: rsp.User{
				AccountID:  m.Send.AccountId,
				Username:   userMap[m.Send.AccountId].Username,
				HeadImgUrl: userMap[m.Send.AccountId].HeadImgUrl,
				Remark:     m.Send.Remark,
			},
			Receive:       m.Receive,
			Content:       m.Content,
			SendTimestamp: m.SendTimestamp,
			ReadStatus:    m.ReadStatus,
			ReadUserTotal: m.ReadUserTotal,
		}
		msgs = append(msgs, msg)
	}
	sort.Sort(rsp.Messages(msgs))
	var response = rsp.SessionMessageResp{
		Messages: msgs,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	go func() {
		for _, msg := range msgs {
			if msg.Send.AccountID != userMeta.AccountId && !msg.ReadStatus {
				var rm = rsp.ReadMessage{
					SessionId: int64(sessionId),
					MessageId: msg.MessageId,
				}
				sendMessageReadNotify(rm, userMeta.GetUser(), msg.Send.AccountID)
			}
		}
	}()
	handle.SuccessResp(c, "", response)
}

// 会话消息用户读取列表查询
func SessionMessageReadUsersQueryHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	sessionId, err := strconv.Atoi(c.Query("session_id"))
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var messageId int
	if messageId, err = strconv.Atoi(c.Query("message_id")); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var imClient im.DaoImClient
	if imClient, err = grpc.GetImClient(); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var resp *im.RspSessionMessageReadUsersQueryDao
	resp, err = imClient.SessionMessageReadUsersQueryDao(context.Background(), &im.ReqSessionMessageReadUsersQueryDao{
		AccountId: userMeta.AccountId,
		SessionId: int64(sessionId),
		MessageId: int64(messageId),
	})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var ids = append(resp.ReadUsers, resp.UnreadUsers...)
	userClient, err := grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var userResp *user.RspUserDaoById
	userResp, err = userClient.UserDaoById(context.Background(), &user.ReqUserDaoById{Ids: ids})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var idsMap = make(map[string]rsp.User)
	for _, u := range userResp.Users {
		idsMap[u.Id] = rsp.User{
			AccountID:  u.Id,
			Username:   u.Username,
			HeadImgUrl: u.HeadImgUrl,
			IsOfficial: u.IsOfficial,
		}
	}
	var readUsers, unreadUsers []rsp.User
	for _, u := range resp.ReadUsers {
		readUsers = append(readUsers, idsMap[u])
	}
	for _, u := range resp.UnreadUsers {
		unreadUsers = append(unreadUsers, idsMap[u])
	}
	var response = rsp.MessageReadUsers{
		MessageId:   resp.MessageId,
		ReadUsers:   readUsers,
		UnreadUsers: unreadUsers,
	}
	handle.SuccessResp(c, "", response)
}

// 会话消息撤回
func SessionMessageWithDrawnHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	sessionId, err := strconv.Atoi(c.Query("session_id"))
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var messageId int
	messageId, err = strconv.Atoi(c.Query("message_id"))
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var imClient im.DaoImClient
	imClient, err = grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	_, err = imClient.SessionMessageWithDrawnDao(context.Background(), &im.ReqSessionMessageWithDrawnDao{
		AccountId: userMeta.AccountId,
		SessionId: int64(sessionId),
		MessageId: int64(messageId),
	})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var resp *im.RspSessionDetailQueryDao
	resp, err = imClient.SessionDetailQueryDao(context.Background(), &im.ReqSessionDetailQueryDao{
		SessionId: int64(sessionId),
		AccountId: userMeta.AccountId,
	})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	go func() {
		for _, u := range resp.Joins {
			var msg = rsp.Message{
				SessionId:   int64(sessionId),
				MessageId:   int64(messageId),
				MessageType: im.SessionMessageType_WithDrawnMessage,
				Send:        userMeta.GetUser(),
				Receive:     u.AccountId,
			}
			sendWithDrawnMessageNotify(msg, u.AccountId)
		}
	}()
	handle.SuccessResp(c, "", nil)
}

// 会话消息已读状态更新
func SessionMessageReadStatusUpdateHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	sessionId, err := strconv.Atoi(c.Query("session_id"))
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	_, err = imClient.SessionMessageReadStatusUpdateDao(context.Background(), &im.ReqSessionMessageReadStatusUpdateDao{
		AccountId: userMeta.AccountId,
		SessionId: int64(sessionId),
	})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 会话消息单条消息已读更新
func SessionSingleMessageReadStatusUpdateHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	sessionId, err := strconv.Atoi(c.Query("session_id"))
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var messageId int
	messageId, err = strconv.Atoi(c.Query("message_id"))
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	accountId := c.Query("account_id")
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	_, err = imClient.SessionMessageReadStatusUpdateDao(context.Background(), &im.ReqSessionMessageReadStatusUpdateDao{
		AccountId:  userMeta.AccountId,
		SessionId:  int64(sessionId),
		MessageIds: []int64{int64(messageId)},
	})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	go func() {
		var rm = rsp.ReadMessage{
			SessionId: int64(sessionId),
			MessageId: int64(messageId),
		}
		sendMessageReadNotify(rm, userMeta.GetUser(), accountId)
	}()
	handle.SuccessResp(c, "", nil)
}

// 会话消息删除
func SessionMessageDeleteHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	sessionId, err := strconv.Atoi(c.Query("session_id"))
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	messageIdsStr := strings.Split(c.Query("message_ids"), ",")
	var messageIds = make([]int64, 0)
	for _, messageIdStr := range messageIdsStr {
		var messageId int
		if messageId, err = strconv.Atoi(messageIdStr); err != nil {
			log.Logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, handle.ParamErrResponse)
			return
		}
		messageIds = append(messageIds, int64(messageId))
	}
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	_, err = imClient.SessionMessageDeleteDao(context.Background(), &im.ReqSessionMessageDeleteDao{
		AccountId:  userMeta.AccountId,
		SessionId:  int64(sessionId),
		MessageIds: messageIds,
	})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 会话消息清空
func SessionMessageFlushHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	sessionId, err := strconv.Atoi(c.Query("session_id"))
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	_, err = imClient.SessionMessageFlushDao(context.Background(), &im.ReqSessionMessageFlushDao{
		AccountId: userMeta.AccountId,
		SessionId: int64(sessionId),
	})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}
