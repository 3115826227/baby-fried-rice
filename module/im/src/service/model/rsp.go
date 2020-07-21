package model

type RspSuccess struct {
	Code int `json:"code"`
}

type RspOkResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RespSuccessData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type RspFriendCategory struct {
	Id      string      `json:"id"`
	Name    string      `json:"name"`
	Friends []RspFriend `json:"friends"`
}

type RspFriend struct {
	Id       string `json:"id"`
	Friend   string `json:"friend"`
	Username string `json:"username"`
	Remark   string `json:"remark"`
	HeadImg  string `json:"head_img"`
}

type RspMessage struct {
	Origin    string `json:"origin"`
	Id        string `json:"id"`
	Types     int    `json:"types"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

type RspChat struct {
	Origin    string `json:"origin"`
	Friend    string `json:"friend"`
	ChatTo    string `json:"chat_to"`
	Id        int    `json:"id"`
	Types     int    `json:"types"`
	Remark    string `json:"remark"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
	Read      bool   `json:"read"`
	More      int    `json:"more"`
}

type RspChats []RspChat

func (rsp RspChats) Len() int {
	return len(rsp)
}

func (rsp RspChats) Swap(i, j int) {
	rsp[i], rsp[j] = rsp[j], rsp[i]
}

func (rsp RspChats) Less(i, j int) bool {
	return rsp[i].Timestamp > rsp[j].Timestamp
}
