package handle

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/3115826227/baby-fried-rice/module/im/src/config"
	"github.com/3115826227/baby-fried-rice/module/im/src/kafka"
	"github.com/3115826227/baby-fried-rice/module/im/src/log"
	"github.com/3115826227/baby-fried-rice/module/im/src/redis"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model/db"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	//ChatMessageMap = make(map[string]chan model.FriendChatMessageReq)
	ChatChan = make(chan *model.FriendMessage, 5000)
	ChatRead = make(chan []string, 5000)
)

func init() {
	go func() {
		/*
			异步接收聊天信息
		*/
		kafka.ReceiveChat(config.ChatTopic)
	}()

	go func() {
		StorageMessage()
	}()

	go func() {
		UpdateMessageRead()
	}()

	go func() {
		for {
			select {
			case message := <-kafka.ChatCh:
				var relation model.FriendRelation
				if err := db.GetDB().Debug().Model(&model.FriendRelation{}).Where("friend = ? and origin = ?", message.Origin, message.Friend).Find(&relation).Error; err != nil {
					log.Logger.Warn(err.Error())
					continue
				}
				message.Remark = relation.FriendRemark
				var msg = &model.FriendMessage{
					Id:        GenerateID(),
					Content:   message.Content,
					Timestamp: message.CreateTime,
					Friend:    message.Friend,
					Origin:    message.Origin,
					Read:      false,
				}
				if _, exist := ConnectionMap[message.Friend]; exist {
					if err := ConnectionMap[message.Friend].WriteJSON(&message); err != nil {
						log.Logger.Warn(err.Error())
					} else {
						msg.Read = true
					}
				}
				ChatChan <- msg
			}
		}
	}()
}

func FriendHistoryMessageGet(c *gin.Context) {
	userMeta := GetUserMeta(c)
	friend := c.Query("friend")
	if friend == "" {
		c.JSON(http.StatusBadRequest, paramErrResponse)
		return
	}
	size := c.Query("size")
	var limit = config.DefaultMessageSize
	if size != "" {
		limit, _ = strconv.Atoi(size)
	}
	fmt.Println(limit)

	object := []string{userMeta.UserId, friend}
	var message = make([]model.FriendMessage, 0)
	if err := db.GetDB().Debug().Model(&model.FriendMessage{}).Raw(`select * from im_friend_message where origin in (?) and friend in (?) 
order by timestamp desc`, object, object).Limit(limit).Find(&message).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	var index = make([]string, 0)
	var rsp = make([]model.FriendMessage, 0)
	for i := len(message) - 1; i >= 0; i-- {
		rsp = append(rsp, message[i])
		if !message[i].Read && message[i].Friend == userMeta.UserId {
			index = append(index, message[i].Id)
		}
	}
	ChatRead <- index
	SuccessResp(c, "", rsp)
}

/*
	好友聊天
*/
func FriendChatHandle(c *gin.Context) {

	token := c.Query("token")
	userMeta, err := GetUserMetaByToken(c, token)
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	defer func() {
		conn.Close()
	}()
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	if err != nil {
		log.Logger.Warn("connect failed")
		return
	} else {
		log.Logger.Info("connect success")
		ConnectionMap[userMeta.UserId] = conn
	}

	for {
		var message model.FriendChatMessageReq
		if err := conn.ReadJSON(&message); err != nil {
			log.Logger.Warn(err.Error())
			return
		}
		userMeta, err = GetUserMetaByToken(c, message.Token)
		if err != nil {
			conn.Close()
			return
		} else {
			message.Token = ""
			message.Connect = true
		}
		var count = 0
		if err := db.GetDB().Model(&model.FriendRelation{}).Where("friend = ? and origin = ?", message.Friend, message.Origin).Count(&count).Error; err != nil {
			log.Logger.Warn(err.Error())
			continue
		}
		if count <= 0 {
			continue
		}
		if err := conn.WriteJSON(&message); err != nil {
			log.Logger.Warn(err.Error())
			return
		}
		kafka.Send(message.ToString(), message.Friend, config.ChatTopic)
	}
}

