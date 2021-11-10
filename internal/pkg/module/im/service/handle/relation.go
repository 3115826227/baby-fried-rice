package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/im/grpc"
	"baby-fried-rice/internal/pkg/module/im/log"
	"baby-fried-rice/internal/pkg/module/im/service/webrtc"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
	"time"
)

// 发送操作通知
func sendOperatorNotify(client im.DaoImClient, operatorId int64, accountId string, sendToOrigin bool) {
	// 获取operator数据，发送给nsq，通知到前端
	userClient, err := grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	resp, err := client.OperatorSingleQueryDao(context.Background(), &im.ReqOperatorSingleQueryDao{
		AccountId:  accountId,
		OperatorId: operatorId,
	})
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var respUser *user.RspUserDaoById
	respUser, err = userClient.UserDaoById(context.Background(), &user.ReqUserDaoById{Ids: []string{resp.Origin}})
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	if len(respUser.Users) != 1 {
		err = constant.QueryUserByIdDaoError
		log.Logger.Error(err.Error())
		return
	}
	var origin = respUser.Users[0]
	var notify = models.WSMessageNotify{
		WSMessageNotifyType: constant.SessionMessageNotify,
		Timestamp:           time.Now().Unix(),
		WSMessage: models.WSMessage{
			SessionMessage: &models.SessionMessage{
				SessionMessageType: constant.OperatorMessage,
				Operator: rsp.Operator{
					OperatorId: resp.Id,
					Origin: rsp.User{
						AccountID:  origin.Id,
						HeadImgUrl: origin.HeadImgUrl,
						Username:   origin.Username,
					},
					Receive:      resp.Receive,
					OptType:      resp.OptType,
					Content:      resp.Content,
					NeedConfirm:  resp.NeedConfirm,
					Confirm:      resp.Confirm,
					OptTimestamp: resp.OptTimestamp,
				},
			},
		},
	}
	if sendToOrigin {
		// 发送给操作者
		notify.Receive = resp.Origin
	} else {
		// 发送给操作接收者
		notify.Receive = resp.Receive
		respUser, err = userClient.UserDaoById(context.Background(), &user.ReqUserDaoById{Ids: []string{resp.Origin}})
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		if len(respUser.Users) != 1 {
			err = constant.QueryUserByIdDaoError
			log.Logger.Error(err.Error())
			return
		}
		var u = respUser.Users[0]
		notify.WSMessage.Send = rsp.User{
			AccountID:  u.Id,
			HeadImgUrl: u.HeadImgUrl,
			Username:   u.Username,
		}
	}
	if err = mq.Send(topic, notify.ToString()); err != nil {
		return
	}
}

// 发送好友添加成功通知
func sendSessionNotify(sessionId int64, sessionType im.SessionType, joins []*im.JoinRemarkDao) {
	var now = time.Now().Unix()
	for _, join := range joins {
		var notify = models.WSMessageNotify{
			WSMessageNotifyType: 2,
			Receive:             join.AccountId,
			WSMessage: models.WSMessage{
				SessionMessage: &models.SessionMessage{
					SessionMessageType: constant.SessionMessage,
					Session: rsp.Session{
						SessionId:   sessionId,
						SessionType: sessionType,
						Name:        join.Remark,
					},
				},
			},
			Timestamp: now,
		}
		if err := mq.Send(topic, notify.ToString()); err != nil {
			continue
		}
	}
}

// 操作添加
func OperatorAddHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var req requests.ReqOperatorAdd
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
	var reqOperator = &im.ReqOperatorAddDao{
		Origin:      userMeta.AccountId,
		Receive:     req.Receive,
		OptType:     im.OptType(req.OptType),
		Content:     req.Content,
		NeedConfirm: req.NeedConfirm,
	}
	var resp *im.RspOperatorAddDao
	resp, err = imClient.OperatorAddDao(context.Background(), reqOperator)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	// 给需要确认的用户发送通知
	if reqOperator.NeedConfirm {
		go sendOperatorNotify(imClient, resp.OperatorId, userMeta.AccountId, false)
	}
	handle.SuccessResp(c, "", nil)
}

