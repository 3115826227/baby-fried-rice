package model

import (
	"github.com/3115826227/baby-fried-rice/module/im/src/log"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model/db"
)

// 获取好友历史消息
func GetFriendHistoryMessage(user, friend string) (messageReceives []ChatMessageReceive) {
	messageReceives = make([]ChatMessageReceive, 0)
	var messages = make([]FriendMessage, 0)
	if err := db.GetDB().Debug().Where("sender in (?)", []string{user, friend}).Find(messages).Error; err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	go FriendMessageReadSync(user, messages)
	for _, m := range messages {
		messageReceives = append(messageReceives, ChatMessageReceive{
			MessageID:   m.ID,
			MessageType: 1,
			Body:        []byte(m.Content),
			Timestamp:   m.Timestamp,
			Sender:      m.Sender,
			Receive:     m.Receive,
		})
	}
	return
}

//获取群历史消息
func GetGroupHistoryMessage(user, group string) (messageReceives []ChatMessageReceive) {
	messageReceives = make([]ChatMessageReceive, 0)
	var messages = make([]FriendGroupMessage, 0)
	var receive = make([]FriendGroupMessageReceive, 0)
	var messageIDs = make([]int, 0)
	if err := db.GetDB().Debug().Where("group_id = ?", group).Find(&messages).Error; err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	for _, m := range messages {
		messageIDs = append(messageIDs, m.ID)
	}
	if err := db.GetDB().Debug().Where("message_id in (?) and user_id = ?", messageIDs, user).Find(&receive).Error; err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	var filterMessageIDs = make([]int, 0)
	var filterMessages = make([]FriendGroupMessage, 0)
	for _, r := range receive {
		filterMessageIDs = append(filterMessageIDs, r.MessageId)
	}
	if err := db.GetDB().Debug().Where("id in (?)", filterMessageIDs).Find(&filterMessages).Error; err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	go GroupMessageReadSync(user, filterMessageIDs)
	for _, m := range filterMessages {
		messageReceives = append(messageReceives, ChatMessageReceive{
			MessageID:   m.ID,
			MessageType: 2,
			Body:        []byte(m.Content),
			Timestamp:   m.Timestamp,
			GroupID:     m.GroupId,
			Sender:      m.Sender,
			Receive:     user,
		})
	}
	return
}

func FriendMessageReadSync(user string, friendMessages []FriendMessage) {
	var filterMessageIDs = make([]int, 0)
	for _, m := range friendMessages {
		if m.Sender == user {
			continue
		}
		filterMessageIDs = append(filterMessageIDs, m.ID)
	}
	if err := db.GetDB().Debug().Model(&FriendMessage{}).Where("id in (?)", filterMessageIDs).Update("read", true).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
}

func GroupMessageReadSync(user string, messageIDs []int) {
	if err := db.GetDB().Debug().Model(&FriendGroupMessageReceive{}).Where("message_id in (?) and user_id = ?", messageIDs, user).Update("read", true).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
}
