package handle

import (
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/kit/models/requests"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"
	"baby-fried-rice/internal/pkg/module/im/cache"
	"baby-fried-rice/internal/pkg/module/im/grpc"
	"baby-fried-rice/internal/pkg/module/im/log"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type sessionWebRTCInfoStatus int64

const (
	//
	SessionWebRTCClose sessionWebRTCInfoStatus = 0
	// 通话已邀请等待回复中
	SessionWebRTCWaiting = 11
	// 邀请等待已超时
	SessionWebRTCTimeout = 12
	// 邀请已被拒绝
	SessionWebRTCReject = 12
	// 通话进行中
	SessionWebRTCPending = 100
)

type sessionWebRTCInfo struct {
	sessionId int64
	status    sessionWebRTCInfoStatus
	notify    chan rsp.SessionWebRTCUserStatus
}

var (
	sessionWebRTCChanLock sync.RWMutex
	sessionWebRTCChanMap  = make(map[int64]sessionWebRTCInfo)
)

func storageVideoMessage(msg models.WSMessageNotify, accountId string, videoTime int64) (messageId int64, err error) {
	var req = im.ReqSessionMessageAddDao{
		MessageType:   msg.WSMessage.SessionMessage.Message.MessageType,
		Send:          accountId,
		SessionId:     msg.WSMessage.SessionMessage.Message.SessionId,
		SendTimestamp: time.Now().Unix(),
	}
	var msgType = msg.WSMessage.SessionMessage.Message.MessageType
	if msgType == im.SessionMessageType_VideoLogMessage || msgType == im.SessionMessageType_AudioLogMessage {
		req.Content = fmt.Sprintf("%v", videoTime)
	}
	var imClient im.DaoImClient
	if imClient, err = grpc.GetImClient(); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var messageAddResp *im.RspSessionMessageAddDao
	if messageAddResp, err = imClient.SessionMessageAddDao(
		context.Background(), &req); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	messageId = messageAddResp.MessageId
	return
}

func doubleSessionInviteHandle(status rsp.SessionWebRTCUserStatus, info sessionWebRTCInfo) {
	var (
		after     = time.After(30 * time.Second)
		notify    models.WSMessageNotify
		sessionId = status.SessionId
		accountId = status.AccountId
	)
	sessionWebRTCChanLock.RLock()
	var exist bool
	info, exist = sessionWebRTCChanMap[sessionId]
	sessionWebRTCChanLock.RUnlock()
	if !exist {
		err := fmt.Errorf("session id isn't exist")
		log.Logger.Error(err.Error())
		return
	}
	var notifyChan = info.notify
	defer func() {
		close(info.notify)
		sessionWebRTCChanLock.Lock()
		delete(sessionWebRTCChanMap, sessionId)
		sessionWebRTCChanLock.Unlock()
		if err := cache.DeleteSessionWebRTC(sessionId); err != nil {
			log.Logger.Error(err.Error())
			return
		}
	}()
	for {
		select {
		case <-after:
			// 请求超时
			// 没有人接受，超过时间，判断为超时
			err := fmt.Errorf("user %v invite session %v video time out", accountId, sessionId)
			log.Logger.Error(err.Error())
			notify.WSMessage.WSMessageType = im.SessionNotifyType_DefaultNotify
			if status.Video {
				notify.WSMessage.SessionMessage.Message.MessageType = im.SessionMessageType_VideoNoReplyMessage
			} else {
				notify.WSMessage.SessionMessage.Message.MessageType = im.SessionMessageType_AudioNoReplyMessage
			}
			notify.WSMessage.SessionMessage.Message.SessionId = sessionId
			var messageId int64
			messageId, err = storageVideoMessage(notify, accountId, 0)
			notify.WSMessage.SessionMessage.Message.MessageId = messageId
			sendWebRTCNotify(notify)
			return
		case returnStatus := <-notifyChan:
			switch returnStatus.Status {
			case im.SessionNotifyType_ReceiveVideoMessage:
				// 同意通话
				sessionWebRTCChanLock.Lock()
				delete(sessionWebRTCChanMap, sessionId)
				sessionWebRTCChanLock.Unlock()
			case im.SessionNotifyType_RejectVideoMessage:
				// 拒绝通话
				err := fmt.Errorf("user %v invite session %v video had rejected", accountId, sessionId)
				log.Logger.Error(err.Error())
				notify.WSMessage.WSMessageType = im.SessionNotifyType_DefaultNotify
				if returnStatus.Video {
					notify.WSMessage.SessionMessage.Message.MessageType = im.SessionMessageType_VideoRejectMessage
				} else {
					notify.WSMessage.SessionMessage.Message.MessageType = im.SessionMessageType_AudioRejectMessage
				}
				notify.WSMessage.SessionMessage.Message.SessionId = sessionId
				var messageId int64
				if messageId, err = storageVideoMessage(notify, accountId, 0); err != nil {
					log.Logger.Error(err.Error())
					return
				}
				notify.WSMessage.SessionMessage.Message.MessageId = messageId
				sendWebRTCNotify(notify)
				if err = cache.DeleteSessionWebRTCTimeInfo(sessionId); err != nil {
					log.Logger.Error(err.Error())
					return
				}
			default:
				err := fmt.Errorf("invalid status")
				log.Logger.Error(err.Error())
			}
			return
		}
	}
}

// 检查用户webrtc状态
func checkUserWebRTCStatus(accountId string) (bool, error) {
	status, err := cache.GetUserOnlineStatus(accountId)
	if err != nil {
		return false, err
	}
	return status.VideoStatus, nil
}

// 用户发起音视频
func InviteVideoHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	sts, err := checkUserWebRTCStatus(userMeta.AccountId)
	if err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, handle.CodeInternalError)
		return
	}
	if sts {
		err = fmt.Errorf(handle.CodeSelfVideoConflictErrorMsg)
		log.Logger.Error(err.Error())
		handle.FailedResp(c, handle.CodeSelfVideoConflictError)
		return
	}
	var req requests.ReqCreateWebRTC
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusOK, handle.ParamErrResponse)
		return
	}
	var imClient im.DaoImClient
	if imClient, err = grpc.GetImClient(); err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, handle.CodeInternalError)
		return
	}
	var detailResp *im.RspSessionDetailQueryDao
	if detailResp, err = imClient.SessionDetailQueryDao(context.Background(), &im.ReqSessionDetailQueryDao{
		SessionId: req.SessionId,
		AccountId: userMeta.AccountId,
	}); err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, handle.CodeInternalError)
		return
	}
	if detailResp.SessionType == im.SessionType_DoubleSession {
		var other string
		for _, u := range detailResp.Joins {
			if u.AccountId != userMeta.AccountId {
				other = u.AccountId
				break
			}
		}
		var otherStatus models.UserOnlineStatus
		if otherStatus, err = cache.GetUserOnlineStatus(other); err != nil {
			log.Logger.Error(err.Error())
			handle.FailedResp(c, handle.CodeInternalError)
			return
		}
		// 用户已离线
		if otherStatus.OnlineType == im.OnlineStatusType_Offline {
			err = fmt.Errorf(handle.CodeUserOfflineErrorMsg)
			log.Logger.Error(err.Error())
			handle.FailedResp(c, handle.CodeSelfVideoConflictError)
			return
		}
		// 用户正在与其他用户通话中
		if otherStatus.VideoStatus {
			err = fmt.Errorf(handle.CodeUserVideoConflictErrorMsg)
			log.Logger.Error(err.Error())
			handle.FailedResp(c, handle.CodeSelfVideoConflictError)
			return
		}
	}
	var swapSdp string
	swapSdp, err = CreateSession(req.Sdp, req.SessionId, userMeta.AccountId, req.Video)
	if err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, handle.CodeInternalError)
		return
	}

	var status = rsp.SessionWebRTCUserStatus{
		SessionId: req.SessionId,
		AccountId: userMeta.AccountId,
		Status:    im.SessionNotifyType_InviteVideoNotify,
		Sdp:       req.Sdp,
		SwapSdp:   swapSdp,
	}
	if err = cache.UpdateSessionWebRTCUserStatus(status); err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, handle.CodeInternalError)
		return
	}

	var ch = make(chan rsp.SessionWebRTCUserStatus)
	var info = sessionWebRTCInfo{
		sessionId: req.SessionId,
		status:    SessionWebRTCWaiting,
		notify:    ch,
	}
	if detailResp.SessionType == im.SessionType_DoubleSession {
		doubleSessionInviteHandle(status, info)
	}

	handle.SuccessResp(c, "", swapSdp)
}

