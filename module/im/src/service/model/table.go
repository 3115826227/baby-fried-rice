package model

import (
	"encoding/json"
	"github.com/3115826227/baby-fried-rice/module/im/src/log"
	"github.com/3115826227/baby-fried-rice/module/im/src/service/model/db"
	"time"
)

func init() {
	err := db.GetDB().AutoMigrate(
		&FriendGroupCategoryMeta{},
		&FriendGroupCategoryRelation{},
		&FriendGroupMeta{},
		&FriendGroupLevelMeta{},
		&FriendGroupRelation{},
		&FriendGroupMessage{},
		&FriendCategoryMeta{},
		&FriendCategoryRelation{},
		&FriendRelation{},
		&FriendMessage{},
		&UserDetail{},
		&FriendAddMeta{},
		&FriendAddPermission{},
	).Error
	if err != nil {
		log.Logger.Warn(err.Error())
	}
}

type IntCommonField struct {
	ID        int       `gorm:"AUTO_INCREMENT;column:id;" json:"id"`
	CreatedAt time.Time `gorm:"column:create_time;type:timestamp" json:"-"`
	UpdatedAt time.Time `gorm:"column:update_time;type:timestamp" json:"-"`
}

type StringCommonField struct {
	ID        string    `gorm:"column:id;type:char(36);primary_key;not null" json:"id"`
	CreatedAt time.Time `gorm:"column:create_time;type:timestamp" json:"-"`
	UpdatedAt time.Time `gorm:"column:update_time;type:timestamp" json:"-"`
}

/*
	用户表
*/
type UserDetail struct {
	UserId    string `gorm:"column:user_id;unique"`
	AccountId string `gorm:"column:account_id;not null"`
}

func (table *UserDetail) TableName() string {
	return "im_user_detail"
}

//群分类元信息表
type FriendGroupCategoryMeta struct {
	IntCommonField

	Name   string `gorm:"column:name;unique_index:idx_friend_group_category"`
	Origin string `gorm:"column:origin;unique_index:idx_friend_group_category"`
}

func (table *FriendGroupCategoryMeta) TableName() string {
	return "im_friend_group_category_meta"
}

//群分类关联表
type FriendGroupCategoryRelation struct {
	FriendGroupCategoryId int    `gorm:"friend_group_category_id;unique_index:idx_friend_group_category_ref_group_category"`
	FriendGroupId         string `gorm:"friend_group_id;unique_index:idx_friend_group_category_ref_group_category"`
}

func (table *FriendGroupCategoryRelation) TableName() string {
	return "im_friend_group_category_relation"
}

//群元信息表
type FriendGroupMeta struct {
	StringCommonField

	Name               string `gorm:"column:name;unique_index:idx_friend_group_name_origin" json:"name"`
	Level              string `gorm:"column:level" json:"level"`
	Official           bool   `gorm:"column:official" json:"official"`
	SchoolDepartmentId string `gorm:"column:school_department_id" json:"school_department_id"`
	Origin             string `gorm:"column:origin;unique_index:idx_friend_group_name_origin" json:"origin"`
}

func (table *FriendGroupMeta) TableName() string {
	return "im_friend_group_meta"
}

type FriendGroupLevelMeta struct {
	IntCommonField

	Name        string `gorm:"column:name;unique"`
	ParentId    int    `gorm:"column:parent_id"`
	PersonLimit int    `gorm:"column:person_limit"`
}

func (table *FriendGroupLevelMeta) TableName() string {
	return "im_friend_group_level_meta"
}

//群好友列表
type FriendGroupRelation struct {
	GroupId         string `gorm:"column:group_id;unique_index:idx_friend_group_ref_group_user"`
	UserId          string `gorm:"column:user_id;unique_index:idx_friend_group_ref_group_user"`
	UserGroupRemark string `gorm:"column:user_group_remark"`
}

func (table *FriendGroupRelation) TableName() string {
	return "im_friend_group_relation"
}

//群消息表
type FriendGroupMessage struct {
	IntCommonField

	GroupId   string `gorm:"column:group_id"`
	Content   string `gorm:"column:content"`
	Timestamp int64  `gorm:"column:timestamp"`
	Sender    string `gorm:"column:sender"`
}

func (table *FriendGroupMessage) TableName() string {
	return "im_friend_group_message"
}

//群消息接收人表
type FriendGroupMessageReceive struct {
	GroupId   string `json:"group_id"`
	MessageId int    `json:"message_id"`
	UserId    string `json:"user_id"`
	Read      bool   `json:"read"`
}

//好友分类元信息表
type FriendCategoryMeta struct {
	IntCommonField

	Name   string `gorm:"column:name;unique_index:idx_friend_category_name_origin" json:"name"`
	Origin string `gorm:"column:origin;unique_index:idx_friend_category_name_origin" json:"-"`
}

func (table *FriendCategoryMeta) TableName() string {
	return "im_friend_category_meta"
}

type FriendCategories []FriendCategoryMeta

func (list FriendCategories) Len() int {
	return len(list)
}

func (list FriendCategories) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list FriendCategories) Less(i, j int) bool {
	return list[i].Name < list[j].Name
}

//好友分类关系表
type FriendCategoryRelation struct {
	CategoryId       int    `gorm:"column:category_id;unique_index:idx_friend_category_ref_category_friend"`
	FriendRelationId string `gorm:"column:friend_relation_id;unique_index:idx_friend_category_ref_category_friend"`
}

func (table *FriendCategoryRelation) TableName() string {
	return "im_friend_category_relation"
}

//好友关系表
type FriendRelation struct {
	StringCommonField

	Friend       string `gorm:"column:friend;unique_index:idx_friend_ref_category_friend_origin"`
	FriendRemark string `gorm:"column:friend_remark"`
	Origin       string `gorm:"column:origin;unique_index:idx_friend_ref_category_friend_origin"`
}

func (table *FriendRelation) TableName() string {
	return "im_friend_relation"
}

//好友消息表
type FriendMessage struct {
	IntCommonField
	Content   string `gorm:"column:content" json:"content"`
	Timestamp int64  `gorm:"column:timestamp" json:"timestamp"`
	Receive   string `gorm:"column:receive" json:"receive"`
	Sender    string `gorm:"column:sender" json:"sender"`
	Read      bool   `gorm:"column:read" json:"read"`
}

func (table *FriendMessage) ToString() string {
	data, _ := json.Marshal(table)
	return string(data)
}

func (table *FriendMessage) TableName() string {
	return "im_friend_message"
}

type UserUnreadMessage struct {
	UserId      string `json:"user_id"`
	MessageType int    `json:"message_type"`
	Friend      string `json:"friend"`
	Group       string `json:"group"`
	Unread      int    `json:"unread"`
	LastMessage int    `json:"last_message"`
}

func (table *UserUnreadMessage) TableName() string {
	return "im_user_unread_message"
}

/*
	添加好友请求记录表
*/
type FriendAddMeta struct {
	IntCommonField

	AskOrReceive bool   `json:"ask_or_receive"`
	Origin       string `json:"origin"`
	Friend       string `json:"friend"`
	//0-忽略，1-同意，2-拒绝
	Agree int `json:"agree"`
}

func (table *FriendAddMeta) TableName() string {
	return "im_friend_add_meta"
}

/*
	添加好友权限表
*/
type FriendAddPermission struct {
	StringCommonField
	//0-无需认证，1-需要同意方可添加，默认0
	Permission int `json:"permission"`
}

func (table *FriendAddPermission) TableName() string {
	return "im_friend_add_permission"
}
