package handle

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/im/grpc"
	"baby-fried-rice/internal/pkg/module/im/log"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

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
	var reqSession = &im.ReqSessionAddDao{
		SessionType:        im.SessionType(req.SessionType),
		JoinPermissionType: im.SessionJoinPermissionType(req.JoinPermissionType),
		Name:               req.Name,
		Origin:             userMeta.AccountId,
		Joins:              req.Joins,
	}
	_, err = imClient.SessionAddDao(context.Background(), reqSession)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 会话列表查询
func SessionQueryHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var reqSession = &im.ReqSessionQueryDao{
		AccountId: userMeta.AccountId,
	}
	var resp *im.RspSessionQueryDao
	resp, err = imClient.SessionQueryDao(context.Background(), reqSession)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var sessionAccountIdMap = make(map[int64]string)
	var idsMap = make(map[string]struct{}, 0)
	var ids = make([]string, 0)
	for _, s := range resp.Sessions {
		if s.Name == "" && s.SessionType == im.SessionType_DoubleSession {
			for _, accountId := range s.Joins {
				if accountId != userMeta.AccountId {
					sessionAccountIdMap[s.SessionId] = accountId
					idsMap[accountId] = struct{}{}
					ids = append(ids, accountId)
					break
				}
			}
		}
		if s.Latest != nil {
			accountId := s.Latest.Send
			if _, exist := idsMap[accountId]; !exist {
				idsMap[accountId] = struct{}{}
				ids = append(ids, accountId)
			}
		}
	}
	var userClient user.DaoUserClient
	userClient, err = grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var usersResp *user.RspUserDaoById
	usersResp, err = userClient.UserDaoById(context.Background(), &user.ReqUserDaoById{Ids: ids})
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var userMap = make(map[string]*user.UserDao)
	for _, u := range usersResp.Users {
		userMap[u.Id] = u
	}
	var sessions = make([]rsp.Session, 0)
	for _, s := range resp.Sessions {
		var session = rsp.Session{
			SessionId:   s.SessionId,
			SessionType: s.SessionType,
			Name:        s.Name,
			Origin:      s.Origin,
			Unread:      s.Unread,
			CreateTime:  s.CreateTime,
			Joins:       s.Joins,
		}
		if accountId, exist := sessionAccountIdMap[s.SessionId]; exist {
			session.Name = userMap[accountId].Username
		}
		if s.Latest != nil {
			lm := s.Latest
			session.LatestMessage = &rsp.Message{
				SessionId:   lm.SessionId,
				MessageId:   lm.MessageId,
				MessageType: lm.MessageType,
				Send: rsp.User{
					AccountID: lm.Send,
					Username:  userMap[lm.Send].Username,
				},
				Receive:       lm.Receive,
				Content:       lm.Content,
				SendTimestamp: lm.SendTimestamp,
				ReadStatus:    lm.ReadStatus,
			}
		}
		sessions = append(sessions, session)
	}
	var res = rsp.SessionsResp{Sessions: sessions}
	handle.SuccessResp(c, "", res)
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
	handle.SuccessResp(c, "", resp)
}

func SessionUpdateHandle(c *gin.Context) {

}

// 加入会话
func SessionJoinHandle(c *gin.Context) {

}

// 离开会话
func SessionLeaveHandle(c *gin.Context) {

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
	handle.SuccessResp(c, "", resp)
}

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
