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
	"github.com/gorilla/websocket"
	"net/http"
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
	ConnectionMap = make(map[string]*websocket.Conn)
	mq            interfaces.MQ
	writeChan     = make(chan models.WSMessageNotify, 2000)
)

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
				handleSession(msg, userMeta.AccountId)
			case constant.SessionMessageMessage:
				handleSessionMessage(msg, userMeta.AccountId)
			}
		default:
			continue
		}
	}
}

// 会话通知消息
func handleSession(msg models.WSMessageNotify, accountId string) {
	imClient, err := grpc.GetImClient()
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	var notify models.WSMessageNotify
	switch msg.WSMessage.WSMessageType {
	case im.SessionNotifyType_InputtingMessage:
		// 对方正在输入
		notify = msg
		for _, u := range notify.WSMessage.SessionMessage.Session.Users {
			if u.AccountID != accountId {
				notify.Receive = u.AccountID
			}
		}
	case im.SessionNotifyType_OnlineStatus:
		// 获取会话中用户在线状态
		var resp *im.RspSessionDetailQueryDao
		resp, err = imClient.SessionDetailQueryDao(context.Background(), &im.ReqSessionDetailQueryDao{
			SessionId: msg.WSMessage.SessionMessage.Session.SessionId,
		})
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		notify = models.WSMessageNotify{
			WSMessageNotifyType: msg.WSMessageNotifyType,
			Receive:             accountId,
			WSMessage:           msg.WSMessage,
			Timestamp:           msg.Timestamp,
		}
		var users = make([]rsp.User, 0)
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
	}
	writeChan <- notify
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
