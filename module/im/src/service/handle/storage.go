package handle

import (
	"github.com/3115826227/baby-fried-rice/module/im/src/log"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model/db"
)

/*
	消息持久化
*/
func MessageStorage(message model.ChatMessageSend) (err error) {
	var beans = make([]interface{}, 0)
	switch message.MessageType {
	case 1:
		var friendMessage = model.FriendMessage{
			Body:        message.Body,
			MessageBody: message.MessageBody,
			MessageType: message.MessageType,
			Image:       message.Image,
			Timestamp:   message.Timestamp,
			Receive:     message.Receive,
			Sender:      message.Sender,
		}
		beans = append(beans, &friendMessage)

		var messageReceive = model.ChatMessageReceive{
			MessageType: message.MessageType,
			Body:        message.Body,
			Image:       message.Image,
			MessageBody: message.MessageBody,
			Timestamp:   message.Timestamp,
			Sender:      message.Sender,
			Receive:     message.Receive,
		}
		SyncMessage(messageReceive)
		return db.CreateMulti(beans...)
	case 2:
		tx := db.GetDB().Begin()
		if err = tx.Error; err != nil {
			return
		}
		defer func() {
			if err != nil {
				tx.Rollback()
				return
			}
			tx.Commit()
		}()
		usersInfo := model.GetGroupUsersInfo(message.GroupID)
		var senderRemark string
		for _, info := range usersInfo {
			if message.Sender == info.UserId {
				senderRemark = info.UserGroupRemark
			}
		}
		var groupMessage = model.FriendGroupMessage{
			GroupId:      message.GroupID,
			MessageType:  message.MessageType,
			MessageBody:  message.MessageBody,
			Body:         message.Body,
			Timestamp:    message.Timestamp,
			Sender:       message.Sender,
			SenderRemark: senderRemark,
		}
		if err = tx.Create(&groupMessage).Error; err != nil {
			log.Logger.Error(err.Error())
			return
		}
		for _, info := range usersInfo {
			var m = model.FriendGroupMessageReceive{
				GroupId:   groupMessage.GroupId,
				MessageId: groupMessage.ID,
				UserId:    info.UserId,
			}
			if err = tx.Create(&m).Error; err != nil {
				log.Logger.Error(err.Error())
				return
			}

			if info.UserId != groupMessage.Sender {
				var messageReceive = model.ChatMessageReceive{
					MessageType:  message.MessageType,
					Body:         message.Body,
					MessageBody:  message.MessageBody,
					Timestamp:    message.Timestamp,
					GroupID:      message.GroupID,
					Sender:       message.Sender,
					SenderRemark: senderRemark,
					Receive:      info.UserId,
				}
				go SyncMessage(messageReceive)
			}
		}
		return
	}
	return
}