// 用户加入音视频
func JoinVideoHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var req requests.ReqJoinWebRTC
	if err := c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusOK, handle.ParamErrResponse)
		return
	}
	swapSdp, err := CreateSession(req.Sdp, req.SessionId, userMeta.AccountId, req.Video)
	if err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, handle.CodeInternalError)
		return
	}
	var status = rsp.SessionWebRTCUserStatus{
		SessionId: req.SessionId,
		AccountId: userMeta.AccountId,
		Status:    im.SessionNotifyType_ReceiveVideoMessage,
		Sdp:       req.Sdp,
		SwapSdp:   swapSdp,
	}
	if err = cache.UpdateSessionWebRTCUserStatus(status); err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, handle.CodeInternalError)
		return
	}
	handle.SuccessResp(c, "", swapSdp)
}

// 用户回应音视频(同意/拒绝)
func ReturnVideoHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	var err error
	var req requests.ReqReturnWebRTC
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusOK, handle.ParamErrResponse)
		return
	}
	var status = rsp.SessionWebRTCUserStatus{
		SessionId: req.SessionId,
		AccountId: userMeta.AccountId,
	}
	if req.Return {
		status.Status = im.SessionNotifyType_ReceiveVideoMessage
	} else {
		status.Status = im.SessionNotifyType_RejectVideoMessage
	}
	if err = cache.UpdateSessionWebRTCUserStatus(status); err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, handle.CodeInternalError)
		return
	}
	if err = cache.SetSessionWebRTCTimeInfo(req.SessionId, req.Video); err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, handle.CodeInternalError)
		return
	}
	sessionWebRTCChanLock.RLock()
	info, exist := sessionWebRTCChanMap[req.SessionId]
	sessionWebRTCChanLock.RUnlock()
	if !exist {
		err = fmt.Errorf("session id isn't exist")
		log.Logger.Error(err.Error())
		handle.FailedResp(c, handle.CodeInternalError)
		return
	}
	info.notify <- status
	if req.Return {
		// 同意 用户创建自己的WebRTC会话
		var swapSdp string
		swapSdp, err = CreateSession(req.Sdp, req.SessionId,
			userMeta.AccountId, req.Video)
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		handle.SuccessResp(c, "", swapSdp)
	} else {
		handle.SuccessResp(c, "", nil)
	}
}

