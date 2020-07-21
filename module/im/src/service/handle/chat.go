package handle

import (
	"github.com/3115826227/baby-fried-rice/module/im/src/log"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type ChatFriend struct {
	First  *websocket.Conn
	Second *websocket.Conn
}

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
)

func init() {
}

/*
	聊天
*/
func ChatHandle(c *gin.Context) {
	token := c.Query("token")
	userMeta, err := GetUserMetaByToken(c, token)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Logger.Warn("connect failed")
		return
	}
	defer conn.Close()
	log.Logger.Info("connect success")
	ConnectionMap[userMeta.UserId] = conn

	for {
		var messageSend model.ChatMessageSend
		if err := conn.ReadJSON(&messageSend); err != nil {
			log.Logger.Error(err.Error())
			return
		}
		if err := MessageStorage(messageSend); err != nil {
			log.Logger.Warn(err.Error())
			continue
		}
		var messageReceive = model.ChatMessageReceive{
			MessageType: messageSend.MessageType,
			Body:        messageSend.Body,
			Timestamp:   messageSend.Timestamp,
			GroupID:     messageSend.GroupID,
			Sender:      messageSend.Sender,
			Receive:     userMeta.UserId,
		}
		if err := conn.WriteJSON(&messageReceive); err != nil {
			log.Logger.Error(err.Error())
			conn.Close()
			return
		}
	}
}
