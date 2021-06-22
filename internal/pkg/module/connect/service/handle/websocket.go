package handle

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"baby-fried-rice/internal/pkg/kit/handle"
	"baby-fried-rice/internal/pkg/kit/interfaces"
	"baby-fried-rice/internal/pkg/kit/models"
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
	mq = nsq.InitNSQMQ(conf.NSQ.Addr)
	err := mq.NewConsumer(conf.NSQ.Topic, conf.NSQ.Channel)
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
				log.Logger.Debug(fmt.Sprintf("write message %v to %v success", msg.WSMessage.ToString(), msg.Receive))
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
	imClient, err := grpc.GetClientGRPC(config.GetConfig().Servers.ImDaoServer)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	accountClient, err := grpc.GetClientGRPC(config.GetConfig().Servers.AccountDaoServer)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, handle.SysErrResponse)
		return
	}
	for {
		var msg models.WSMessageNotify
		if err = conn.ReadJSON(&msg); err != nil {
			log.Logger.Error(err.Error())
			continue
		}
		switch msg.WSMessageNotifyType {
		case constant.SessionMessageNotify:
			// 从accountDao服务获取消息发送者的用户名和头像
			var userResp *user.RspDaoUserDetail
			userResp, err = user.NewDaoUserClient(accountClient.GetRpcClient()).
				UserDaoDetail(context.Background(), &user.ReqDaoUserDetail{AccountId: userMeta.AccountId})
			if err != nil {
				log.Logger.Error(err.Error())
				continue
			}
			msg.WSMessage.Send.AccountId = userResp.Detail.AccountId
			msg.WSMessage.Send.Username = userResp.Detail.Username
			msg.WSMessage.Send.HeadImgUrl = userResp.Detail.HeadImgUrl
			// 从imDao服务中获取会话中的所有用户id
			var resp *im.RspSessionDetailQueryDao
			resp, err = im.NewDaoImClient(imClient.GetRpcClient()).
				SessionDetailQueryDao(context.Background(), &im.ReqSessionDetailQueryDao{SessionId: msg.WSMessage.SessionId})
			if err != nil {
				log.Logger.Error(err.Error())
				continue
			}
			// 发送给nsq
			for _, u := range resp.Session.Joins {
				var notify = models.WSMessageNotify{
					WSMessageNotifyType: msg.WSMessageNotifyType,
					Receive:             u,
					WSMessage:           msg.WSMessage,
					Timestamp:           msg.Timestamp,
				}
				writeChan <- notify
			}
			// 将会话消息发给imDao服务存入数据库中
			var req = &im.ReqSessionMessageAddDao{
				MessageType:   msg.WSMessage.WSMessageType,
				Send:          msg.WSMessage.Send.AccountId,
				SessionId:     msg.WSMessage.SessionId,
				Content:       []byte(msg.WSMessage.Content),
				SendTimestamp: msg.Timestamp,
			}
			_, err = im.NewDaoImClient(imClient.GetRpcClient()).
				SessionMessageAddDao(context.Background(), req)
			if err != nil {
				log.Logger.Error(err.Error())
				continue
			}
		default:
			continue
		}
	}
}
