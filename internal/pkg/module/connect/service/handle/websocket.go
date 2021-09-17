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
		closeChan <- true
		conn.Close()
	}()
	log.Logger.Info(userMeta.AccountId + " connect success")
	ConnectionMap[userMeta.AccountId] = conn
	for {
		var msg models.WSMessageNotify
		if err = conn.ReadJSON(&msg); err != nil {
			log.Logger.Error(err.Error())
			continue
		}
		switch msg.WSMessageNotifyType {
		case constant.SessionMessageNotify:
			switch msg.WSMessage.SessionMessage.SessionMessageType {
			case constant.SessionMessageMessage:
				handleSessionMessage(msg, userMeta.AccountId)
			}
		case constant.SpaceMessageNotify:
		default:
			continue
		}
	}
}

func handleSessionMessage(msg models.WSMessageNotify, accountId string) {
	imClient, err := grpc.GetClientGRPC(config.GetConfig().Rpc.SubServers.ImDaoServer)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	accountClient, err := grpc.GetClientGRPC(config.GetConfig().Rpc.SubServers.AccountDaoServer)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	// 从accountDao服务获取消息发送者的用户名和头像
	var userResp *user.RspDaoUserDetail
	userResp, err = user.NewDaoUserClient(accountClient.GetRpcClient()).
		UserDaoDetail(context.Background(), &user.ReqDaoUserDetail{AccountId: accountId})
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	msg.WSMessage.Send.AccountId = userResp.Detail.AccountId
	msg.WSMessage.Send.Username = userResp.Detail.Username
	msg.WSMessage.Send.HeadImgUrl = userResp.Detail.HeadImgUrl
	msg.Timestamp = time.Now().Unix()
	// 从imDao服务中获取会话中的所有用户id
	var resp *im.RspSessionDetailQueryDao
	resp, err = im.NewDaoImClient(imClient.GetRpcClient()).
		SessionDetailQueryDao(context.Background(), &im.ReqSessionDetailQueryDao{SessionId: msg.WSMessage.SessionMessage.Message.SessionId})
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
		switch notify.WSMessage.WSMessageType {

		}
		writeChan <- notify
	}
	// 将会话消息发给imDao服务存入数据库中
	var req = &im.ReqSessionMessageAddDao{
		MessageType:   msg.WSMessage.WSMessageType,
		Send:          msg.WSMessage.Send.AccountId,
		SessionId:     msg.WSMessage.SessionMessage.Message.SessionId,
		Content:       msg.WSMessage.SessionMessage.Message.Content,
		SendTimestamp: msg.Timestamp,
	}
	_, err = im.NewDaoImClient(imClient.GetRpcClient()).
		SessionMessageAddDao(context.Background(), req)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
}
