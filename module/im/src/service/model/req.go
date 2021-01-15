package model

import "encoding/json"

type ChatMessageSend struct {
	Token       string `json:"token"`        //用来校验消息是否真实
	MessageType int    `json:"message_type"` //消息类型 1-个人消息，2-群组消息
	BodyType    int    `json:"body_type"`    //消息内容类型
	MessageBody string `json:"message_body"` //文字类消息内容
	Image       bool   `json:"image"`        //是否为图片
	Body        []byte `json:"body"`         //其他类消息内容
	Timestamp   int64  `json:"timestamp"`    //消息产生时间
	GroupID     string `json:"group_id"`     //消息群组id
	Sender      string `json:"sender"`       //消息发送者ID
	Receive     string `json:"receive"`      //消息接受者ID
}

func (message *ChatMessageSend) ToString() string {
	data, _ := json.Marshal(message)
	return string(data)
}

type ChatMessageReceive struct {
	MessageID    int    `json:"message_id"`
	MessageType  int    `json:"message_type"`
	Image        bool   `json:"image"`
	Body         []byte `json:"body"`          //消息内容
	MessageBody  string `json:"message_body"`  //文字类消息内容
	Timestamp    int64  `json:"timestamp"`     //消息产生时间
	GroupID      string `json:"group_id"`      //消息群组id
	Sender       string `json:"sender"`        //消息发送者ID
	SenderRemark string `json:"sender_remark"` //消息发送者昵称
	Receive      string `json:"receive"`       //消息接受者ID
}

func (message *ChatMessageReceive) ToString() string {
	data, _ := json.Marshal(message)
	return string(data)
}

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
	CategoryId int    `json:"category_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
}

type FriendAddReq struct {
	AccountId string `json:"account_id"`
	Remark    string `json:"remark"`
	Category  string `json:"category"`
}

type FriendRemarkUpdateReq struct {
	Id     string `json:"id"`
	Remark string `json:"remark"`
}

type GroupAddReq struct {
	Friends []struct {
		Id string `json:"id"`
	} `json:"friends"`
	Name string `json:"name"`
}

type ReqOfficialGroupAdd struct {
	Organize string `json:"organize" binding:"required"`
	Name     string `json:"name" binding:"required"`
}
