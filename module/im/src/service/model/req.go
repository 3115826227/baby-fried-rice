package model

import "encoding/json"

type FriendChatMessageReq struct {
	Origin     string `json:"origin"`
	Friend     string `json:"friend"`
	Token      string `json:"token"`
	Content    string `json:"content"`
	Remark     string `json:"remark"`
	CreateTime int64  `json:"create_time"`
	IsFriend   bool   `json:"is_friend"`
	Status     bool   `json:"status"`
	Connect    bool   `json:"connect"`
}

func (msg *FriendChatMessageReq) ToString() string {
	data, _ := json.Marshal(msg)
	return string(data)
}

type FriendCategoryAddReq struct {
	Name string `json:"name"`
}

type FriendCategoryUpdateReq struct {
	CategoryId int    `json:"category_id"`
	Name       string `json:"name"`
}

type FriendAddReq struct {
	UserId string `json:"user_id"`
	Remark string `json:"remark"`
}

type FriendRemarkUpdateReq struct {
	Id     string `json:"id"`
	Remark string `json:"remark"`
}