func UpdateMessageRead() {
	for {
		select {
		case index := <-ChatRead:
			if err := db.GetDB().Debug().Model(&model.FriendMessage{}).Where("id in (?)", index).Update("read", true).Error; err != nil {
				log.Logger.Warn(err.Error())
			}
		}
	}
}

func StorageMessage() {
	for {
		select {
		case msg := <-ChatChan:
			if err := db.GetDB().Debug().Create(&msg).Error; err != nil {
				log.Logger.Warn(err.Error())
				continue
			}
			var relation model.FriendRelation
			if err := db.GetDB().Model(&model.FriendRelation{}).Where("friend = ? and origin = ?", msg.Friend, msg.Origin).Find(&relation).Error; err != nil {
				log.Logger.Warn(err.Error())
				continue
			}
			redis.HashAdd(fmt.Sprintf("%v:%v", config.ChatNewMessageKey, msg.Origin), relation.ID, msg.ToString())

			var opRelation model.FriendRelation
			if err := db.GetDB().Model(&model.FriendRelation{}).Where("friend = ? and origin = ?", msg.Origin, msg.Friend).Find(&opRelation).Error; err != nil {
				log.Logger.Warn(err.Error())
				continue
			}
			redis.HashAdd(fmt.Sprintf("%v:%v", config.ChatNewMessageKey, msg.Friend), opRelation.ID, msg.ToString())
		}
	}
}

func ChatMessageNewGet(c *gin.Context) {
	userMeta := GetUserMeta(c)

	var chatMap = make(map[string]model.RspChat)
	var rsp = make([]model.RspChat, 0)
	mp, err := redis.HashGet(fmt.Sprintf("%v:%v", config.ChatNewMessageKey, userMeta.UserId))
	if err != nil {
		log.Logger.Warn(err.Error())
		SuccessResp(c, "", rsp)
		return
	}
	var relationIndex = make([]string, 0)
	for key, value := range mp {
		var msg model.FriendMessage
		if err := json.Unmarshal([]byte(value), &msg); err != nil {
			log.Logger.Warn(err.Error())
			continue
		}
		relationIndex = append(relationIndex, key)
		var chatTo = msg.Friend
		if msg.Friend == userMeta.UserId {
			chatTo = msg.Origin
		}
		chat := model.RspChat{
			Origin:    msg.Origin,
			Friend:    msg.Friend,
			Id:        msg.Id,
			ChatTo:    chatTo,
			Types:     0,
			Remark:    "",
			Read:      msg.Read,
			Content:   msg.Content,
			Timestamp: msg.Timestamp,
		}
		chatMap[key] = chat
	}
	var friends = make([]model.FriendRelation, 0)
	if err := db.GetDB().Debug().Where("id in (?)", relationIndex).Find(&friends).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	var friendIndex = make([]string, 0)
	for _, item := range friends {
		friendIndex = append(friendIndex, item.Friend)
	}
	var reads = make([]struct {
		Friend string
		UnRead int
	}, 0)
	if err := db.GetDB().Debug().Model(&model.FriendMessage{}).Raw(`select a.friend, count(*) from im_friend_message as a
where a.friend in (?) and a.origin = ? and a.read=0 group by a.friend`, friendIndex, userMeta.UserId).Find(&reads).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	var friendUnreadMap = make(map[string]int)
	for _, data := range reads {
		friendUnreadMap[data.Friend] = data.UnRead
	}
	for _, item := range friends {
		var data = chatMap[item.ID]
		data.Remark = item.FriendRemark
		chatMap[item.ID] = data
		rsp = append(rsp, chatMap[item.ID])
	}
	for index, data := range rsp {
		rsp[index].More = friendUnreadMap[data.Friend]
	}
	sort.Sort(model.RspChats(rsp))
	SuccessResp(c, "", rsp)
}
