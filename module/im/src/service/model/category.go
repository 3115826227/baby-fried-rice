package model

import (
	"github.com/3115826227/baby-fried-rice/module/im/src/log"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model/db"
)

func FindFriendCategory(categoryId int) (category FriendCategoryMeta, err error) {
	if err = db.GetDB().Debug().Where("id = ?", categoryId).First(&categoryId).Error; err != nil {
		log.Logger.Error(err.Error())
	}
	return
}

func FindFriendCategoryRelationByRelationId(relationIds ...string) (category []FriendCategoryRelation, err error) {
	if err = db.GetDB().Debug().Where("friend_relation_id in (?)", relationIds).Find(&category).Error; err != nil {
		log.Logger.Error(err.Error())
	}
	return
}

func GetFriendCategoryRelation(categoryId string) (relations []FriendCategoryRelation, err error) {
	if err = db.GetDB().Debug().Where("category_id = ?", categoryId).Find(&relations).Error; err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}

func GetFriendCategory(origin string) (category []FriendCategoryMeta, err error) {
	if err = db.GetDB().Debug().Where("origin = ?", origin).Find(&category).Error; err != nil {
		log.Logger.Error(err.Error())
	}
	return
}
