package handle

import (
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/im/src/log"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetHistoryMessageHandle(c *gin.Context) {
	userMeta := GetUserMeta(c)
	messageTypeStr := c.Query("message_type")
	messageType, err := strconv.Atoi(messageTypeStr)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	switch messageType {
	case 1:
		friend := c.Query("friend")
		if friend == "" {
			return
		}
		messageReceives := model.GetFriendHistoryMessage(userMeta.UserId, friend)
		fmt.Println(messageReceives)
	case 2:
		group := c.Query("group")
		if group == "" {
			return
		}
		messageReceives := model.GetGroupHistoryMessage(userMeta.UserId, group)
		fmt.Println(messageReceives)
	}
}
