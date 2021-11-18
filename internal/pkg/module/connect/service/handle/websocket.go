package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/models"
	"baby-fried-rice/internal/pkg/kit/models/rsp"
	"baby-fried-rice/internal/pkg/kit/mq/nsq"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/im"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/user"
	"baby-fried-rice/internal/pkg/module/connect/cache"
	"baby-fried-rice/internal/pkg/module/connect/config"
	"baby-fried-rice/internal/pkg/module/connect/grpc"
	"baby-fried-rice/internal/pkg/module/connect/log"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

var (
	upGrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ConnectionMap         = make(map[string]*websocket.Conn)
	mq                    interfaces.MQ
	writeChan             = make(chan models.WSMessageNotify, 2000)
	sessionWebRTCChanLock sync.RWMutex
	sessionWebRTCChanMap  = make(map[int64]sessionWebRTCInfo)
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

func (info *sessionWebRTCInfo) close() {
	info.status = SessionWebRTCClose
	close(info.notify)
}

func Init() {
	go handleWrite()

	conf := config.GetConfig()
	mq = nsq.InitNSQMQ(conf.MessageQueue.NSQ.Cluster)
	tc := conf.MessageQueue.ConsumeTopics.WebsocketNotify
	err := mq.NewConsumer(tc.Topic, tc.Channel)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	go runConsume()
}

func handleWrite() {
	for {
		select {
		case msg := <-writeChan:
			conn, exist := ConnectionMap[msg.Receive]
			if exist {
				if err := conn.WriteJSON(msg); err != nil {
					log.Logger.Error(err.Error())
					continue
				}
				log.Logger.Debug(fmt.Sprintf("write smsDao %v to %v success", msg.WSMessage.ToString(), msg.Receive))
			}
		}
	}
}

// 消费MQ，推送给前端
func runConsume() {
	for {
		value, err := mq.Consume()
		if err != nil {
			log.Logger.Error(err.Error())
			continue
		}
		var msg models.WSMessageNotify
		if err = json.Unmarshal([]byte(value), &msg); err != nil {
			log.Logger.Error(err.Error())
			continue
		}
		writeChan <- msg
	}
}

func WebSocketHandle(c *gin.Context) {
	userMeta := handle.GetUserMeta(c)
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Logger.Error("connect failed")
		return
	}
	closeChan := make(chan bool, 1)
	defer func() {
		// 用户断开连接
		if err = cache.UpdateUserOnlineStatus(userMeta.AccountId, im.OnlineStatusType_Offline); err != nil {
			log.Logger.Error(err.Error())
		}
		closeChan <- true
		conn.Close()
	}()
	log.Logger.Info(userMeta.AccountId + " connect success")
	ConnectionMap[userMeta.AccountId] = conn
	if err = cache.UpdateUserOnlineStatus(userMeta.AccountId, im.OnlineStatusType_PCOnline); err != nil {
		log.Logger.Error(err.Error())
	}
	for {
		var msg models.WSMessageNotify
		if err = conn.ReadJSON(&msg); err != nil {
			log.Logger.Error(err.Error())
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				return
			}
			continue
		}
		switch msg.WSMessageNotifyType {
		case constant.SessionMessageNotify:
			switch msg.WSMessage.SessionMessage.SessionMessageType {
			case constant.SessionMessage:
				handleSession(msg, userMeta)
			case constant.SessionMessageMessage:
				handleSessionMessage(msg, userMeta.AccountId)
			}
		default:
			continue
		}
	}
}

