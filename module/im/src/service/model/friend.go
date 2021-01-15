package model

import (
	"github.com/3115826227/baby-fried-rice/module/im/src/log"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model/db"
)

type FriendCategoryResult struct {
	Id           string `json:"id"`
	Friend       string `json:"friend"`
	FriendRemark string `json:"friend_remark"`
	CategoryName string `json:"category_name"`
	Username     string `json:"username"`
}

func GetFriend(origin string) (res []FriendCategoryResult, err error) {
	err = db.GetDB().Raw(`select a.id, a.friend,a.friend_remark, c.name as category_name, d.username from im_friend_relation as a 
left join im_friend_category_relation as b on a.id = b.friend_relation_id
left join im_friend_category_meta as c on b.category_id = c.id  
left join im_user_detail as d on a.friend = d.user_id where a.origin = ?`, origin).Scan(&res).Error
	if err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}

func FindFriendRelationById(id string) (relation FriendRelation, err error) {
	if err = db.GetDB().Debug().Where("id = ?", id).First(&relation).Error; err != nil {
		log.Logger.Error(err.Error())
	}
	return
}

func FindFriendRelation(origin, friend string) (relation FriendRelation, err error) {
	if err = db.GetDB().Debug().Where("origin = ? and friend = ?", origin, friend).First(&relation).Error; err != nil {
		log.Logger.Error(err.Error())
	}
	return
}
