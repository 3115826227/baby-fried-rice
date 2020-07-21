package handle

import (
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model/db"
)

/*
	消息持久化
*/
func MessageStorage(message model.ChatMessageSend) error {
	var beans = make([]interface{}, 0)
	switch message.MessageType {
	case 1:
		for _, r := range message.Receive {
			var friendMessage = model.FriendMessage{
				Content:   string(message.Body),
				Timestamp: message.Timestamp,
				Receive:   r,
				Sender:    message.Sender,
			}
			beans = append(beans, &friendMessage)

			if r != message.Sender {
				var messageReceive = model.ChatMessageReceive{
					MessageType: message.MessageType,
					Body:        message.Body,
					Timestamp:   message.Timestamp,
					GroupID:     message.GroupID,
					Sender:      message.Sender,
					Receive:     r,
				}
				go SyncMessage(messageReceive)
			}
		}
	case 2:
		var groupMessage = model.FriendGroupMessage{
			GroupId:   message.GroupID,
			Content:   string(message.Body),
			Timestamp: message.Timestamp,
			Sender:    message.Sender,
		}
		beans = append(beans, &groupMessage)
		usersInfo := model.GetGroupUsersInfo(message.GroupID)
		for _, info := range usersInfo {
			var m = model.FriendGroupMessageReceive{
				GroupId:   groupMessage.GroupId,
				MessageId: groupMessage.ID,
				UserId:    info.UserId,
			}
			beans = append(beans, &m)

			if info.UserId != groupMessage.Sender {
				var messageReceive = model.ChatMessageReceive{
					MessageType: message.MessageType,
					Body:        message.Body,
					Timestamp:   message.Timestamp,
					GroupID:     message.GroupID,
					Sender:      message.Sender,
					Receive:     info.UserId,
				}
				go SyncMessage(messageReceive)
			}
		}
	}
	return db.CreateMulti(beans...)
}
