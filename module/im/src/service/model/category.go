package model

import (
	"github.com/3115826227/baby-fried-rice/module/im/src/log"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model/db"
)

func FindFriendCategory(categoryId int) (category FriendCategoryMeta, err error) {
	if err = db.GetDB().Where("id = ?", categoryId).First(&categoryId).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}

func FindFriendCategoryRelationByRelationId(relationIds ...string) (category []FriendCategoryRelation, err error) {
	if err = db.GetDB().Where("friend_relation_id in (?)", relationIds).Find(&category).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}

func GetFriendCategoryRelation(categoryId int) (relations []FriendCategoryRelation, err error) {
	if err = db.GetDB().Where("category_id = ?", categoryId).Find(&relations).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}

func GetFriendCategory(origin string) (category []FriendCategoryMeta, err error) {
	if err = db.GetDB().Where("origin = ?", origin).Find(&category).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}