// 操作确认
func OperatorConfirmHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var req requests.ReqOperatorConfirm
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
	var reqOperator = &im.ReqOperatorConfirmDao{
		AccountId:  userMeta.AccountId,
		OperatorId: req.OperatorId,
		Confirm:    req.Confirm,
	}
	_, err = imClient.OperatorConfirmDao(context.Background(), reqOperator)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	// 获取operator数据，发送给nsq，通知到前端
	go sendOperatorNotify(imClient, req.OperatorId, userMeta.AccountId, true)
	if req.Confirm {
		// 如果同意了操作，需要根据具体情况做后续动作
		var od *im.OperatorDao
		od, err = imClient.OperatorSingleQueryDao(context.Background(), &im.ReqOperatorSingleQueryDao{
			OperatorId: req.OperatorId,
			AccountId:  userMeta.AccountId,
		})
		if err != nil {
			log.Logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
			return
		}
		switch od.OptType {
		case im.OptType_AddFriend:
			// 请求添加好友验证通过
			var remark string
			remark, err = GetUsername(od.Origin)
			if err != nil {
				c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
				return
			}
			var reqFriend = im.ReqFriendAddDao{
				Origin:     od.Origin,
				AccountId:  od.Receive,
				OperatorId: od.Id,
				Remark:     remark,
				OriRemark:  userMeta.Username,
			}
			_, err = imClient.FriendAddDao(context.Background(), &reqFriend)
			if err != nil {
				log.Logger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
				return
			}
			if err = addFriendSession(od.Origin, remark, od.Receive, userMeta.Username, imClient); err != nil {
				log.Logger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
				return
			}
		case im.OptType_JoinSession:
			// 会话加入成功
			var reqSession = im.ReqSessionJoinDao{
				AccountId:  od.Origin,
				SessionId:  od.SessionId,
				OperatorId: od.Id,
			}
			_, err = imClient.SessionJoinDao(context.Background(), &reqSession)
			if err != nil {
				log.Logger.Error(err.Error())
				c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
				return
			}
		}
	}
	handle.SuccessResp(c, "", nil)
}

// 操作读取状态更新
func OperatorReadStatusUpdateHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var req requests.ReqOperatorReadStatusUpdate
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
	var reqOperator = &im.ReqOperatorReadStatusUpdateDao{
		AccountId:   userMeta.AccountId,
		OperatorIds: req.Operators,
	}
	_, err = imClient.OperatorReadStatusUpdateDao(context.Background(), reqOperator)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 操作列表查询
func OperatorQueryHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	reqPage, err := handle.PageHandle(c)
	if err != nil {
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
	var reqOperator = &im.ReqOperatorsQueryDao{
		AccountId: userMeta.AccountId,
		Page:      reqPage.Page,
		PageSize:  reqPage.PageSize,
	}
	var resp *im.RspOperatorsQueryDao
	resp, err = imClient.OperatorsQueryDao(context.Background(), reqOperator)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var list = make([]rsp.Operator, 0)
	var userMap = make(map[string]*user.UserDao)
	var ids = make([]string, 0)
	for _, o := range resp.List {
		if _, exist := userMap[o.Origin]; !exist {
			userMap[o.Origin] = new(user.UserDao)
			ids = append(ids, o.Origin)
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
	for _, u := range userResp.Users {
		userMap[u.Id] = u
	}
	for _, o := range resp.List {
		var operator = rsp.Operator{
			OperatorId: o.Id,
			Origin: rsp.User{
				AccountID:  o.Origin,
				Username:   userMap[o.Origin].Username,
				HeadImgUrl: userMap[o.Origin].HeadImgUrl,
			},
			Receive:      o.Receive,
			OptType:      o.OptType,
			Content:      o.Content,
			NeedConfirm:  o.NeedConfirm,
			Confirm:      o.Confirm,
			OptTimestamp: o.OptTimestamp,
		}
		// 只有操作接收者才能看到该字段
		if userMeta.AccountId == operator.Receive {
			if o.ReceiveReadStatus {
				// 已读
				operator.ReceiveReadStatus = 2
			} else {
				// 未读
				operator.ReceiveReadStatus = 1
			}
		}
		list = append(list, operator)
	}
	var res = rsp.OperatorResp{Operators: list}
	handle.SuccessResp(c, "", res)
}

// 操作删除
func OperatorDeleteHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	operatorId, err := strconv.Atoi(c.Query("operator_id"))
	if err != nil {
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
	var reqOperator = &im.ReqOperatorDeleteDao{
		AccountId:  userMeta.AccountId,
		OperatorId: int64(operatorId),
	}
	_, err = imClient.OperatorDeleteDao(context.Background(), reqOperator)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

func GetUsername(accountId string) (username string, err error) {
	userClient, err := grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var resp *user.RspUserDaoById
	resp, err = userClient.UserDaoById(context.Background(), &user.ReqUserDaoById{Ids: []string{accountId}})
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	if len(resp.Users) != 1 {
		err = constant.QueryUserByIdDaoError
		log.Logger.Error(err.Error())
		return
	}
	username = resp.Users[0].Username
	return
}

// 添加好友
func FriendAddHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var req requests.ReqAddFriend
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, handle.ParamErrResponse)
		return
	}
	var remark string
	if req.Remark != "" {
		remark = req.Remark
	} else {
		// 获取用户名
		var err error
		remark, err = GetUsername(req.AccountId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
			return
		}
	}
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var reqFriend = &im.ReqFriendAddDao{
		Origin:    userMeta.AccountId,
		AccountId: req.AccountId,
		Remark:    remark,
		OriRemark: userMeta.Username,
	}
	_, err = imClient.FriendAddDao(context.Background(), reqFriend)
	if err != nil {
		if err.Error() == fmt.Sprintf("rpc error: code = Unknown desc = %v", constant.NeedApplyAddFriendError.Error()) {
			// 需要发出好友申请
			handle.ErrorResp(c, http.StatusOK, handle.CodeNeedApplyAddFriend, handle.CodeNeedApplyAddFriendMsg)
			return
		} else {
			log.Logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
			return
		}
	}
	// 好友添加成功后，创建会话，给会话成员发送会话通知
	if err = addFriendSession(userMeta.AccountId, userMeta.Username, req.AccountId, remark, imClient); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

func addFriendSession(origin, remark, friend, friendRemark string, imClient im.DaoImClient) error {
	var joins = []*im.JoinRemarkDao{
		{
			AccountId: origin,
			Remark:    remark,
		}, {
			AccountId: friend,
			Remark:    friendRemark,
		},
	}
	var reqSession = &im.ReqSessionAddDao{
		SessionType: im.SessionType_DoubleSession,
		Origin:      origin,
		Joins:       joins,
		Level:       im.SessionLevel_SessionBaseLevel,
	}
	respSession, err := imClient.SessionAddDao(context.Background(), reqSession)
	if err != nil {
		return err
	}
	go sendSessionNotify(respSession.SessionId, im.SessionType_DoubleSession, joins)
	return nil
}

// 好友列表查询
func FriendQueryHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	blackListStr := c.Query("black_list")
	remark := c.Query("remark")
	var blackList bool
	if blackListStr == "true" {
		blackList = true
	}
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var reqFriend = &im.ReqFriendQueryDao{
		Origin:     userMeta.AccountId,
		BlackList:  blackList,
		RemarkLike: remark,
	}
	var resp *im.RspFriendQueryDao
	resp, err = imClient.FriendQueryDao(context.Background(), reqFriend)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var list = make([]rsp.Friend, 0)
	for _, f := range resp.List {
		var friend = rsp.Friend{
			AccountId:  f.AccountId,
			Remark:     f.Remark,
			BlackList:  f.BlackList,
			Timestamp:  f.Timestamp,
			OnlineType: f.OnlineType,
		}
		list = append(list, friend)
	}
	sort.Sort(rsp.Friends(list))
	var res = rsp.FriendResp{Friends: list}
	handle.SuccessResp(c, "", res)
}