func storageVideoMessage(msg models.WSMessageNotify, accountId string, videoTime int64) (messageId int64, err error) {
	var req = im.ReqSessionMessageAddDao{
		MessageType:   msg.WSMessage.SessionMessage.Message.MessageType,
		Send:          accountId,
		SessionId:     msg.WSMessage.SessionMessage.Message.SessionId,
		SendTimestamp: time.Now().Unix(),
	}
	if msg.WSMessage.SessionMessage.Message.MessageType == im.SessionMessageType_VideoLogMessage {
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

func sendWebRTCNotify(resp *im.RspSessionDetailQueryDao, notify models.WSMessageNotify, accountId string) {
	// 给邀请者发送通知
	if resp.SessionType == im.SessionType_DoubleSession {
		for _, u := range resp.Joins {
			if u.AccountId != accountId {
				notify.Receive = u.AccountId
				writeChan <- notify
				break
			}
		}
	} else {
		var inviteUserMap = make(map[string]struct{})
		for _, u := range notify.WSMessage.SessionMessage.WebRtc.InviteUsers {
			inviteUserMap[u] = struct{}{}
		}
		for _, u := range resp.Joins {
			if u.AccountId == accountId {
				continue
			}
			if _, exist := inviteUserMap[u.AccountId]; exist {
				var newNotify = notify
				newNotify.Receive = u.AccountId
				writeChan <- newNotify
			}
		}
	}
}

// 会话通知消息
func handleSession(msg models.WSMessageNotify, userMeta *handle.UserMeta) {
	var accountId = userMeta.AccountId
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var notify = msg
	var sessionId = msg.WSMessage.SessionMessage.Session.SessionId
	switch msg.WSMessage.WSMessageType {
	case im.SessionNotifyType_InputtingMessage:
		// 对方正在输入
		for _, u := range notify.WSMessage.SessionMessage.Session.Users {
			if u.AccountID != accountId {
				notify.Receive = u.AccountID
				break
			}
		}
		writeChan <- notify
	case im.SessionNotifyType_OnlineStatus:
		// 获取会话中用户在线状态
		notify = models.WSMessageNotify{
			WSMessageNotifyType: msg.WSMessageNotifyType,
			Receive:             accountId,
			WSMessage:           msg.WSMessage,
			Timestamp:           msg.Timestamp,
		}
		var users = make([]rsp.User, 0)
		var resp *im.RspSessionDetailQueryDao
		resp, err = imClient.SessionDetailQueryDao(context.Background(), &im.ReqSessionDetailQueryDao{
			SessionId: sessionId,
		})
		if err != nil {
			log.Logger.Error(err.Error())
			notify.WSMessage.SessionMessage.Session.SessionId = sessionId
			writeChan <- notify
			return
		}
		for _, u := range resp.Joins {
			users = append(users, rsp.User{
				AccountID:  u.AccountId,
				Remark:     u.Remark,
				OnlineType: u.OnlineType,
			})
		}
		notify.WSMessage.SessionMessage.Session = rsp.Session{
			SessionId: resp.SessionId,
			Users:     users,
		}
		writeChan <- notify
	case im.SessionNotifyType_InviteVideoNotify:
		// 用户邀请视频
		notify = models.WSMessageNotify{
			WSMessageNotifyType: msg.WSMessageNotifyType,
			WSMessage:           msg.WSMessage,
			Timestamp:           msg.Timestamp,
		}
		notify.WSMessage.SessionMessage.Session.SessionId = sessionId
		// 邀请用户创建视频会话
		var swapSdp string
		swapSdp, err = CreateSession(msg.WSMessage.SessionMessage.WebRtc.Sdp, sessionId, accountId)
		if err != nil {
			log.Logger.Error(err.Error())
			notify.WSMessage.WSMessageType = im.SessionNotifyType_InviteVideoFailedNotify
			writeChan <- notify
			return
		}
		notify.WSMessage.SessionMessage.WebRtc.InviteAccount = accountId
		notify.WSMessage.SessionMessage.WebRtc.SwapSdp = swapSdp
		// 给邀请用户发送本地sdp回执
		notify.WSMessage.WSMessageType = im.SessionNotifyType_LocalVideoMessage
		notify.Receive = accountId
		writeChan <- notify

		// 给其他受邀请用户发送邀请，并带上sdp
		var resp *im.RspSessionDetailQueryDao
		resp, err = imClient.SessionDetailQueryDao(context.Background(), &im.ReqSessionDetailQueryDao{
			SessionId: sessionId,
		})
		if err != nil {
			log.Logger.Error(err.Error())
			notify.WSMessage.WSMessageType = im.SessionNotifyType_InviteVideoFailedNotify
			writeChan <- notify
			return
		}
		notify.WSMessage.WSMessageType = im.SessionNotifyType_InviteVideoNotify
		sendWebRTCNotify(resp, notify, accountId)

		var ch = make(chan rsp.SessionWebRTCUserStatus, len(msg.WSMessage.SessionMessage.WebRtc.InviteUsers)+2)
		var info = sessionWebRTCInfo{
			sessionId: sessionId,
			status:    SessionWebRTCWaiting,
			notify:    ch,
		}
		var status = rsp.SessionWebRTCUserStatus{
			SessionId: sessionId,
			AccountId: accountId,
			Status:    im.SessionNotifyType_InviteVideoNotify,
			Sdp:       msg.WSMessage.SessionMessage.WebRtc.Sdp,
			SwapSdp:   swapSdp,
		}
		if err = cache.UpdateSessionWebRTCUserStatus(status); err != nil {
			log.Logger.Error(err.Error())
			return
		}
		sessionWebRTCChanLock.Lock()
		sessionWebRTCChanMap[sessionId] = info
		info.notify <- status
		sessionWebRTCChanLock.Unlock()
		var after = time.After(30 * time.Second)
		for {
			select {
			case <-after:
				// 请求超时
				// 没有人接受，超过时间，判断为超时
				err = fmt.Errorf("user %v invite session %v video time out", accountId, sessionId)
				log.Logger.Error(err.Error())
				notify.WSMessage.WSMessageType = im.SessionNotifyType_DefaultNotify
				notify.WSMessage.SessionMessage.Message.MessageType = im.SessionMessageType_VideoNoReplyMessage
				notify.WSMessage.SessionMessage.Message.SessionId = sessionId
				var messageId int64
				messageId, err = storageVideoMessage(notify, accountId, 0)
				notify.WSMessage.SessionMessage.Message.MessageId = messageId
				writeChan <- notify
				sessionWebRTCChanLock.Lock()
				delete(sessionWebRTCChanMap, sessionId)
				sessionWebRTCChanLock.Unlock()
				if err = cache.DeleteSessionWebRTC(sessionId); err != nil {
					log.Logger.Error(err.Error())
					return
				}
				return
			default:
				// webrtc会话通知
				// 更新状态
				sessionWebRTCChanLock.RLock()
				var exist bool
				if info, exist = sessionWebRTCChanMap[sessionId]; !exist {
					break
				}
				if len(info.notify) == 0 {
					sessionWebRTCChanLock.RUnlock()
					break
				}
				status = <-info.notify
				sessionWebRTCChanLock.RUnlock()
				// 判断是否可以提前结束视频会话邀请
				var receive bool
				var number int
				if resp.SessionType == im.SessionType_DoubleSession {
					number = 2
				} else {
					number = len(msg.WSMessage.SessionMessage.WebRtc.InviteUsers)
				}
				var statusMap map[string]rsp.SessionWebRTCUserStatus
				if statusMap, err = cache.GetSessionWebRTC(sessionId); err != nil {
					return
				}
				receive, err = cache.JudgeSessionReceiveWebRTC(statusMap)
				if err != nil {
					log.Logger.Error(err.Error())
					return
				}
				if !receive && len(statusMap) == number {
					// 会话已被拒绝
					err = fmt.Errorf("user %v invite session %v video had rejected", accountId, sessionId)
					log.Logger.Error(err.Error())
					notify.WSMessage.WSMessageType = im.SessionNotifyType_DefaultNotify
					notify.WSMessage.SessionMessage.Message.MessageType = im.SessionMessageType_VideoRejectMessage
					notify.WSMessage.SessionMessage.Message.SessionId = sessionId
					var messageId int64
					if messageId, err = storageVideoMessage(notify, accountId, 0); err != nil {
						log.Logger.Error(err.Error())
						return
					}
					notify.WSMessage.SessionMessage.Message.MessageId = messageId
					writeChan <- notify
					sessionWebRTCChanLock.Lock()
					delete(sessionWebRTCChanMap, sessionId)
					sessionWebRTCChanLock.Unlock()
					if err = cache.DeleteSessionWebRTC(sessionId); err != nil {
						log.Logger.Error(err.Error())
						return
					}
					if err = cache.DeleteSessionWebRTCStartTime(sessionId); err != nil {
						log.Logger.Error(err.Error())
						return
					}
					return
				} else if receive {
					// 通话已经接通
					sessionWebRTCChanLock.Lock()
					delete(sessionWebRTCChanMap, sessionId)
					sessionWebRTCChanLock.Unlock()
					return
				}
			}
		}

	case im.SessionNotifyType_ReceiveVideoMessage:
		// 接受视频通话
		// 接受用户加入邀请用户的视频通话
		var remotesSwapSdp string
		remotesSwapSdp, err = JoinSession(msg.WSMessage.SessionMessage.WebRtc.RemoteSdp, sessionId, notify.WSMessage.SessionMessage.WebRtc.InviteAccount)
		if err != nil {
			log.Logger.Error(err.Error())
			notify.WSMessage.WSMessageType = im.SessionNotifyType_InviteVideoFailedNotify
			writeChan <- notify
			return
		}
		notify.WSMessage.SessionMessage.WebRtc.RemoteSwapSdp = remotesSwapSdp
		// 给自己发送会话加入的远程回执
		notify.WSMessage.WSMessageType = im.SessionNotifyType_RemoteVideoMessage
		notify.Receive = accountId
		writeChan <- notify

		// 发送接受会话通知
		notify.WSMessage.WSMessageType = im.SessionNotifyType_ReceiveVideoMessage
		var status = rsp.SessionWebRTCUserStatus{
			SessionId:     sessionId,
			AccountId:     accountId,
			Status:        im.SessionNotifyType_ReceiveVideoMessage,
			RemoteSdp:     msg.WSMessage.SessionMessage.WebRtc.RemoteSdp,
			RemoteSwapSdp: remotesSwapSdp,
		}
		notify.WSMessage.Send = userMeta.GetUser()
		if err = cache.UpdateSessionWebRTCUserStatus(status); err != nil {
			log.Logger.Error(err.Error())
			return
		}
		sessionWebRTCChanLock.RLock()
		if info, exist := sessionWebRTCChanMap[sessionId]; exist {
			info.notify <- status
		}
		sessionWebRTCChanLock.RUnlock()
		// 判断当前会话视频通话是否已经进行
		if _, err = cache.GetSessionWebRTCStartTime(sessionId); err != nil {
			if err == redis.Nil {
				// 如果Redis中不存在则表示还未进行，初始化通话开始时间
				if err = cache.SetSessionWebRTCStartTime(sessionId); err != nil {
					log.Logger.Error(err.Error())
					return
				}
			}
		}
		var resp *im.RspSessionDetailQueryDao
		resp, err = imClient.SessionDetailQueryDao(context.Background(), &im.ReqSessionDetailQueryDao{
			SessionId: sessionId,
		})
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		notify.WSMessage.SessionMessage.WebRtc.InviteAccount = accountId
		// 接受用户创建视频会话
		var ownSessionSwapSdp string
		ownSessionSwapSdp, err = CreateSession(msg.WSMessage.SessionMessage.WebRtc.Sdp, sessionId, accountId)
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		sendWebRTCNotify(resp, notify, accountId)

		// 给自己发送本地回执
		notify.WSMessage.SessionMessage.WebRtc.SwapSdp = ownSessionSwapSdp
		notify.WSMessage.WSMessageType = im.SessionNotifyType_LocalVideoMessage
		notify.Receive = accountId
		writeChan <- notify
	case im.SessionNotifyType_RejectVideoMessage:
		// 拒绝视频通话
		notify.WSMessage.WSMessageType = im.SessionNotifyType_RejectVideoMessage
		var status = rsp.SessionWebRTCUserStatus{
			SessionId: sessionId,
			AccountId: accountId,
			Status:    im.SessionNotifyType_RejectVideoMessage,
		}
		if err = cache.UpdateSessionWebRTCUserStatus(status); err != nil {
			log.Logger.Error(err.Error())
			return
		}
		sessionWebRTCChanLock.RLock()
		if info, exist := sessionWebRTCChanMap[sessionId]; exist {
			info.notify <- status
		}
		sessionWebRTCChanLock.RUnlock()
		var resp *im.RspSessionDetailQueryDao
		resp, err = imClient.SessionDetailQueryDao(context.Background(), &im.ReqSessionDetailQueryDao{
			SessionId: sessionId,
		})
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		sendWebRTCNotify(resp, notify, accountId)
	case im.SessionNotifyType_HangupVideoMessage:
		// 挂断视频通话
		log.Logger.Debug("hangup video 挂断通话")
		var status = rsp.SessionWebRTCUserStatus{
			SessionId: sessionId,
			AccountId: accountId,
			Status:    im.SessionNotifyType_HangupVideoMessage,
		}
		if err = cache.UpdateSessionWebRTCUserStatus(status); err != nil {
			log.Logger.Error(err.Error())
			return
		}
		var resp *im.RspSessionDetailQueryDao
		resp, err = imClient.SessionDetailQueryDao(context.Background(), &im.ReqSessionDetailQueryDao{
			SessionId: sessionId,
		})
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		if resp.SessionType == im.SessionType_DoubleSession {
			// 如果是双人会话，就直接结束该视频会话，保存通话记录
			var startTime int64
			startTime, err = cache.GetSessionWebRTCStartTime(sessionId)
			if err != nil {
				log.Logger.Error(err.Error())
				return
			}
			notify.WSMessage.WSMessageType = im.SessionNotifyType_HangupVideoMessage
			notify.WSMessage.SessionMessage.Message.MessageType = im.SessionMessageType_VideoLogMessage
			notify.WSMessage.SessionMessage.Message.SessionId = sessionId
			var messageId int64
			if messageId, err = storageVideoMessage(notify, accountId, time.Now().Unix()-startTime); err != nil {
				log.Logger.Error(err.Error())
				return
			}
			for _, u := range resp.Joins {
				if u.AccountId == accountId {
					continue
				}
				var newNotify = notify
				newNotify.Receive = u.AccountId
				newNotify.WSMessage.WSMessageType = im.SessionNotifyType_HangupVideoMessage
				newNotify.WSMessage.SessionMessage.Message.MessageId = messageId
				writeChan <- newNotify
			}
			if err = cache.DeleteSessionWebRTC(sessionId); err != nil {
				log.Logger.Error(err.Error())
				return
			}
			if err = cache.DeleteSessionWebRTCStartTime(sessionId); err != nil {
				log.Logger.Error(err.Error())
				return
			}
			sessionWebRTCChanLock.Lock()
			delete(sessionWebRTCChanMap, sessionId)
			sessionWebRTCChanLock.Unlock()
		} else {
			// todo 如果是多人视频会话，则另做处理
		}
	case im.SessionNotifyType_JoinVideoMessage:
		fmt.Println("join video")
		var remotesSwapSdp string
		remotesSwapSdp, err = JoinSession(msg.WSMessage.SessionMessage.WebRtc.RemoteSdp, sessionId, notify.WSMessage.SessionMessage.WebRtc.InviteAccount)
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		notify.WSMessage.SessionMessage.WebRtc.RemoteSwapSdp = remotesSwapSdp
		// 给自己发送会话加入的远程回执
		notify.WSMessage.WSMessageType = im.SessionNotifyType_RemoteVideoMessage
		notify.Receive = accountId
		writeChan <- notify
	}
}

func handleSessionMessage(msg models.WSMessageNotify, accountId string) {
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var userClient user.DaoUserClient
	userClient, err = grpc.GetUserClient()
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	// 从accountDao服务获取消息发送者的用户名和头像
	var userResp *user.RspDaoUserDetail
	userResp, err = userClient.UserDaoDetail(context.Background(),
		&user.ReqDaoUserDetail{AccountId: accountId})
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	msg.WSMessage.Send.AccountID = userResp.Detail.AccountId
	msg.WSMessage.Send.Username = userResp.Detail.Username
	msg.WSMessage.Send.HeadImgUrl = userResp.Detail.HeadImgUrl
	msg.Timestamp = time.Now().Unix()
	// 从imDao服务中获取会话中的所有用户id
	var resp *im.RspSessionDetailQueryDao
	resp, err = imClient.SessionDetailQueryDao(context.Background(),
		&im.ReqSessionDetailQueryDao{SessionId: msg.WSMessage.SessionMessage.Message.SessionId})
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	// 将会话消息发给imDao服务存入数据库中
	var req = &im.ReqSessionMessageAddDao{
		MessageType:   msg.WSMessage.SessionMessage.Message.MessageType,
		Send:          msg.WSMessage.Send.AccountID,
		SessionId:     msg.WSMessage.SessionMessage.Message.SessionId,
		Content:       msg.WSMessage.SessionMessage.Message.Content,
		SendTimestamp: msg.Timestamp,
	}
	var messageAddResp *im.RspSessionMessageAddDao
	messageAddResp, err = imClient.SessionMessageAddDao(context.Background(), req)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	// 发送给nsq
	for _, u := range resp.Joins {
		var notify = models.WSMessageNotify{
			WSMessageNotifyType: msg.WSMessageNotifyType,
			Receive:             u.AccountId,
			WSMessage:           msg.WSMessage,
			Timestamp:           msg.Timestamp,
		}
		notify.WSMessage.SessionMessage.Message = rsp.Message{
			MessageId:   messageAddResp.MessageId,
			SessionId:   resp.SessionId,
			MessageType: msg.WSMessage.SessionMessage.Message.MessageType,
			Send: rsp.User{
				AccountID:  userResp.Detail.AccountId,
				Username:   userResp.Detail.Username,
				HeadImgUrl: userResp.Detail.HeadImgUrl,
			},
			Receive:       u.AccountId,
			Content:       msg.WSMessage.SessionMessage.Message.Content,
			SendTimestamp: msg.Timestamp,
		}
		writeChan <- notify
	}
}