// 用户交换webrtc sdp
func SwapWebRTCSdpHandle(c *gin.Context) {
	var err error
	var req requests.ReqSwapWebRTCSdp
	if err = c.ShouldBind(&req); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusOK, handle.ParamErrResponse)
		return
	}
	var remoteSwapSdp string
	remoteSwapSdp, err = JoinSession(req.RemoteSdp, req.SessionId,
		req.Origin, req.Video)
	if err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, handle.CodeInternalError)
		return
	}
	handle.SuccessResp(c, "", remoteSwapSdp)
}

// 用户(挂断)退出音视频通话
func HangupVideoHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	sessionId, err := strconv.Atoi("session_id")
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusOK, handle.ParamErrResponse)
		return
	}
	if err = cache.RemoveSessionWebRTCUserStatus(int64(sessionId), userMeta.AccountId); err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, handle.CodeInternalError)
		return
	}
	var imClient im.DaoImClient
	if imClient, err = grpc.GetImClient(); err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, handle.CodeInternalError)
		return
	}
	var detailResp *im.RspSessionDetailQueryDao
	if detailResp, err = imClient.SessionDetailQueryDao(context.Background(), &im.ReqSessionDetailQueryDao{
		SessionId: int64(sessionId),
		AccountId: userMeta.AccountId,
	}); err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, handle.CodeInternalError)
		return
	}
	var notify models.WSMessageNotify
	if detailResp.SessionType == im.SessionType_DoubleSession {
		// 如果是双人会话，就直接结束该视频会话，保存通话记录
		var timeInfo rsp.SessionWebRTCTimeInfo
		timeInfo, err = cache.GetSessionWebRTCTimeInfo(int64(sessionId))
		if err != nil {
			log.Logger.Error(err.Error())
			handle.FailedResp(c, handle.CodeInternalError)
			return
		}
		notify.WSMessage.WSMessageType = im.SessionNotifyType_HangupVideoMessage
		if timeInfo.Video {
			notify.WSMessage.SessionMessage.Message.MessageType = im.SessionMessageType_VideoLogMessage
		} else {
			notify.WSMessage.SessionMessage.Message.MessageType = im.SessionMessageType_AudioLogMessage
		}

		notify.WSMessage.SessionMessage.Message.SessionId = int64(sessionId)
		var messageId int64
		if messageId, err = storageVideoMessage(notify, userMeta.AccountId, time.Now().Unix()-timeInfo.StartTime); err != nil {
			log.Logger.Error(err.Error())
			handle.FailedResp(c, handle.CodeInternalError)
			return
		}
		for _, u := range detailResp.Joins {
			if u.AccountId == userMeta.AccountId {
				continue
			}
			var newNotify = notify
			newNotify.Receive = u.AccountId
			newNotify.WSMessage.WSMessageType = im.SessionNotifyType_HangupVideoMessage
			newNotify.WSMessage.SessionMessage.Message.MessageId = messageId
			sendWebRTCNotify(newNotify)
		}
		if err = cache.DeleteSessionWebRTC(int64(sessionId)); err != nil {
			log.Logger.Error(err.Error())
			handle.FailedResp(c, handle.CodeInternalError)
			return
		}
		if err = cache.DeleteSessionWebRTCTimeInfo(int64(sessionId)); err != nil {
			log.Logger.Error(err.Error())
			handle.FailedResp(c, handle.CodeInternalError)
			return
		}
		sessionWebRTCChanLock.Lock()
		delete(sessionWebRTCChanMap, int64(sessionId))
		sessionWebRTCChanLock.Unlock()
	} else {
		// todo 如果是多人视频会话，则另做处理
	}
	handle.SuccessResp(c, "", nil)
}

// 查询当前webrtc状态
func VideoStatusHandle(c *gin.Context) {
	sessionId, err := strconv.Atoi("session_id")
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusOK, handle.ParamErrResponse)
		return
	}
	var mp map[string]rsp.SessionWebRTCUserStatus
	mp, err = cache.GetSessionWebRTC(int64(sessionId))
	if err != nil {
		log.Logger.Error(err.Error())
		handle.FailedResp(c, handle.CodeInternalError)
		return
	}
	var users = make([]rsp.SessionWebRTCUserStatus, 0)
	for _, user := range mp {
		users = append(users, user)
	}
	var resp = rsp.SessionWebRTCInfo{
		Users: users,
	}
	if len(users) > 0 {
		resp.Status = true
	}
	handle.SuccessResp(c, "", resp)
}
