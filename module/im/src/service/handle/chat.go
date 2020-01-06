package handle

import (
	"github.com/3115826227/baby-fried-rice/module/im/src/log"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model/db"
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
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ConnectionMap = make(map[string]*websocket.Conn)
)

/*
	好友聊天
*/
func FriendChatHandle(c *gin.Context) {
	userMeta, err := GetUserMetaByToken(c)
	if err != nil {
		return
	}

	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	defer conn.Close()

	if _, exist := ConnectionMap[userMeta.UserId]; !exist {
		ConnectionMap[userMeta.UserId] = conn
	}

	var mainDB = db.GetDB()
	for {
		message := &model.FriendChatMessageReq{}
		var friendRelations = make([]model.FriendRelation, 0)
		if err := mainDB.Where("origin = ? and friend = ?", userMeta.UserId, message.Friend).Find(&friendRelations).Error; err != nil {
			log.Logger.Warn(err.Error())
			message.Status = false
		}
		if len(friendRelations) == 0 {
			message.IsFriend = false
		}

		if err := conn.ReadJSON(message); err != nil {
			log.Logger.Warn(err.Error())
			return
		}
		if err := conn.WriteJSON(message); err != nil {
			log.Logger.Warn(err.Error())
			return
		}
		if !message.IsFriend || !message.Status {
			continue
		}
		if _, exist := ConnectionMap[message.Friend]; exist {
			if err := ConnectionMap[message.Friend].WriteJSON(message); err != nil {
				log.Logger.Warn(err.Error())
				return
			}
		}
	}
}
