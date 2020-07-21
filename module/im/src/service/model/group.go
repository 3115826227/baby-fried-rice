package model

import (
	"github.com/3115826227/baby-fried-rice/module/im/src/log"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model/db"
)

type GroupUserInfo struct {
	UserId          string `json:"user_id"`
	UserGroupRemark string `json:"user_group_remark"`
}

func GetGroupUsersInfo(group string) (users []GroupUserInfo) {
	users = make([]GroupUserInfo, 0)
	relations := make([]FriendGroupRelation, 0)
	if err := db.GetDB().Debug().Where("group_id = ?", group).Find(&relations).Error; err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	for _, r := range relations {
		users = append(users, GroupUserInfo{
			UserId:          r.UserId,
			UserGroupRemark: r.UserGroupRemark,
		})
	}
	return
}