// 好友黑名单操作
func FriendBlackListUpdateHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var req requests.ReqUpdateFriendBlackList
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
	var reqFriend = &im.ReqFriendBlackListDao{
		Origin:    userMeta.AccountId,
		Friend:    req.Friend,
		BlackList: req.BlackList,
	}
	_, err = imClient.FriendBlackListDao(context.Background(), reqFriend)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 好友备注修改
func FriendRemarkUpdateHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var req requests.ReqUpdateFriendRemark
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
	var reqFriend = &im.ReqFriendRemarkDao{
		Origin: userMeta.AccountId,
		Friend: req.Friend,
		Remark: req.Remark,
	}
	_, err = imClient.FriendRemarkDao(context.Background(), reqFriend)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 好友关系删除
func FriendDeleteHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	friend := c.Query("friend")
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var reqFriend = &im.ReqFriendDeleteDao{
		Origin: userMeta.AccountId,
		Friend: friend,
	}
	_, err = imClient.FriendDeleteDao(context.Background(), reqFriend)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

func UserManageQueryHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var reqUserManage = &im.ReqUserManageQueryDao{
		AccountId: userMeta.AccountId,
	}
	var resp *im.RspUserManageQueryDao
	resp, err = imClient.UserManageQueryDao(context.Background(), reqUserManage)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var res = rsp.UserManageResp{
		AccountId:               resp.AccountId,
		AddFriendPermissionType: int32(resp.AddFriendPermissionType),
		UpdateTimestamp:         resp.UpdateTimestamp,
	}
	handle.SuccessResp(c, "", res)
}

func UserManageUpdateHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var req requests.ReqUserManageUpdate
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
	var reqUserManage = &im.ReqUserManageUpdateDao{
		AccountId:               userMeta.AccountId,
		AddFriendPermissionType: im.AddFriendPermissionType(req.AddFriendPermissionType),
	}
	_, err = imClient.UserManageUpdateDao(context.Background(), reqUserManage)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	handle.SuccessResp(c, "", nil)
}

// 创建webrtc
func SessionCreateWebRTC(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var req requests.ReqCreateWebRTC
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
	var imReq = im.ReqSessionDetailQueryDao{
		AccountId: userMeta.AccountId,
		SessionId: req.SessionId,
	}
	var resp *im.RspSessionDetailQueryDao
	resp, err = imClient.SessionDetailQueryDao(context.Background(), &imReq)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	var sdp string
	sdp, err = webrtc.CreateSession(req.Sdp, fmt.Sprintf("%v", req.SessionId))
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	go func() {
		//  邀请Session中好友视频通话
		for _, u := range resp.Joins {
			if u.AccountId == userMeta.AccountId {
				continue
			}
			var notify = models.WSMessageNotify{
				WSMessageNotifyType: constant.SessionMessageNotify,
				Receive:             u.AccountId,
				WSMessage: models.WSMessage{
					WSMessageType: im.SessionNotifyType_InviteNotify,
					Send:          userMeta.GetUser(),
					SessionMessage: &models.SessionMessage{
						SessionMessageType: constant.SessionMessageMessage,
						Message: rsp.Message{
							SessionId: req.SessionId,
						},
					},
				},
				Timestamp: time.Now().Unix(),
			}
			if err = mq.Send(topic, notify.ToString()); err != nil {
				log.Logger.Error(err.Error())
				return
			}
		}

	}()
	handle.SuccessResp(c, "", map[string]interface{}{
		"sdp": sdp,
	})
}

// webrtc通话加入回复
func SessionJoinWebRTC(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var req requests.ReqReturnWebRTC
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
	var imReq = im.ReqSessionDetailQueryDao{
		AccountId: userMeta.AccountId,
		SessionId: req.SessionId,
	}
	var resp *im.RspSessionDetailQueryDao
	resp, err = imClient.SessionDetailQueryDao(context.Background(), &imReq)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	if !req.Return {
		// 取消通话加入
		if resp.SessionType == im.SessionType_DoubleSession {

		}
	}
}
