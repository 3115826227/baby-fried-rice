package model

type FriendChatMessageReq struct {
	Friend     string `json:"friend"`
	Content    string `json:"content"`
	CreateTime int    `json:"create_time"`
	IsFriend   bool   `json:"is_friend"`
	Status     bool   `json:"status"`
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
